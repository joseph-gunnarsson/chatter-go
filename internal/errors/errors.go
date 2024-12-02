package errors

import (
	"errors"
	"fmt"
	"os"
)

type ServerError struct {
	Message string
}

func (e ServerError) Error() string {
	return e.Message
}

type ClientErrorExit struct {
	Message string
}

func (e ClientErrorExit) Error() string {
	return e.Message
}

type ClientErrorSilent struct {
	Message string
}

func (e ClientErrorSilent) Error() string {
	return e.Message
}

type ClientGUILoadError struct {
	Message string
}

func (e ClientGUILoadError) Error() string {
	return e.Message
}

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.As(err, &ServerError{}):
		return fmt.Errorf("server error: %w", err)
	case errors.As(err, &ClientErrorExit{}):
		fmt.Println(fmt.Errorf("client error: %w", err))
		os.Exit(1)
	case errors.As(err, &ClientGUILoadError{}):
		fmt.Println(fmt.Errorf("client GUI load error: %w", err))
		os.Exit(1)
	case errors.As(err, &ClientErrorSilent{}):
		return nil
	default:
		return fmt.Errorf("unexpected error: %w", err)
	}
	return nil
}
