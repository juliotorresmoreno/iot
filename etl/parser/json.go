package parser

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type JSONProvider struct{}

func NewJSONProvider() *JSONProvider {
	return &JSONProvider{}
}

func (el *JSONProvider) Encode(value interface{}) (string, error) {
	buff, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}

func (el *JSONProvider) Decode(value string, obj interface{}) error {
	return bson.Unmarshal([]byte(value), obj)
}
