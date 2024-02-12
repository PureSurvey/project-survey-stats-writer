package events

import (
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
)

var (
	ParserPool = &fastjson.ParserPool{}
)

func GetEventsInsertStatement(messages []string) string {
	sb := strings.Builder{}
	sb.WriteString(`INSERT INTO public.SurveyEvents (Date, Hour, UnitId, Platform, IsMobile, SurveyStart, SurveyExit) VALUES `)

	parser := ParserPool.Get()
	defer ParserPool.Put(parser)

	for i := 0; i < len(messages); i++ {
		val, _ := parser.Parse(messages[i])

		sb.WriteString("('")
		sb.WriteString(string(val.GetStringBytes(`date`)))
		sb.WriteString("',")
		sb.WriteString(strconv.Itoa(val.GetInt(`hour`)))
		sb.WriteString(",")
		sb.WriteString(strconv.Itoa(val.GetInt(`unitId`)))
		sb.WriteString(",'")
		sb.WriteString(string(val.GetStringBytes(`platform`)))
		sb.WriteString("',")
		sb.WriteString(strconv.Itoa(val.GetInt(`isMobile`)))
		sb.WriteString(",")
		eventType := val.GetInt(`eventType`)
		if eventType == 0 {
			sb.WriteString("1,0")
		} else if eventType == 1 {
			sb.WriteString("0,1")
		}
		sb.WriteString(")")

		if i != len(messages)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
