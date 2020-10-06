package controller

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/dheerajinganti/brokerservice/kubeclient"
	"github.com/dheerajinganti/brokerservice/model"
	appsv1 "k8s.io/api/apps/v1"
)

//Controller ...
type Controller struct {
	K8sClient kubernetes.Interface
	Config    *model.Catalog
}

// CreateController ...
func CreateController() (*Controller, error) {

	config, err := InitBroker()
	if err != nil {
		log.Println("Error getting kubeconfig")
	}

	kubeConfig, err := kubeclient.GetkClient()
	if err != nil {
		log.Println("Error getting kubeconfig")
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	// //test list pods
	// podlist, err := clientset.CoreV1().Pods("kubecf").List(context.TODO(), metav1.ListOptions{})

	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("******Listing all pods******")
	// for _, pod := range podlist.Items {
	// 	fmt.Println(pod.Name)
	// }

	return &Controller{
		K8sClient: clientset,
		Config:    config,
	}, nil
}

// InitBroker initialize broker by loading catalog config
func InitBroker() (*model.Catalog, error) {

	log.Println("init: load broker config file")
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	var m model.Catalog
	//unmarshal yaml config into Catalog data structure
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, err
	}

	//TODO: this is a test code, remove it before
	for _, s := range m.Services {
		for _, t := range s.Tags {
			fmt.Println(t)
		}
	}
	for _, ps := range m.PlanSettings {
		fmt.Println(ps.ID)
	}
	return &m, nil
}

//Catalog ...
func (c *Controller) Catalog(w http.ResponseWriter, r *http.Request) {

	log.Println("getting catalog for service")

	d1, err := json.Marshal(c.Config.Services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	d2 := ("{\"services\":") + string(d1) + "}"

	var svc model.SvcCatalog
	err = json.Unmarshal([]byte(d2), &svc)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteResponse(w, http.StatusOK, svc)
}

//WriteResponse ...
func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(data))
}

