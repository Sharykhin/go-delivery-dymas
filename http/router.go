package http

import (
	"fmt"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	uuidRegExp := "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"
	url := "/courier/{courier_id:%s}/location"
	url = fmt.Sprintf(url, uuidRegExp)
	router.HandleFunc(url, handlerCouriersLocation).Methods("POST")

	return router
}
