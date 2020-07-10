package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/yuin/goldmark"
)

var apiKey = os.Getenv("MAILGUN_APIKEY")

type Issue struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (issue Issue) BodyHTML() template.HTML {
	buf := bytes.Buffer{}
	err := goldmark.Convert([]byte(issue.Body), &buf)
	if err != nil {
		return template.HTML(issue.Body)
	}

	return template.HTML(buf.String())
}

type issueTemplateValues struct {
	List       List
	Subscriber Subscriber
	Issue      Issue
}

func logSendErr(issue Issue, scriber Subscriber, err error) {
	log.Printf("Error sending issue \"%s\" to %s: %v\n", issue.Subject, scriber.Email, err)
}

func (list List) Send(issue Issue) {
	tmpl, err := useTemplate("issue")
	if err != nil {
		log.Printf("Error sending issue \"%s\": %v", issue.Subject, err)
	}

	mg := mailgun.NewMailgun(list.Name, apiKey)

	vals := issueTemplateValues{}
	for _, scriber := range list.ActiveSubscribers() {
		vals = issueTemplateValues{
			List:       list,
			Subscriber: scriber,
			Issue:      issue,
		}
		buf := bytes.Buffer{}
		err := tmpl.Execute(&buf, vals)
		if err != nil {
			logSendErr(issue, scriber, err)
			continue
		}

		mail := mg.NewMessage(
			issue.From,
			issue.Subject,
			"",
			scriber.Email,
		)
		mail.SetHtml(buf.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		if apiKey == "" {
			log.Printf("No API Key found, mock-sending issue \"%s\" to %s\n", issue.Subject, scriber.Email)
			continue
		}

		resp, id, err := mg.Send(ctx, mail)
		if err != nil {
			logSendErr(issue, scriber, err)
		} else {
			log.Printf("Sent issue \"%s\" to %s: resp(%s), id(%s)\n",
				issue.Subject, scriber.Email, resp, id)
		}

		// reasonable rate limit -- this is 120 emails per minute.
		// which is enough for me but not too high to be dangerous
		time.Sleep(500 * time.Millisecond)
	}
}
