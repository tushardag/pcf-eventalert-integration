package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tushardag/pcf-eventalert-integration/helpers"
)

const (
	teamsType = "teams"
)

//MSTeamsAlert : Interface with MS Team and publish the alert
func (rh *RequestHandler) MSTeamsAlert(w http.ResponseWriter, r *http.Request) {
	//pulling mux variable
	vars := mux.Vars(r)
	var route *routes
	var err error
	if rh.applConfig.EnableMysql {
		route, err = rh.dbConn.getRoute(vars["identifier"], teamsType)
	} else {
		route, err = rh.applConfig.getRoute(vars["identifier"], teamsType)
	}
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
	fmt.Printf("EventAlert message received for: %s", incomingMsg.Metadata.EventDescription)

	//Building the message body to post a call for MSTeams webhook
	//Reference fields https://docs.microsoft.com/en-us/outlook/actionable-messages/card-reference
	fmt.Println("Publishing message to Teams " + vars["identifier"] + " with URL - " + route.PostURL)

	// if err := helpers.CompileTeamsMessage(incomingMsg).PostMessage(route.PostURL); err != nil {
	// 	log.Printf("Error in publishing message to Teams: %s\n", err)
	// 	http.Error(w, "Unable to publish message to Teams", http.StatusInternalServerError)
	// 	return
	// }
}

//BuildMessage ... building the message based on the incoming msg fields
func BuildMessage(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
