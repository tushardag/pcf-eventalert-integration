package handlers

import (
	"fmt"
	"log"
)

//RequestHandler : Application-wide configuration to allow passing already established DB interface
type RequestHandler struct {
	dbConn *mysqlDB
}

//RequestHandlerInit : Initializing the DB session
func RequestHandlerInit(config MySQLConfig) (*RequestHandler, error) {
	var rh RequestHandler
	var err error
	rh.dbConn, err = newDBConnection(config)
	if err != nil {
		log.Println("Unable to get DB connection")
		return nil, err
	}
	fmt.Println("Successfully established DB connection")
	//fmt.Println(dbConn)
	return &rh, nil
}

//CloseDB : Free up the DB resource
func (rh *RequestHandler) CloseDB() {
	fmt.Println("Closing the DB connection.")
	rh.dbConn.close()
}
