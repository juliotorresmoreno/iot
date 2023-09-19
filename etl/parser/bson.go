package parser

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
)

type BSONProvider struct{}

func NewBSONProvider() *BSONProvider {
	return &BSONProvider{}
}

func (el *BSONProvider) Encode(value interface{}) (string, error) {
	buff, err := bson.Marshal(value)
	if err != nil {
		return "", err
	}
	b64 := base64.RawStdEncoding.EncodeToString(buff)

	return b64, nil
}

func (el *BSONProvider) Decode(value string, obj interface{}) error {
	decoded, err := base64.RawStdEncoding.DecodeString(value)

	if err != nil {
		return err
	}

	return bson.Unmarshal(decoded, obj)
}
