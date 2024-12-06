package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"project-survey-stats-writer/internal/configuration"
	"project-survey-stats-writer/internal/db"
	"project-survey-stats-writer/internal/enums"
	"project-survey-stats-writer/internal/events/model"
	"project-survey-stats-writer/internal/events/model/pb"
	"strconv"
	"time"
)

var (
	trackingEventsHeader       = []byte("date,timestamp,hour,unit_id,geo,lang,gender,unit_requests,unit_responses,unit_views,unit_ends,survey_id,survey_responses,mismatched_survey_responses,mismatched_survey_reason,survey_starts,survey_ends\n")
	completionEventsHeader     = []byte("date,timestamp,hour,survey_id,question_id,option_ids,geo,lang,gender,answered,views\n")
	completionFlatEventsHeader = []byte("date,timestamp,hour,survey_id,question_id,option_id,geo,lang,gender,answered,views\n")
)

type Proceeder struct {
	dbClient db.Client
	messages chan model.Message
	config   *configuration.AppConfiguration

	trackingEvents       int
	completionEvents     int
	completionFlatEvents int

	trackingEventsBuffer       *bytes.Buffer
	completionEventsBuffer     *bytes.Buffer
	completionFlatEventsBuffer *bytes.Buffer
}

func NewProceeder(dbClient db.Client, config *configuration.AppConfiguration, messages chan model.Message) Proceeder {
	return Proceeder{
		dbClient:                   dbClient,
		config:                     config,
		messages:                   messages,
		trackingEventsBuffer:       bytes.NewBuffer([]byte{}),
		completionEventsBuffer:     bytes.NewBuffer([]byte{}),
		completionFlatEventsBuffer: bytes.NewBuffer([]byte{}),
	}
}

func (p *Proceeder) Proceed() {
	for {
		select {
		case msg := <-p.messages:
			switch msg.Topic {
			case p.config.EventsConfiguration.TrackingEventsTopic:
				p.proceedTrackingEvent(msg.Value)
			case p.config.EventsConfiguration.CompletionEventsTopic:
				p.proceedCompletionEvent(msg.Value)
			}
		}
	}
}

func (p *Proceeder) Finalise() {
	if p.trackingEventsBuffer.Len() != 0 {
		p.dbClient.BulkInsert(p.trackingEventsBuffer, "pure_survey.tracking_stats")
		p.trackingEvents = 0
	}

	if p.completionEventsBuffer.Len() != 0 {
		p.dbClient.BulkInsert(p.completionEventsBuffer, "pure_survey.completion_stats")
		p.completionEvents = 0
	}

	if p.completionFlatEventsBuffer.Len() != 0 {
		p.dbClient.BulkInsert(p.completionFlatEventsBuffer, "pure_survey.completion_stats_flat")
		p.completionFlatEvents = 0
	}
}

