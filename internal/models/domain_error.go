package models

import "fmt"

type DomainError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Why         string `json:"why"`
	NextStep    string `json:"next_step"`
	Recoverable bool   `json:"recoverable"`
}

func (e DomainError) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return e.Message
}

func NewDomainError(code, message, why, nextStep string, recoverable bool) DomainError {
	return DomainError{
		Code:        code,
		Message:     message,
		Why:         why,
		NextStep:    nextStep,
		Recoverable: recoverable,
	}
}

func WrapDomainError(code, message, why, nextStep string, recoverable bool, err error) DomainError {
	if err == nil {
		return NewDomainError(code, message, why, nextStep, recoverable)
	}
	if why == "" {
		why = err.Error()
	} else {
		why = fmt.Sprintf("%s: %s", why, err.Error())
	}
	return NewDomainError(code, message, why, nextStep, recoverable)
}
