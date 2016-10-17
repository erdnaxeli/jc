package jc

import "fmt"

type ConfigError struct {
	Field string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Missing field '%s' in configuration", e.Field)
}

type ConnectionError struct {
	Err error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("Connection error: %s", e.Err)
}
