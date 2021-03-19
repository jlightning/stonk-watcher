package util

import (
	"encoding/json"

	"github.com/tidwall/pretty"
)

func MustJSONStringify(data interface{}, indent bool) string {
	var bytes []byte
	var err error

	bytes, err = json.Marshal(data)
	if err != nil {
		panic(err)
	}
	if indent {
		bytes = pretty.Color(pretty.PrettyOptions(bytes, &pretty.Options{
			Width:  180,
			Indent: "  ",
		}), nil)
	}

	return string(bytes)
}
