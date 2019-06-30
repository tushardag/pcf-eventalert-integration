package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

//Index : Default landing to provide list of existing mappings and sample requests
func (rh *RequestHandler) ListMappings(w http.ResponseWriter, r *http.Request) {
	routes, err := rh.dbConn.listRoutes()
	if err != nil {
		log.Printf("Unable to fetch the list of route mapping. %s\n", err)
		http.Error(w, "Unable to fetch route mapping from DB", http.StatusInternalServerError)
	}
	//routesJSON, err := json.Marshal(routes)
	// if err != nil {
	// 	log.Fatal("Cannot encode to JSON ", err)
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(routes)
}
