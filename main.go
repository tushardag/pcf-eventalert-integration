package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/tushardag/pcf-eventalert-integration/handlers"
)

func main() {
	fmt.Println("Starting up the pcf-eventalert-integration server...")
	fmt.Println("Validating Application config.")

	ymlFile, err := ioutil.ReadFile("application.yml")
	if err != nil {
		log.Println("Error in reading application config.")
		log.Fatalln(err)
	}

	// For the local instance of mysql, update the username,
	// password and instance connection string. When running locally,
	// localhost:3306 is used
	var mysqlDBconfig handlers.MySQLConfig
	mysqlDBconfig, err = configureMySQL()
	//fmt.Printf("Connecting to MySQL Host %s on port %d \n", mysqlDBconfig.Host, mysqlDBconfig.Port)
	requestHandler, err := handlers.RequestHandlerInit(mysqlDBconfig, ymlFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer requestHandler.CloseDB()

	fmt.Println("Initializing the webserver process.")
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
	// Supress the mapping management for non-db mode
	if requestHandler.DBinUse() {
		router.HandleFunc("/{type}/{identifier}", requestHandler.CreatMapping).Methods("PUT")
		router.HandleFunc("/{type}/{identifier}", requestHandler.RemoveMapping).Methods("DELETE")
	}
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
		fmt.Printf("Enabling route: %s for method %s\n", path, method)
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
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		fmt.Println("Shutting down the server process")
		os.Exit(0)
	}()

	//Starting the webserver
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Unable to start the pcf-eventalert-integration server.")
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

type mysqlDBInfo struct {
	Credentials mysqlDBCredentials `json:"credentials"`
}

type mysqlDBCredentials struct {
	Host     string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"username"`
	Password string `json:"password"`
}

func configureMySQL() (handlers.MySQLConfig, error) {
	if os.Getenv("VCAP_SERVICES") != "" {
		// Running in PCF.
		s := os.Getenv("VCAP_SERVICES")
		services := make(map[string][]mysqlDBInfo)
		err := json.Unmarshal([]byte(s), &services)
		if err != nil {
			log.Printf("Error parsing MySQL connection information: %v\n", err.Error())
			return handlers.MySQLConfig{}, err
		}
		info := services["p.mysql"]
		if len(info) == 0 {
			log.Printf("No MySQL databases are bound to this application.\n")
			return handlers.MySQLConfig{}, fmt.Errorf("unable to find service with name p.mysql")
		}
		// Assumes only a single MySQLDB is bound to this application
		creds := info[0].Credentials
		//fmt.Println(creds)
		return handlers.MySQLConfig{
			Username: creds.User,
			Password: creds.Password,
			Host:     creds.Host,
			Port:     creds.Port,
		}, nil
	}

	// Running locally.
	return handlers.MySQLConfig{
		Username: "mapper",
		Password: "mapper",
		Host:     "localhost",
		Port:     3306,
	}, nil
}
