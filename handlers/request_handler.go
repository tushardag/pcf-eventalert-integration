package handlers

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

//RequestHandler : Application-wide configuration to allow passing already established DB interface
type RequestHandler struct {
	dbConn     *mysqlDB
	applConfig *applicationConfig
}

// Routes holds metadata about a route mapping records.
type routes struct {
	Identifier  string
	RouteType   string
	PostURL     string
	Description string
}

//RequestHandlerInit : Initializing the DB session
func RequestHandlerInit(config MySQLConfig, yamlFile []byte) (*RequestHandler, error) {
	var rh RequestHandler
	var err error
	err = yaml.Unmarshal(yamlFile, &rh.applConfig)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Teams name: " + rh.applConfig.Notifications[0].Name)
	if rh.applConfig.EnableMysql {
		fmt.Println("Establishing MySQL DB Connection")
		rh.dbConn, err = newDBConnection(config)
		if err != nil {
			log.Println("Unable to get DB connection")
			return nil, err
		}
		fmt.Println("Successfully established DB connection")
		//fmt.Println(dbConn)
	} else {
		fmt.Println("Application is being configured to run with NO MySQL DB instance.")
	}
	return &rh, nil
}

//CloseDB : Free up the DB resource
func (rh *RequestHandler) CloseDB() {
	if rh.applConfig.EnableMysql {
		fmt.Println("Closing the DB connection.")
		rh.dbConn.close()
	}
}

//DBinUse : Return if the MySQL DB is in use
func (rh *RequestHandler) DBinUse() bool {
	return rh.applConfig.EnableMysql
}
