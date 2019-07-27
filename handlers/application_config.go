package handlers

import "fmt"

type applicationConfig struct {
	EnableMysql   bool `yaml:"enable_mysql"`
	Notifications []struct {
		Name      string `yaml:"name"`
		Teams     string `yaml:"teams"`
		Pagerduty string `yaml:"pagerduty"`
	} `yaml:"notifications"`
}

func (applConfig *applicationConfig) listRoutes() ([]*routes, error) {
	var routeEntries []*routes
	for _, notify := range applConfig.Notifications {
		route := &routes{
			Identifier:  notify.Name,
			RouteType:   "teams",
			PostURL:     notify.Teams,
			Description: notify.Name,
		}
		routeEntries = append(routeEntries, route)
		route = &routes{
			Identifier:  notify.Name,
			RouteType:   "pagerduty",
			PostURL:     notify.Pagerduty,
			Description: notify.Name,
		}
		routeEntries = append(routeEntries, route)
	}
	if routeEntries == nil {
		return nil, fmt.Errorf("unable to parse the entries from application configs %v", nil)
	}
	return routeEntries, nil
}

func (applConfig *applicationConfig) getRoute(identifier string, routeType string) (*routes, error) {
	var route *routes
	for _, notify := range applConfig.Notifications {
		if notify.Name == identifier && routeType == "pagerduty" {
			route = &routes{
				Identifier:  notify.Name,
				RouteType:   "pagerduty",
				PostURL:     notify.Pagerduty,
				Description: notify.Name,
			}
		}
		if notify.Name == identifier && routeType == "teams" {
			route = &routes{
				Identifier:  notify.Name,
				RouteType:   "teams",
				PostURL:     notify.Teams,
				Description: notify.Name,
			}
		}
	}
	if route == nil {
		return nil, fmt.Errorf("unable to find entry for %s with type %s", identifier, routeType)
	}
	return route, nil
}
