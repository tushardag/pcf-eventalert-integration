package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type requestJSON struct {
	URL         string `json:"URL,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	supportedTypes = "teams,pagerduty"
)

//CreatMapping PUT request to create new mapping into the route_mapping table
func (rh *RequestHandler) CreatMapping(wr http.ResponseWriter, req *http.Request) {

	fmt.Printf("Received a PUT request. Creating new route mapping entry.")
	vars := mux.Vars(req)
	if !strings.Contains(supportedTypes, vars["type"]) {
		log.Printf("Invalid Entry type received. Type received: " + vars["type"])
		// Write an error and stop the handler chain
		http.Error(wr, "Not a valid Type.", http.StatusNotAcceptable)
		return
	}
	decoder := json.NewDecoder(req.Body)
	var reqJSON requestJSON
	if err := decoder.Decode(&reqJSON); err != nil {
		log.Printf("Invalid PUT Request")
		http.Error(wr, "Invalid JSON Request. Please verify and resubmit", http.StatusNotAcceptable)
		return
	}
	route := &routes{
		Identifier:  vars["identifier"],
		RouteType:   vars["type"],
		PostURL:     reqJSON.URL,
		Description: reqJSON.Description,
	}

	if route.RouteType == "teams" {
		_, err := url.ParseRequestURI(route.PostURL)
		if err != nil {
			log.Printf("Invalid URL received in PUT Request for Teams")
			http.Error(wr, "Invalid URL received for Teams. Please verify and resubmit", http.StatusNotAcceptable)
			return
		}
	}

	if err := rh.dbConn.addRoute(route); err != nil {
		log.Printf("Unable to add given route mapping Identifier:" + route.Identifier + " Type:" + route.RouteType)
		log.Println(err)
		http.Error(wr, "Internal server error. Please check the logs for more information", http.StatusInternalServerError)
		return
	}
	//wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	fmt.Printf("Successfully added a new mapping entry for " + route.Identifier + " with type as " + route.RouteType)
	return
}

//RemoveMapping DELETE request to remove the mapping from route_mapping table
func (rh *RequestHandler) RemoveMapping(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if !strings.Contains(supportedTypes, vars["type"]) {
		log.Printf("Invalid Entry type received for removal: " + vars["type"])
		// Write an error and stop the handler chain
		http.Error(wr, "Not a valid Type in the request.", http.StatusNotAcceptable)
		return
	}

	if err := rh.dbConn.deleteRoute(vars["identifier"], vars["type"]); err != nil {
		log.Printf("Unable to remove route mapping Identifier:" + vars["identifier"] + " Type:" + vars["type"])
		log.Println(err)
		http.Error(wr, "Internal server error. Please check the logs for more information", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully removed " + vars["type"] + " mapping for identifier " + vars["identifier"])
	return
}
