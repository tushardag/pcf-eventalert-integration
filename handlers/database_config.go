package handlers

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net/url"

	// importing mysql driver in conjunction with database/sql package
	_ "github.com/go-sql-driver/mysql"
)

//createMappingTable : verify and create db and table for first time setup
var createMappingTable = []string{
	`CREATE DATABASE IF NOT EXISTS event_router_mapping DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE event_router_mapping;`,
	`CREATE TABLE IF NOT EXISTS route_mapping (
		identifier VARCHAR(30) NOT NULL,
		routeType VARCHAR(10) NOT NULL,
		postURL VARCHAR(255) NOT NULL,
		description TEXT NULL,
		PRIMARY KEY (identifier, routeType)
	)`,
}

//MysqlDB : persists event mapping to MySQL interface
type mysqlDB struct {
	conn *sql.DB

	fetchAll   *sql.Stmt
	retriveOne *sql.Stmt
	createNew  *sql.Stmt
	removeOne  *sql.Stmt
}

//mappingObject : Ensure mysqlDB conforms to the interface.
var _ mappingDatabase = &mysqlDB{}

// dbConnectionString : Returns a connection string suitable for sql.Open
func (config *MySQLConfig) dbConnectionString(schemaName string) string {
	var returnString string
	if config.Username != "" {
		returnString = config.Username
		if config.Password != "" {
			returnString = returnString + ":" + config.Password
		}
		returnString = returnString + "@"
	}
	return fmt.Sprintf("%stcp([%s]:%d)/%s", returnString, config.Host, config.Port, schemaName)
}

//NewDBConnection : Initiating new DB connection instance
func newDBConnection(config MySQLConfig) (*mysqlDB, error) {
	// Check database and table exists. If not, create it.
	if err := config.ensureTableExists("event_router_mapping"); err != nil {
		return nil, err
	}
	var databaseConn mysqlDB
	var err error
	databaseConn.conn, err = sql.Open("mysql", config.dbConnectionString("event_router_mapping"))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	if err := databaseConn.conn.Ping(); err != nil {
		databaseConn.conn.Close()
		return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
	}

	// Prepared statements. The actual SQL queries are in the code near the relevant method (e.g. ListRoutes).
	if databaseConn.fetchAll, err = databaseConn.conn.Prepare(listStatement); err != nil {
		log.Println("Failed to prepare list statement")
		return nil, fmt.Errorf("mysql: prepare list: %v", err)
	}
	if databaseConn.createNew, err = databaseConn.conn.Prepare(insertStatement); err != nil {
		log.Println("Failed to prepare insert statement")
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}
	if databaseConn.retriveOne, err = databaseConn.conn.Prepare(getStatement); err != nil {
		log.Println("Failed to prepare get statement")
		return nil, fmt.Errorf("mysql: prepare get: %v", err)
	}
	if databaseConn.removeOne, err = databaseConn.conn.Prepare(deleteStatement); err != nil {
		log.Println("Failed to prepare delete statement")
		return nil, fmt.Errorf("mysql: prepare delete: %v", err)
	}
	fmt.Println("Returning the DB instance")
	return &databaseConn, nil
}

func (config MySQLConfig) ensureTableExists(dbName string) error {
	conn, err := sql.Open("mysql", config.dbConnectionString(""))
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
			"could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE " + dbName); err != nil {
		fmt.Println("Creating event_router_mapping DB and route_mapping Table")
		return createTable(conn)
	}

	if _, err := conn.Exec("DESCRIBE route_mapping"); err != nil {
		fmt.Println("Found event_router_mapping DB. Creating route_mapping Table")
		return createTable(conn)
	}
	return nil
}

// Close closes the database, freeing up any resources.
func (db *mysqlDB) close() {
	db.conn.Close()
}

func urlToString(dbURL *url.URL) string {
	return fmt.Sprintf(
		"%v@tcp(%v)%v?parseTime=true", dbURL.User, dbURL.Host, dbURL.Path,
	)
}

// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB) error {
	for _, stmt := range createMappingTable {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// rowScanner is implemented by sql.Row and sql.Rows
type rowScanner interface {
	Scan(dest ...interface{}) error
}

// scanRoute reads a book from a sql.Row or sql.Rows
func scanRoute(s rowScanner) (*routes, error) {
	var (
		identifier  sql.NullString
		routeType   sql.NullString
		postURL     sql.NullString
		description sql.NullString
	)
	if err := s.Scan(&identifier, &routeType, &postURL, &description); err != nil {
		return nil, err
	}

	route := &routes{
		Identifier:  identifier.String,
		RouteType:   routeType.String,
		PostURL:     postURL.String,
		Description: description.String,
	}
	return route, nil
}

//
const listStatement = `SELECT * FROM route_mapping`

// ListRoutes returns a list of mapping records
func (db *mysqlDB) listRoutes() ([]*routes, error) {
	rows, err := db.fetchAll.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routeEntries []*routes
	for rows.Next() {
		route, err := scanRoute(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		routeEntries = append(routeEntries, route)
	}

	return routeEntries, nil
}

const getStatement = "SELECT * FROM route_mapping WHERE identifier = ? and routeType = ?"

// GetRoute retrieves a Route by its identifier.
func (db *mysqlDB) getRoute(identifier string, routeType string) (*routes, error) {
	book, err := scanRoute(db.retriveOne.QueryRow(identifier, routeType))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mysql: could not find route with identifier %s of type %s", identifier, routeType)
	}
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get route: %v", err)
	}
	return book, nil
}

const insertStatement = `
  INSERT INTO route_mapping (
	  identifier, routeType, postURL, description) 
	  VALUES (?, ?, ?, ?)`

// AddRoute saves a new Route mapping.
func (db *mysqlDB) addRoute(rt *routes) error {
	_, err := execAffectingOneRow(db.createNew, rt.Identifier, rt.RouteType, rt.PostURL, rt.Description)
	if err != nil {
		return err
	}
	return nil
}

const deleteStatement = `DELETE FROM route_mapping WHERE identifier = ? and routeType = ?`

// DeleteRoute : removes a given route by its identifier.
func (db *mysqlDB) deleteRoute(identifier string, routeType string) error {
	_, err := execAffectingOneRow(db.removeOne, identifier, routeType)
	return err
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}
	return r, nil
}
