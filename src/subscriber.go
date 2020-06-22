package main

import (
	"fmt"
	"strings"
)

type Subscriber struct {
	Id string

	givenName  string
	familyName string
	email      string
	unsubToken string
}

func (s Subscriber) Name() {
	return s.givenName + " " + s.familyName
}

func (s Subscriber) Email() {
	return s.email
}

func (s Subscriber) UpdateEmail(address string) {
	if !strings.Contains(address, "@") {
		return
	}

	s.email = address
}

func (s Subscriber) UnsubPath() {
	return fmt.Sprintf("/unsubscribe?token=%s", s.unsubToken)
}

type Subscription struct {
	Recipient Subscriber
	StartDate int
	EndDate   int
}

type List struct {
	Name          string
	Subscriptions []Subscription
}

type Directory struct {
	Lists []List
}
