package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tushardag/pcf-eventalert-integration/helpers"
)

const (
	pagerdutyType = "pagerduty"
)

//PagerDutyAlert ...
func (rh *RequestHandler) PagerDutyAlert(w http.ResponseWriter, r *http.Request) {
	//pulling mux variable
	vars := mux.Vars(r)
	route, err := rh.dbConn.getRoute(vars["identifier"], pagerdutyType)
	if err != nil {
		log.Printf("Unable to pull webhook URL for : " + vars["identifier"])
		// Write an error and stop the handler chain
		http.Error(w, "Unable to pull webhook URL for "+vars["identifier"]+". Please create the mapping or validate the identifier.", http.StatusPreconditionRequired)
		return
	}

	//Un-marshalling JSON through incoming request from Event Alert
	incomingMsg := new(helpers.EventAlert)
	if err := incomingMsg.ParseEventAlert(json.NewDecoder(r.Body)); err != nil {
		log.Printf("Error in parsing the request object.")
		log.Println(err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// fmt.Fprint(w, "Value of Person in function :", incomingMsg)
	fmt.Printf("%s EventAlert message received for: %s", incomingMsg.Metadata.Status, incomingMsg.Metadata.EventDescription)
	fmt.Println("Publishing message to PagerDuty " + vars["identifier"])
	//Building the message body to post a call for MSTeams webhook
	//Reference fields https://docs.microsoft.com/en-us/outlook/actionable-messages/card-reference
	err = helpers.CompilePagerDutyMessage(incomingMsg, route.PostURL).CreateIncident()
	if err != nil {
		log.Printf("Error in opening Incident in PagerDuty: %s\n", err)
		http.Error(w, "Unable to open Incident in PagerDuty.", http.StatusInternalServerError)
		return
	}
}
