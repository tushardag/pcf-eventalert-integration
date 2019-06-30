package handlers

import "log"

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
	log.Println("Successfully established DB connection")
	//log.Println(dbConn)
	return &rh, nil
}

//CloseDB : Free up the DB resource
func (rh *RequestHandler) CloseDB() {
	log.Println("Closing the DB connection.")
	rh.dbConn.close()
}
