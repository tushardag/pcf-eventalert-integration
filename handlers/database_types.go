package handlers

//MySQLConfig connection construct information for MySQL DB Connection
type MySQLConfig struct {
	// Optional.
	Username, Password string
	// Host of the MySQL instance.
	Host string
	// Port of the MySQL instance.
	Port int
}

// MappingDatabase provides thread-safe access to a database of mapping records.
type mappingDatabase interface {
	// ListRoutes returns a list of all available route mapping
	listRoutes() ([]*routes, error)

	// GetRoute retrieves a route by its identifier and type.
	getRoute(string, string) (*routes, error)

	// AddRoute saves a new route
	addRoute(rt *routes) error

	// DeleteRoute removes a given route by its identifier and type.
	deleteRoute(string, string) error

	// close closes the database, freeing up any available resources.
	// TODO: close() should return an error.
	close()
}
