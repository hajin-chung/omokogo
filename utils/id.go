package utils

import "github.com/nlepage/go-cuid2"

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
