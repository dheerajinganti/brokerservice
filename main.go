package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/dheerajinganti/brokerservice/controller"
	"github.com/dheerajinganti/brokerservice/model"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to my web api with mux router implementation")
	//fmt.Println("this is a homepage!")
}

//Server ...
type Server struct {
	controller *controller.Controller
}

// //Controller ...
// type Controller struct {
// 	k8sClient kubernetes.Interface
// }

var m model.Catalog

func main() {

	//err := controller.InitBroker()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	log.Println("creating server")
	s, err := createserver()
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.HandleFunc("/", homePage)
	appRouter.HandleFunc("/v2/catalog", s.controller.Catalog).Methods("GET")
	appRouter.HandleFunc("/v2/service_instances/{service_instance_guid}", s.controller.GetServiceInstance).Methods("GET")
	appRouter.HandleFunc("/v2/service_instances/{instance_guid}", s.controller.CreateServiceInstance).Methods("PUT")
	appRouter.HandleFunc("/v2/service_instances/{instance_id}/last_operation", s.controller.LastOperation).Methods("GET")
	appRouter.HandleFunc("/v2/service_instances/{instance_guid}", s.controller.RemoveServiceInstance).Methods("DELETE")
	appRouter.HandleFunc("/v2/service_instances/{instance_guid}/service_bindings/{service_binding_guid}", s.controller.Bind).Methods("PUT")
	appRouter.HandleFunc("/v2/service_instances/{instance_guid}/service_bindings/{service_binding_guid}", s.controller.UnBind).Methods("DELETE")

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8888"
	}
	log.Fatal(http.ListenAndServe(":"+port, appRouter))
}

func createserver() (*Server, error) {

	controller, err := controller.CreateController()
	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, errors.New("error")
	}
	return &Server{
		controller: controller,
	}, nil
}
