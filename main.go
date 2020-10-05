package main

import (
	//"bytes"
	//"context"
	//"encoding/json"
	"fmt"

	//"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	//"gopkg.in/yaml.v2"
	//"github.com/smallfish/simpleyaml"
	//"github.com/ghodss/yaml"
	"github.com/gorilla/mux"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"

	"github.com/dheerajinganti/brokerservice/controller"
	//"github.com/dheerajinganti/brokerservice/kubeclient"
	"github.com/dheerajinganti/brokerservice/model"
	//appsv1 "k8s.io/api/apps/v1"
	//apiv1 "k8s.io/api/core/v1"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to my web api with mux router implementation")
	fmt.Println("this is a homepage!")
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
	appRouter.HandleFunc("/v2/service_instances/{instance_guid}", s.controller.CreateServiceInstance).Methods("PUT")
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
		return nil, err
	}
	return &Server{
		controller: controller,
	}, nil
}
