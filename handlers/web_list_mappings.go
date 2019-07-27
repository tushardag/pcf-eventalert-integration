package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

//ListMappings : Default landing to provide list of existing mappings and sample requests
func (rh *RequestHandler) ListMappings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var routes []*routes
	var err error
	if rh.applConfig.EnableMysql {
		routes, err = rh.dbConn.listRoutes()
		if err != nil {
			log.Printf("Unable to fetch the list of route mapping. %s\n", err)
			http.Error(w, "Unable to fetch route mapping from DB", http.StatusInternalServerError)
			return
		}
	} else {
		routes, err = rh.applConfig.listRoutes()
		if err != nil {
			log.Printf("Unable to fetch the list of route mapping. %s\n", err)
			http.Error(w, "Unable to fetch route mapping from DB", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(routes)
}
