package json

import ji "github.com/json-iterator/go"

// Json func alias.
var (
	json          = ji.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

// MarshalToString ...
func MarshalToString(v interface{}) string {
	s, err := ji.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}
