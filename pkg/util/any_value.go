package _util

import (
	"encoding/json"
	"reflect"

	"github.com/gookit/goutil/strutil"
)

// AnyValueT ...
type AnyValueT struct {
	Raw    interface{}
	Text   string
	LineNo int
}

// AnyValue ...
type AnyValue = *AnyValueT

// Format ...
func (i AnyValue) String() string {
	return i.Text
}

// AnyValueFromRaw ...
func AnyValueFromRaw(ctx JogContext, lineNo int, raw interface{}, replace map[string]string) AnyValue {
	var text string

	alreadyNormalized := false

	if raw == nil {
		return &AnyValueT{Raw: nil, Text: "", LineNo: lineNo}
	}

	kind := reflect.TypeOf(raw).Kind()
	if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Array {
		json, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			// TODO: should we panic with this err ?
			ctx.LogError("failed to json format", "lineNo", lineNo, "raw", raw)
		} else {
			text = string(json)
		}
		alreadyNormalized = true
	} else {
		text = strutil.MustString(raw)
	}

	/* TODO: why?
	if len(text) >= 1 {
		if text[:1] == "\"" || text[:1] == "'" {
			text = text[1:]
		}
	}
	if len(text) >= 1 {
		if text[len(text)-1:] == "\"" || text[len(text)-1:] == "'" {
			text = text[:len(text)-1]
		}
	}
	*/
	text = strutil.Replaces(text, replace)

	if alreadyNormalized == false {
		var obj interface{}
		if err := json.Unmarshal([]byte(text), &obj); err == nil {
			json, err := json.MarshalIndent(obj, "", "  ")
			if err == nil {
				text = string(json)
			}
		}
		// TODO: should we panic with above err ?
	}

	return &AnyValueT{Raw: raw, Text: text, LineNo: lineNo}
}
