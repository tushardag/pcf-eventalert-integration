package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/tushardag/webhook-handler/handlers"
)

func main() {
	log.Printf("Starting up the webhook-handler server...")
	log.Printf("Connecting to the EventRouteMapping DB and initializing the base")
	// For the local instance of mysql, update the username,
	// password and instance connection string. When running locally,
	// localhost:3306 is used
	requestHandler, err := handlers.RequestHandlerInit(handlers.MySQLConfig{
		Username: "mapper",
		Password: "mapper",
		Host:     "localhost",
		Port:     3306,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer requestHandler.CloseDB()
	log.Printf("Initializing the webserver process.")
	router := mux.NewRouter()

	// list out the routes and usages information
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Print the routes and help information for app usages
		fmt.Fprintln(w, "{type} ==> teams or pagerduty")
		fmt.Fprintln(w, "{identifier} ==> unique tag for respective teams/pagerduty endpoint")
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				fmt.Fprintf(w, "ROUTE \"%s\" is servicing ", pathTemplate)
			}
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Fprintf(w, "on HTTP method %s", strings.Join(methods, ","))
			}
			fmt.Fprintln(w)
			return nil
		})
	}).Methods("GET")

	// Fetch the list of existing route mappings from DB in JSON format
	router.HandleFunc("/routes", requestHandler.ListMappings).Methods("GET")
	// Teams routes
	router.HandleFunc("/{type}/{identifier}", requestHandler.CreatMapping).Methods("PUT")
	router.HandleFunc("/{type}/{identifier}", requestHandler.RemoveMapping).Methods("DELETE")

	//MS Teams Event routing
	router.HandleFunc("/teams/{identifier}", requestHandler.MSTeamsAlert).Methods("POST")
	//PagerDuty Event routing
	router.HandleFunc("/pagerduty/{identifier}", requestHandler.PagerDutyAlert).Methods("POST")

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		method, err := route.GetMethods()
		if err != nil {
			return err
		}
		log.Printf("Enabling route: %s for method %s\n", path, method)
		return nil
	})

	//Configuring out HTTP server
	server := &http.Server{
		Addr: "0.0.0.0:" + getPort(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Handling gracefull shutdown of the server
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		log.Println("Shutting down the server process")
		os.Exit(0)
	}()

	//Starting the webserver
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Unable to start the webhook-handler server.")
		log.Fatalln(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}
