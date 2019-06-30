package helpers

import (
	"encoding/json"
	"fmt"
)

//EventAlert ... the structure of Event Alert message published
type EventAlert struct {
	Publisher string `json:"publisher,omitempty"`
	Topic     string `json:"topic"`
	Metadata  struct {
		Status           string `json:"status"`
		StatusColor      string `json:"statusColor,omitempty"`
		Value            string `json:"value,omitempty"`
		Job              string `json:"job,omitempty"`
		Index            string `json:"index,omitempty"`
		IP               string `json:"ip,omitempty"`
		Deployment       string `json:"deployment,omitempty"`
		Foundation       string `json:"foundation,omitempty"`
		EventType        string `json:"eventType,omitempty"`
		EventDescription string `json:"eventDescription,omitempty"`
		URL              string `json:"url,omitempty"`
		DocsURL          string `json:"docsUrl,omitempty"`
	} `json:"metadata,omniemtpy"`
}

//ParseEventAlert ... Parsing and mapping the fields to predefined type.
func (eventAlert *EventAlert) ParseEventAlert(request *json.Decoder) error {
	if err := request.Decode(&eventAlert); err != nil {
		return err
	}
	//Verify mandatory fields
	if eventAlert.Topic == "" || eventAlert.Metadata.Status == "" || eventAlert.Metadata.EventDescription == "" {
		return fmt.Errorf("Missing mandatory fields i.e. Topic/Status/EventDescription")
	}
	return nil
}