func (p *Proceeder) proceedTrackingEvent(message []byte) {
	if p.trackingEventsBuffer.Len() == 0 {
		p.trackingEventsBuffer.Write(trackingEventsHeader)
	}

	event := &pb.TrackingEvent{}
	err := proto.Unmarshal(message, event)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	timestamp := time.Unix(event.Timestamp, 0)

	p.trackingEventsBuffer.Write([]byte(timestamp.Format("2006-01-02")))
	p.trackingEventsBuffer.WriteByte(',')

	p.trackingEventsBuffer.Write([]byte(timestamp.Format("2006-01-02 15:04:05")))
	p.trackingEventsBuffer.WriteByte(',')

	binary.Write(p.trackingEventsBuffer, binary.LittleEndian, timestamp.Hour())
	p.trackingEventsBuffer.WriteByte(',')

	binary.Write(p.trackingEventsBuffer, binary.LittleEndian, event.UnitId)
	p.trackingEventsBuffer.WriteByte(',')

	p.trackingEventsBuffer.Write([]byte(event.Geo))
	p.trackingEventsBuffer.WriteByte(',')

	p.trackingEventsBuffer.Write([]byte(event.Lang))
	p.trackingEventsBuffer.WriteByte(',')

	binary.Write(p.trackingEventsBuffer, binary.LittleEndian, event.Gender)
	p.trackingEventsBuffer.WriteByte(',')

	if int(event.EventType) == int(enums.ETUnitRequest) {
		p.trackingEventsBuffer.Write([]byte("1,"))
	} else {
		p.trackingEventsBuffer.Write([]byte("0,"))
	}

	if int(event.EventType) == int(enums.ETUnitResponse) {
		p.trackingEventsBuffer.Write([]byte("1,"))
	} else {
		p.trackingEventsBuffer.Write([]byte("0,"))
	}

	if int(event.EventType) == int(enums.ETUnitView) {
		p.trackingEventsBuffer.Write([]byte("1,"))
	} else {
		p.trackingEventsBuffer.Write([]byte("0,"))
	}

	if int(event.EventType) == int(enums.ETUnitEnd) {
		p.trackingEventsBuffer.Write([]byte("1,"))
	} else {
		p.trackingEventsBuffer.Write([]byte("0,"))
	}

	p.getStringArray(event.SurveyIds).WriteTo(p.trackingEventsBuffer)

	var surveyResponses []int32
	var mismatchedResponses []int32
	for _, rsn := range event.MismatchReason {
		if int(rsn) == int(enums.ETUnknown) {
			surveyResponses = append(surveyResponses, 1)
			mismatchedResponses = append(mismatchedResponses, 0)
		} else {
			surveyResponses = append(surveyResponses, 0)
			mismatchedResponses = append(mismatchedResponses, 1)
		}
	}
	p.getStringArray(surveyResponses).WriteTo(p.trackingEventsBuffer)
	p.getStringArray(mismatchedResponses).WriteTo(p.trackingEventsBuffer)
	p.getStringArray(event.MismatchReason).WriteTo(p.trackingEventsBuffer)

	if int(event.EventType) == int(enums.ETSurveyStart) {
		p.trackingEventsBuffer.Write([]byte("'{1}',"))
	} else {
		p.trackingEventsBuffer.Write([]byte("'{0}',"))
	}

	if int(event.EventType) == int(enums.ETSurveyEnd) {
		p.trackingEventsBuffer.Write([]byte("'{1}'"))
	} else {
		p.trackingEventsBuffer.Write([]byte("'{0}'"))
	}
	p.trackingEventsBuffer.WriteByte('\n')
}

func (p *Proceeder) proceedCompletionEvent(message []byte) {
	if p.completionFlatEventsBuffer.Len() == 0 {
		p.completionFlatEventsBuffer.Write(completionFlatEventsHeader)
	}

	event := &pb.CompletionEvent{}
	err := proto.Unmarshal(message, event)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	for _, option := range event.OptionIds {
		timestamp := time.Unix(event.Timestamp, 0)

		writeDate(p.completionFlatEventsBuffer, timestamp)
		writeSeparator(p.completionFlatEventsBuffer)

		writeTimestamp(p.completionFlatEventsBuffer, timestamp)
		writeSeparator(p.completionFlatEventsBuffer)

		writeInt(p.completionFlatEventsBuffer, timestamp.Hour())
		writeSeparator(p.completionFlatEventsBuffer)

		writeInt(p.completionFlatEventsBuffer, int(event.SurveyId))
		writeSeparator(p.completionFlatEventsBuffer)

		writeInt(p.completionFlatEventsBuffer, int(event.QuestionId))
		writeSeparator(p.completionFlatEventsBuffer)

		writeInt(p.completionFlatEventsBuffer, int(option))
		writeSeparator(p.completionFlatEventsBuffer)

		writeString(p.completionFlatEventsBuffer, event.Geo)
		writeSeparator(p.completionFlatEventsBuffer)

		writeString(p.completionFlatEventsBuffer, event.Lang)
		writeSeparator(p.completionFlatEventsBuffer)

		writeInt(p.completionFlatEventsBuffer, int(event.Gender))
		writeSeparator(p.completionFlatEventsBuffer)

		if int(event.EventType) == int(enums.ETQuestionAnswer) {
			write(p.completionFlatEventsBuffer, "1,")
		} else {
			write(p.completionFlatEventsBuffer, "0,")
		}

		if int(event.EventType) == int(enums.ETQuestionView) {
			write(p.completionFlatEventsBuffer, "1")
		} else {
			write(p.completionFlatEventsBuffer, "0")
		}
		writeEndLine(p.completionFlatEventsBuffer)

		p.completionFlatEvents++
	}

	if p.completionFlatEvents == p.config.DbConfiguration.WritingBatchSize {
		p.dbClient.BulkInsert(p.completionFlatEventsBuffer, "pure_survey.completion_stats_flat")
		p.completionFlatEvents = 0
	}
}