//CreateServiceInstance ...requesy body payload
//{
//	"service_id": "00fc4084-4ea1-40b2-8db7-55d040c8c683",
//  "plan_id": "23332639-fbc1-49e7-ab24-52b586860fef",
//  "organization_guid": "ba414500-ea59-4851-bc05-1a5c40bbfde8",
//  "context": {
//  	"platform": "cloudfoundry",
//      "organization_guid": "ba414500-ea59-4851-bc05-1a5c40bbfde8",
//      "space_guid": "3e8e3238-8d2a-4646-9f6e-fc398908586a",
//      "organization_name": "hsop",
//      "space_name": "dev",
//      "instance_name": "test-postgres"
//     }
//}
func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Service Instance...")

	instance_id := mux.Vars(r)["instance_guid"]
	log.Println("Instance_id is: " + instance_id)
	log.Println("read request body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("received body: " + string(body))

	//unmarshal body into ProvisionDetails type
	var pd model.ProvisionDetails
	err = json.Unmarshal(body, &pd)
	if err != nil {
		fmt.Println(err)
	}
	//filter the plan based on the planid
	var ps model.PlanSetting
	var bFound bool
	for _, p := range c.Config.PlanSettings {
		fmt.Println(p.ID)
		if pd.PlanID == p.ID {
			bFound = true
			ps = p
			log.Println("Found plan id for creating instance")
			break
		}
	}
	if !bFound {
		log.Println("No plans found for creating instance, return")
		return
	}

	log.Println("printing plan details", ps.CPULimit, ps.CPURequest, ps.ImageName)

	//log.Println(tmpl)
	password := "#test123"
	encodedPassword := b64.StdEncoding.EncodeToString([]byte(password))
	//create config parameters to render template
	si := model.ServiceInstance{
		instance_id,
		ps.MemoryRequest,
		ps.CPURequest,
		ps.ImageName,
		ps.MemoryLimit,
		ps.CPULimit,
		ps.Storage,
		encodedPassword,
		"vmware-str",
	}

	//pg config creation
	tmpl := template.Must(template.ParseFiles("./templates/postgres-configmap.yaml"))
	out := new(bytes.Buffer)
	err = tmpl.Execute(out, si)
	if err != nil {
		panic(err)
	}
	var pgConfig apiv1.ConfigMap
	err = yaml.Unmarshal([]byte(out.String()), &pgConfig)
	result, err := c.K8sClient.CoreV1().ConfigMaps("default").Create(context.TODO(), &pgConfig, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())

	//pg init-config creation
	tmpl = template.Must(template.ParseFiles("./templates/postgres-init-configmap.yaml"))
	out = new(bytes.Buffer)
	err = tmpl.Execute(out, si)
	if err != nil {
		panic(err)
	}
	var pgInitConfig apiv1.ConfigMap
	err = yaml.Unmarshal([]byte(out.String()), &pgInitConfig)
	result, err = c.K8sClient.CoreV1().ConfigMaps("default").Create(context.TODO(), &pgInitConfig, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("Created init-configmap %q.\n", result.GetObjectMeta().GetName())

	//pg statefulset creation
	tmpl = template.Must(template.ParseFiles("./templates/postgres-statefulset.yaml"))
	out = new(bytes.Buffer)
	err = tmpl.Execute(out, si)
	if err != nil {
		panic(err)
	}
	var pgStatefulSet appsv1.StatefulSet
	err = yaml.Unmarshal([]byte(out.String()), &pgStatefulSet)
	result1, err1 := c.K8sClient.AppsV1().StatefulSets("default").Create(context.TODO(), &pgStatefulSet, metav1.CreateOptions{})
	if err1 != nil {
		panic(err1)
	}
	log.Printf("Created Statefulset %q.\n", result1.GetObjectMeta().GetName())

	//pg service creation
	tmpl = template.Must(template.ParseFiles("./templates/postgres-service.yaml"))
	out = new(bytes.Buffer)
	err = tmpl.Execute(out, si)
	if err != nil {
		panic(err)
	}
	var svcConfig apiv1.Service
	err = yaml.Unmarshal([]byte(out.String()), &svcConfig)
	svcresult, err2 := c.K8sClient.CoreV1().Services("default").Create(context.TODO(), &svcConfig, metav1.CreateOptions{})
	if err2 != nil {
		panic(err2)
	}
	log.Printf("Created Service %q.\n", svcresult.GetObjectMeta().GetName())

	//pg secret cretion
	tmpl = template.Must(template.ParseFiles("./templates/postgres-secret.yaml"))
	out = new(bytes.Buffer)
	err = tmpl.Execute(out, si)
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())

	var secretConfig apiv1.Secret
	err = yaml.Unmarshal([]byte(out.String()), &secretConfig)
	secresult, err3 := c.K8sClient.CoreV1().Secrets("default").Create(context.TODO(), &secretConfig, metav1.CreateOptions{})
	if err3 != nil {
		panic(err3)
	}
	log.Printf("Created Service %q.\n", secresult.GetObjectMeta().GetName())

	WriteResponse(w, http.StatusOK, "CreateServiceInstance is in works")

}

//GetServiceInstance ...
func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
	log.Println("Get Service Instance State....")
	instance_id := mux.Vars(r)["instance_guid"]
	log.Println("Instance_id is: " + instance_id)
	log.Println("read request body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("received body: " + string(body))
	WriteResponse(w, http.StatusOK, "GetServiceInstance is in works")
}

//RemoveServiceInstance ...
func (c *Controller) RemoveServiceInstance(w http.ResponseWriter, r *http.Request) {
	log.Println("Remove Service Instance...")
	instance_id := mux.Vars(r)["instance_guid"]
	log.Println("Instance_id is: " + instance_id)
	log.Println("read request body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("received body: " + string(body))
	WriteResponse(w, http.StatusOK, "RemoveServiceInstance is in works")
}

//Bind ...
func (c *Controller) Bind(w http.ResponseWriter, r *http.Request) {
	log.Println("Bind Service Instance...")
	instance_id := mux.Vars(r)["instance_guid"]
	log.Println("Instance_id is: " + instance_id)
	log.Println("read request body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("received body: " + string(body))

	WriteResponse(w, http.StatusOK, "Bind is in works")
}

//UnBind ...
func (c *Controller) UnBind(w http.ResponseWriter, r *http.Request) {
	log.Println("Unbind Service Instance...")
	instance_id := mux.Vars(r)["instance_guid"]
	log.Println("Instance_id is: " + instance_id)
	log.Println("read request body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("received body: " + string(body))

	WriteResponse(w, http.StatusOK, "UnBind is in works")

}
