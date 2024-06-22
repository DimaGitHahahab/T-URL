package service

import "net/url"

type URLValidator struct {
}

func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

func (v URLValidator) IsValidURL(long string) bool {
	_, err := url.ParseRequestURI(long)
	return err == nil
}
