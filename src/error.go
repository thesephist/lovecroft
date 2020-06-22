package main

import "fmt"

type notFoundError struct {
	subject string
}

func (err notFoundError) Error() string {
	return fmt.Sprintf("%s not found", err.subject)
}

func IsNotFound(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}
