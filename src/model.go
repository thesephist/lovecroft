package main

import (
	"fmt"
	"strings"
	"time"
)

type Subscriber struct {
	ID         string
	GivenName  string
	FamilyName string
	StartDate  time.Time
	EndDate    time.Time

	email      string
	unsubToken string
}

func (s Subscriber) Name() string {
	return s.GivenName + " " + s.FamilyName
}

func (s Subscriber) Email() string {
	return s.email
}

func (s Subscriber) UpdateEmail(address string) {
	if !strings.Contains(address, "@") {
		return
	}

	s.email = address
}

func (s Subscriber) UnsubPath() string {
	return fmt.Sprintf("/unsubscribe?token=%s", s.unsubToken)
}

func (s Subscriber) IsActive() bool {
	now := time.Now()
	return now.After(s.StartDate) && now.Before(s.EndDate)
}

type List struct {
	Name        string
	Subscribers []Subscriber
}

func (list List) ActiveSubscribers() (scribers []Subscriber) {
	for _, scriber := range list.Subscribers {
		if scriber.IsActive() {
			scribers = append(scribers, scriber)
		}
	}
	return
}

func (list List) RenderToCSV() string {
	rows := []string{
		"ID,Email,Given Name,Family Name",
	}
	for _, scriber := range list.ActiveSubscribers() {
		rows = append(rows, fmt.Sprintf("%s,\"%s\",\"%s\",\"%s\"",
			scriber.ID,
			scriber.GivenName,
			scriber.FamilyName,
			scriber.email,
		))
	}
	return strings.Join(rows, "\n")
}

type Directory struct {
	Lists []List
}

func (d Directory) FindList(name string) (List, error) {
	for _, list := range d.Lists {
		if list.Name == name {
			return list, nil
		}
	}
	return List{}, notFoundError{subject: "List"}
}
