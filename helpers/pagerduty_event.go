package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const pagerDutyURL = "https://events.pagerduty.com/v2/enqueue"

//PDOutgoingMsg : type for PagerDuty message post request
type pdOutgoingMsg struct {
	Payload     payload `json:"payload"`
	RoutingKey  string  `json:"routing_key"`
	Links       []link  `json:"links"`
	EventAction string  `json:"event_action"`
	Client      string  `json:"client"`
	ClientURL   string  `json:"client_url"`
}

type payload struct {
	Summary       string       `json:"summary"`
	Source        string       `json:"source"`
	Severity      string       `json:"severity"`
	Component     string       `json:"component"`
	Group         string       `json:"group"`
	Class         string       `json:"class"`
	CustomDetails customDetail `json:"custom_details"`
}

type customDetail struct {
	Value string `json:"value"`
	IP    string `json:"ip"`
	Index string `json:"index"`
}

type link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

//CompilePagerDutyMessage : Parsing and mapping the fields to predefined PagerDuty type.
func CompilePagerDutyMessage(eventAlert *EventAlert, routingKey string) pdOutgoingMsg {
	return pdOutgoingMsg{
		RoutingKey:  routingKey,
		EventAction: "trigger",
		Payload: payload{
			Summary:   eventAlert.Metadata.Foundation + ": " + eventAlert.Metadata.EventDescription,
			Source:    eventAlert.Metadata.Foundation,
			Severity:  strings.ToLower(eventAlert.Metadata.Status),
			Component: eventAlert.Topic,
			Group:     eventAlert.Metadata.Job,
			Class:     eventAlert.Metadata.EventType,
			CustomDetails: customDetail{
				Value: eventAlert.Metadata.Value,
				IP:    eventAlert.Metadata.IP,
				Index: eventAlert.Metadata.Index,
			},
		},
		Links: []link{
			link{
				Href: "Refer Documentation",
				Text: eventAlert.Metadata.DocsURL,
			},
		},
		Client:    eventAlert.Publisher,
		ClientURL: eventAlert.Metadata.URL,
	}
}

//CreateIncident : Posting the message to MSTeams
func (pd pdOutgoingMsg) CreateIncident() error {
	enc, err := json.Marshal(pd)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(enc)
	res, err := http.Post(pagerDutyURL, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		log.Printf("PD Request: %s", enc)
		return fmt.Errorf("error in posting incident to PD: %s", res.Status)
	}
	fmt.Printf("Successfully posted the message to PagerDuty. Response code: %s", res.Status)
	return nil
}
