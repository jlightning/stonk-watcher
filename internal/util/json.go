package util

import "encoding/json"

func MustJSONStringify(data interface{}, indent bool) string {
	var bytes []byte
	var err error
	if indent {
		bytes, err = json.MarshalIndent(data, "", "  ")
	} else {
		bytes, err = json.Marshal(data)
	}
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
