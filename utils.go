package main

import (
	"crypto/sha256"
	"github.com/nlepage/go-cuid2"
)

func Hash(raw string) string {
	h := sha256.New()

	h.Write([]byte(raw))
	bs := h.Sum(nil)

	return string(bs[:])
}

var idFunction func() (string, error)

func InitId() error {
	var err error
	idFunction, err = cuid2.Init(cuid2.Options{
		Length: 10,
	})

	return err
}

func CreateId() string {
	if idFunction == nil {
		InitId()
	}
	id, _ := idFunction()
	return id
}
