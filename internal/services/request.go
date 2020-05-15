package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jz222/loggy/internal/keys"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/templates"
	"github.com/mailgun/mailgun-go/v4"
)

// InterfaceRequest represents the interface for the request service.
type InterfaceRequest interface {
	SendSlackAlert(models.Service, models.Error) error
	Post(payload interface{}, url string) error
	SendEmail(recipient, event string, data map[string]interface{}) error
}

// Request contains methods to send HTTP requests.
type Request struct{}

// SendSlackAlert sends a formatted error notification as Slack webhook.
func (r *Request) SendSlackAlert(service models.Service, errorEvent models.Error) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"mrkdwn_in":   []string{"text"},
				"color":       "#FF0055",
				"pretext":     fmt.Sprintf("An error occurred in %s", service.Name),
				"author_name": errorEvent.Type,
				"title":       errorEvent.Message,
				"title_link":  fmt.Sprintf("%s/services/%s/error/%s", keys.GetKeys().CLIENT_URL, service.ID.Hex(), errorEvent.ID.Hex()),
				"text":        "Visit your LOGGY dashboard for more details",
				"fields": []map[string]interface{}{
					{
						"title": "In Service",
						"value": service.Name,
						"short": true,
					},
					{
						"title": "Occurrences",
						"value": fmt.Sprintf("%d", errorEvent.Count),
						"short": true,
					},
					{
						"title": "Resolved",
						"value": strconv.FormatBool(errorEvent.Resolved),
						"short": true,
					},
					{
						"title": "Adapter",
						"value": fmt.Sprintf("%s %s", errorEvent.Adapter.Name, errorEvent.Adapter.Version),
						"short": true,
					},
				},
				"footer": "LOGGY",
				"ts":     errorEvent.Timestamp,
			},
		},
	})
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", service.SlackWebhookURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

// Post sends a given payload as a post request to a given URL.
func (r *Request) Post(payload interface{}, url string) error {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

// SendEmail sends an email to the given recipient.
func (r *Request) SendEmail(recipient, event string, data map[string]interface{}) error {
	var emailTemplate = ""
	var emailRawBody = ""
	var subject = ""

	// Load Mailgun settings
	mailgunPrivateKey := keys.GetKeys().MAILGUN_PRIVATE_KEY
	mailgunDomain := keys.GetKeys().MAILGUN_DOMAIN
	if mailgunPrivateKey == "" || mailgunDomain == "" {
		return nil
	}

	// Determine email template
	switch event {
	case "invitation":
		subject = "You were invited to LOGGY"
		emailTemplate = templates.Invitation
		emailRawBody = templates.InvitationRaw
	default:
		return errors.New("the provided event " + event + " is not available")
	}

	// Parse raw email template
	t := template.Must(template.New("email").Parse(emailRawBody))

	builder := &strings.Builder{}

	err := t.Execute(builder, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	parsedBody := builder.String()

	// Parse email template
	t = template.Must(template.New("email").Parse(emailTemplate))

	builder = &strings.Builder{}

	err = t.Execute(builder, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	parsedHTML := builder.String()

	// Setup Mailgun and send message
	mg := mailgun.NewMailgun(mailgunDomain, mailgunPrivateKey)

	message := mg.NewMessage("LOGGY <no-reply@loggy.io>", subject, parsedBody, recipient)

	message.SetHtml(parsedHTML)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
