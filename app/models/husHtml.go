package models

import (
	"github.com/revel/revel"
	"github.com/dustin/go-humanize"
	"html/template"
	"time"
	"strconv"
	"reflect"
	"encoding/json"
)

type HusHtml struct {
	*revel.Controller
}

func (husHtml HusHtml) Html(str string) template.HTML {
	return template.HTML(str)
}

func (husHtml HusHtml) Mod(a, b int) int {
	 return a%b
}

func (husHtml HusHtml) FormatFValue(number float64, precision int) string {
	format := "#,###."
	for i := 0; i < precision; i++ {
		format = format + "#"
	}

	return humanize.FormatFloat(format, number)
}

func (husHtml HusHtml) Add(a, b int) int64 {
	return int64(a) + int64(b)
}

func (husHtml HusHtml) Subtract(a, b int) int64 {
	return int64(a) - int64(b)
}

func (husHtml HusHtml) FormatCurrTime(key string) string {
	if (key == "year") {
		return strconv.Itoa(time.Now().Year())
	}
	if (key == "month") {
		return time.Now().Month().String()[:3]
	}
	if (key == "day") {
		return strconv.Itoa(time.Now().Day())
	}
	return ""
}

func (husHtml HusHtml) InArray(val interface{}, arrayString string) bool {
	if arrayString == "" {
		return false
	}
	var slice interface{}
	json.Unmarshal([]byte(arrayString), &slice)

	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}