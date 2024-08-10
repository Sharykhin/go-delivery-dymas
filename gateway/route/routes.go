package route

import (
	"fmt"
	"net/http"
	"os"

	"gateway/handler"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"

	"gopkg.in/yaml.v3"
)

type Routes struct {
	Routes []Route
}

type Route struct {
	Path         string   `yaml:"path"`
	Methods      []string `yaml:"methods"`
	HostRedirect string   `yaml:"host_redirect"`
	Parameters   []Parameter
}

type Parameter struct {
	Name     string   `yaml:"name"`
	Patterns []string `yaml:"patterns"`
}

// CreateServicesRoutesFromConfig create routes based on routes from config yml and add according handler  for route
func CreateServicesRoutesFromConfig(middleWares map[string]func(paramName string) func(next http.Handler) http.Handler) (map[string]pkghttp.Route, error) {
	var routes Routes
	yamlRoutesFile, err := os.ReadFile("routes.yaml")
	err = yaml.Unmarshal(yamlRoutesFile, &routes)
	if err != nil {
		return nil, err
	}

	fmt.Println(routes)
	var routers = make(map[string]pkghttp.Route)
	for _, route := range routes.Routes {
		router := pkghttp.Route{
			Methods: route.Methods,
			Handler: handler.NewGateWayProxyHandler(route.HostRedirect),
		}
		if route.Parameters != nil {
			for _, routeParameter := range route.Parameters {
				if routeParameter.Patterns != nil {
					for _, pattern := range routeParameter.Patterns {
						if routeParameter.Name != "" && pattern != "" {
							middleWare, ok := middleWares[pattern]
							if ok && middleWare != nil {

								router.Middlewares = append(router.Middlewares, middleWare(routeParameter.Name))
							}
						}
					}
				}
			}
		}
		routers[route.Path] = router
	}

	return routers, nil
}
