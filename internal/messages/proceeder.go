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
	trackingEventsHeader   = []byte("date,timestamp,hour,unit_id,geo,lang,gender,unit_requests,unit_responses,unit_views,unit_ends,survey_id,survey_responses,mismatched_survey_responses,mismatched_survey_reason,survey_starts,survey_ends\n")
	completionEventsHeader = []byte("date,timestamp,hour,survey_id,question_id,option_ids,geo,lang,gender,answered,views\n")
)

type Proceeder struct {
	dbClient db.Client
	messages chan model.Message
	config   *configuration.AppConfiguration

	trackingEvents   int
	completionEvents int

	trackingEventsBuffer   *bytes.Buffer
	completionEventsBuffer *bytes.Buffer
}

func NewProceeder(dbClient db.Client, config *configuration.AppConfiguration, messages chan model.Message) Proceeder {
	return Proceeder{dbClient: dbClient, config: config, messages: messages, trackingEventsBuffer: bytes.NewBuffer([]byte{}), completionEventsBuffer: bytes.NewBuffer([]byte{})}
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

	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.Write([]byte(timestamp.Format("2006-01-02")))
	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.Write([]byte(timestamp.Format("2006-01-02 15:04:05")))
	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.Write([]byte(strconv.Itoa(timestamp.Hour())))
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.Write([]byte(strconv.Itoa(int(event.SurveyId))))
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.Write([]byte(strconv.Itoa(int(event.QuestionId))))
	p.completionEventsBuffer.WriteByte(',')

	p.getStringArray(event.OptionIds).WriteTo(p.completionEventsBuffer)
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.Write([]byte(event.Geo))
	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.Write([]byte(event.Lang))
	p.completionEventsBuffer.WriteByte('\'')
	p.completionEventsBuffer.WriteByte(',')

	p.completionEventsBuffer.Write([]byte(strconv.Itoa(int(event.Gender))))
	p.completionEventsBuffer.WriteByte(',')

	if int(event.EventType) == int(enums.ETQuestionAnswer) {
		p.completionEventsBuffer.Write([]byte("1,"))
	} else {
		p.completionEventsBuffer.Write([]byte("0,"))
	}

	if int(event.EventType) == int(enums.ETQuestionView) {
		p.completionEventsBuffer.Write([]byte("1"))
	} else {
		p.completionEventsBuffer.Write([]byte("0"))
	}
	p.completionEventsBuffer.WriteByte('\n')

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
