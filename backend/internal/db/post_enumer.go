// Code generated by "enumer --type=Post --json"; DO NOT EDIT.

package db

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _PostName = "UnknownTeamPlayer"

var _PostIndex = [...]uint8{0, 7, 11, 17}

const _PostLowerName = "unknownteamplayer"

func (i Post) String() string {
	if i >= Post(len(_PostIndex)-1) {
		return fmt.Sprintf("Post(%d)", i)
	}
	return _PostName[_PostIndex[i]:_PostIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PostNoOp() {
	var x [1]struct{}
	_ = x[Unknown-(0)]
	_ = x[Team-(1)]
	_ = x[Player-(2)]
}

var _PostValues = []Post{Unknown, Team, Player}

var _PostNameToValueMap = map[string]Post{
	_PostName[0:7]:        Unknown,
	_PostLowerName[0:7]:   Unknown,
	_PostName[7:11]:       Team,
	_PostLowerName[7:11]:  Team,
	_PostName[11:17]:      Player,
	_PostLowerName[11:17]: Player,
}

var _PostNames = []string{
	_PostName[0:7],
	_PostName[7:11],
	_PostName[11:17],
}

// PostString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PostString(s string) (Post, error) {
	if val, ok := _PostNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PostNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Post values", s)
}

// PostValues returns all values of the enum
func PostValues() []Post {
	return _PostValues
}

// PostStrings returns a slice of all String values of the enum
func PostStrings() []string {
	strs := make([]string, len(_PostNames))
	copy(strs, _PostNames)
	return strs
}

// IsAPost returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Post) IsAPost() bool {
	for _, v := range _PostValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Post
func (i Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Post
func (i *Post) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Post should be a string, got %s", data)
	}

	var err error
	*i, err = PostString(s)
	return err
}
