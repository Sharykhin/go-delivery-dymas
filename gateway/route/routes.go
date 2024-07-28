package route

import "fmt"

func CreateServicesRoutes() map[string]string {
	courierService := fmt.Sprintf(
		"/couriers/{id:%s}",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	courierLocationService := fmt.Sprintf(
		"/courier/{courier_id:%s}/location",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	orderService := fmt.Sprintf(
		"/orders/%s",
		"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	)

	routesServices := map[string]string{
		courierService:         "http://localhost:9667",
		courierLocationService: "http://localhost:8081",
		orderService:           "http://localhost:6661",
		"/orders":              "http://localhost:6661",
	}

	return routesServices
}
