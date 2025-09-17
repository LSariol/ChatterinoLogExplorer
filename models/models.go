package models

import (
	"fmt"
	"strings"
)

type FormData struct {
	Channel    string
	Terms      []string
	Duration   int
	User       string
	ExactMatch bool
}

func (f FormData) String() string {
	return fmt.Sprintf("Channel: %s, Terms: %s, Duration: %d, User: %s, ExactMatch: %t", f.Channel, strings.Join(f.Terms, ", "), f.Duration, f.User, f.ExactMatch)
}

type Record struct {
	Date    string
	User    string
	Message string
}
