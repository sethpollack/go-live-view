package handler

import (
	"encoding/base64"
	"encoding/json"
)

type tokenizer interface {
	Encode(any) (string, error)
	Decode(string, any) error
}

type defaultTokenizer struct{}

func (d *defaultTokenizer) Encode(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func (d *defaultTokenizer) Decode(s string, v any) error {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
