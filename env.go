package main

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type EnvType struct {
	DB_URL        string
}

var Env EnvType

func LoadEnv() error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	content := string(data)
	for _, line := range strings.Split(content, "\n") {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]
		log.Printf("%s=%s", key, value)

		reflect.ValueOf(&Env).Elem().FieldByName(key).SetString(value)
	}

	return nil
}
