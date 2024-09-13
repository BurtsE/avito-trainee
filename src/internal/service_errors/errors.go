package service_errors

import "fmt"

type UserNotFound struct {
	Err error
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("user not found: %v", e.Err)
}

type AuthError struct {
	Err error
}

func (e AuthError) Error() string {
	return fmt.Sprintf("authentification error: %v", e.Err)
}

type TenderError struct {
	Err error
}

func (e TenderError) Error() string {
	return fmt.Sprintf("data validation error: %v", e.Err)
}