func (p *Proceeder) proceedCompletionFlatEvent(message []byte) {
	if p.completionEventsBuffer.Len() == 0 {
		p.completionEventsBuffer.Write(completionEventsHeader)
	}

	event := &pb.CompletionEvent{}
	err := proto.Unmarshal(message, event)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	timestamp := time.Unix(event.Timestamp, 0)

	writeDate(p.completionEventsBuffer, timestamp)
	writeSeparator(p.completionEventsBuffer)

	writeTimestamp(p.completionEventsBuffer, timestamp)
	writeSeparator(p.completionEventsBuffer)

	writeInt(p.completionEventsBuffer, timestamp.Hour())
	writeSeparator(p.completionEventsBuffer)

	writeInt(p.completionEventsBuffer, int(event.SurveyId))
	writeSeparator(p.completionEventsBuffer)

	writeInt(p.completionEventsBuffer, int(event.QuestionId))
	writeSeparator(p.completionEventsBuffer)

	p.getStringArray(event.OptionIds).WriteTo(p.completionEventsBuffer)
	writeSeparator(p.completionEventsBuffer)

	writeString(p.completionEventsBuffer, event.Geo)
	writeSeparator(p.completionEventsBuffer)

	writeString(p.completionEventsBuffer, event.Lang)
	writeSeparator(p.completionEventsBuffer)

	writeInt(p.completionEventsBuffer, int(event.Gender))
	writeSeparator(p.completionEventsBuffer)

	if int(event.EventType) == int(enums.ETQuestionAnswer) {
		write(p.completionEventsBuffer, "1,")
	} else {
		write(p.completionEventsBuffer, "0,")
	}

	if int(event.EventType) == int(enums.ETQuestionView) {
		write(p.completionEventsBuffer, "1")
	} else {
		write(p.completionEventsBuffer, "0")
	}
	writeEndLine(p.completionEventsBuffer)

	p.completionEvents++

	if p.completionEvents == p.config.DbConfiguration.WritingBatchSize {
		p.dbClient.BulkInsert(p.completionEventsBuffer, "pure_survey.completion_stats")
		p.completionEvents = 0
	}
}

func (p *Proceeder) getStringArray(arr []int32) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{'\'', '['})

	for _, el := range arr {
		buf.Write([]byte(strconv.Itoa(int(el))))
	}

	buf.WriteByte(']')
	buf.WriteByte('\'')
	return buf
}

func writeDate(buffer *bytes.Buffer, timestamp time.Time) {
	buffer.WriteByte('\'')
	buffer.Write([]byte(timestamp.Format("2006-01-02")))
	buffer.WriteByte('\'')
}

func writeTimestamp(buffer *bytes.Buffer, timestamp time.Time) {
	buffer.WriteByte('\'')
	buffer.Write([]byte(timestamp.Format("2006-01-02 15:04:05")))
	buffer.WriteByte('\'')
}

func writeInt(buffer *bytes.Buffer, val int) {
	buffer.Write([]byte(strconv.Itoa(val)))
}

func writeString(buffer *bytes.Buffer, val string) {
	buffer.WriteByte('\'')
	buffer.Write([]byte(val))
	buffer.WriteByte('\'')
}

func write(buffer *bytes.Buffer, val string) {
	buffer.Write([]byte(val))
}

func writeSeparator(buffer *bytes.Buffer) {
	buffer.WriteByte(',')
}

func writeEndLine(buffer *bytes.Buffer) {
	buffer.WriteByte('\n')
}
