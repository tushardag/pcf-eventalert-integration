package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//TeamsOutgoingMsg : type for Teams message post request
type teamsOutgoingMsg struct {
	Type             string            `json:"@type"`
	Context          string            `json:"@context"`
	ThemeColor       string            `json:"themeColor"`
	Title            string            `json:"title"`
	Summary          string            `json:"summary,omitempty"`
	Sections         []section         `json:"sections,omitempty"`
	PotentialActions []potentialAction `json:"potentialAction,omitempty"`
}

//Section : Specific activity section need to be posted to Teams
type section struct {
	ActivityTitle string `json:"activityTitle,omitempty"`
	Facts         []fact `json:"facts,omitempty"`
	Markdown      bool   `json:"markdown,omitempty"`
}

//Fact : Key and value attributes to appear on Teams message
type fact struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

//PotentialAction : Adding any additional action block to the Teams message
type potentialAction struct {
	Type    string   `json:"@type,omitempty"`
	Name    string   `json:"name,omitempty"`
	Targets []target `json:"targets,omitempty"`
}

//Target : targeted key-value for the suggested action in the Teams message
type target struct {
	OS  string `json:"os,omitempty"`
	URI string `json:"uri,omitempty"`
}

//CompileTeamsMessage : Parsing and mapping the fields to predefined Teams type.
func CompileTeamsMessage(eventAlert *EventAlert) teamsOutgoingMsg {
	return teamsOutgoingMsg{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: eventAlert.Metadata.StatusColor,
		Title:      eventAlert.Metadata.Status + ": " + eventAlert.Metadata.EventDescription,
		Summary:    eventAlert.Publisher,
		Sections: []section{
			section{
				ActivityTitle: eventAlert.Metadata.Foundation,
				Facts: []fact{
					fact{
						Name:  "Topic",
						Value: eventAlert.Topic,
					},
					fact{
						Name:  "Job",
						Value: eventAlert.Metadata.Job,
					},
					fact{
						Name:  "Value",
						Value: eventAlert.Metadata.Value,
					},
					fact{
						Name:  "Event Type",
						Value: eventAlert.Metadata.EventType,
					},
					fact{
						Name:  "Publisher",
						Value: eventAlert.Publisher,
					},
				},
				Markdown: false,
			},
		},
		PotentialActions: []potentialAction{
			potentialAction{
				Type: "OpenUri",
				Name: "View in HealthWatch",
				Targets: []target{
					target{
						OS:  "default",
						URI: eventAlert.Metadata.URL,
					},
				},
			},
			potentialAction{
				Type: "OpenUri",
				Name: "Refer Documentation",
				Targets: []target{
					target{
						OS:  "default",
						URI: eventAlert.Metadata.DocsURL,
					},
				},
			},
		},
	}
}

//PostMessage : Posting the message to MSTeams
func (msg teamsOutgoingMsg) PostMessage(endpoint string) error {
	enc, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(enc)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("error on message: %s", res.Status)
	}
	log.Printf("Successfully posted the message to MSTeam. Response code: %s", res.Status)
	return nil
}
