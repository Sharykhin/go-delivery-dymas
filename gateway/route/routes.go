package route

import (
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Sharykhin/go-delivery-dymas/gateway/handler"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
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
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
}

// CreateServicesRoutesFromConfig create routes based on routes from config yml and add according handler  for route
func CreateServicesRoutesFromConfig(middleWares map[string]func(paramName string) func(next http.Handler) http.Handler) (map[string]pkghttp.Route, error) {
	var routes Routes
	yamlRoutesFile, err := os.ReadFile("routes.yaml")
	err = yaml.Unmarshal(yamlRoutesFile, &routes)
	if err != nil {
		return nil, err
	}
	var routers = make(map[string]pkghttp.Route)
	for _, route := range routes.Routes {
		router := pkghttp.Route{
			Methods: route.Methods,
			Handler: handler.NewGateWayProxyHandler(route.HostRedirect),
		}

		for _, routeParameter := range route.Parameters {
			if routeParameter.Pattern != "" && routeParameter.Name != "" {
				middleWare, ok := middleWares[routeParameter.Pattern]
				if ok && middleWare != nil {
					router.Middlewares = append(router.Middlewares, middleWare(routeParameter.Name))
				}
			}
		}

		routers[route.Path] = router
	}

	return routers, nil
}
