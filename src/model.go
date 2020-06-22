package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscriber struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
	StartDate  time.Time
	EndDate    time.Time
	UnsubToken string
}

func (s Subscriber) Name() string {
	return s.GivenName + " " + s.FamilyName
}

func (s Subscriber) IsActive() bool {
	now := time.Now()
	return now.After(s.StartDate) && (s.EndDate.IsZero() || now.Before(s.EndDate))
}

type List struct {
	Name        string
	Subscribers []Subscriber
}

func newUnsubToken() string {
	return strings.ToLower(uuid.New().String())
}

func (list *List) Subscribe(scriber Subscriber) {
	scriber.StartDate = time.Now()
	scriber.EndDate = time.Time{}
	scriber.UnsubToken = newUnsubToken()

	for i, existing := range list.Subscribers {
		if existing.Email == scriber.Email {
			list.Subscribers[i] = scriber
			return
		}
	}

	list.Subscribers = append(list.Subscribers, scriber)
}

func (list *List) Unsubscribe(token string) error {
	for i, existing := range list.Subscribers {
		if existing.UnsubToken == token {
			list.Subscribers[i].EndDate = time.Now()
			list.Subscribers[i].UnsubToken = newUnsubToken()
			return nil
		}
	}

	return notFoundError{subject: "Subscriber"}
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
	rows := []string{}
	for _, scriber := range list.Subscribers {
		items := []string{
			scriber.Email,
			scriber.GivenName,
			scriber.FamilyName,
			scriber.StartDate.Format(time.RFC3339),
			scriber.EndDate.Format(time.RFC3339),
			scriber.UnsubToken,
		}
		for i, it := range items {
			items[i] = "\"" + it + "\""
		}
		rows = append(rows, strings.Join(items, ","))
	}
	return strings.Join(rows, "\n")
}

type Directory struct {
	Lists []List
}

func (d *Directory) FindList(name string) (*List, error) {
	for i, list := range d.Lists {
		if list.Name == name {
			return &d.Lists[i], nil
		}
	}
	return nil, notFoundError{subject: "List"}
}
