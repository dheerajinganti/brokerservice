### Manage go module dependency

        go mod init github.com/dheerajinganti/brokerservice
        go get "github.com/gorilla/mux"

        list dependencies
        go list -m all
        
### push golang app to cloud foundry
#### manifest.yaml
---
        applications:
        - name: brokerservice
        env:
            GO_INSTALL_PACKAGE_SPEC: github.com/dheerajinganti/brokerservice

### create cf service broker 
        cf create-service-broker test-broker test test123 https://brokerservice.130.147.139.72.nip.io

#### enable service access
        cf enable-service-access hsdp-rds -b test-broker

### Dockerize go application 

#### docker build
        docker build --tag hsop-local-docker.artifactory.pic.philips.com/pg-broker:0.1 .
#### docker run 
        docker run -d --name pg-broker --restart=always --publish 8888:8888 hsop-local-docker.artifactory.pic.philips.com/pg-broker:0.2
#### push docker image to registry
        docker image push hsop-local-docker.artifactory.pic.philips.com/pg-broker:0.2 

#### create imagepull secret for kubernetes deployment

        kubectl create secret docker-registry afpullsecret -n pg-test --docker-server=https://hsop-local-docker.artifactory.pic.philips.com --docker-username=ing08763 --docker-email=admin@philips.com --docker-password=""

#### helm charts for golang application
        helm create charts

#### helm dry run before deployment
        helm install pg-broker --dry-run --debug ./charts


#### http requests from Cloud Foundry as per OSBI specification

        PUT /v2/service_instances/abb0736d-8941-47f0-aca9-fc449c0d067d?accepts_incomplete=true
    
        DELETE /v2/service_instances/abb0736d-8941-47f0-aca9-fc449c0d067d?accepts_incomplete=true&plan_id=23332639-fbc1-49e7-ab24-52b586860fef&service_id=00fc4084-4ea1-40b2-8db7-55d040c8c683

#### request body payload for createinstance:

        {
          "service_id": "00fc4084-4ea1-40b2-8db7-55d040c8c683",
          "plan_id": "23332639-fbc1-49e7-ab24-52b586860fef",
          "organization_guid": "ba414500-ea59-4851-bc05-1a5c40bbfde8",
          "context": {
          	"platform": "cloudfoundry",
              "organization_guid": "ba414500-ea59-4851-bc05-1a5c40bbfde8",
              "space_guid": "3e8e3238-8d2a-4646-9f6e-fc398908586a",
              "organization_name": "hsop",
              "space_name": "dev",
              "instance_name": "test-postgres"
             }
        }

## InitBroker function

        func InitBroker() error {

            log.Println("init: load broker config file")
            data, err := ioutil.ReadFile("config.yml")
            if err != nil {
                log.Fatalln(err)
            }

            //unmarshal yaml config into Catalog data structure
            err = yaml.Unmarshal([]byte(data), &m)
            if err != nil {
                log.Printf("err: %v\n", err)
                return err
            }
            //fmt.Println(m.Services[0])

            // d, err := yaml.Marshal(m.Services)
            // if err != nil {
            // 	log.Fatalf("error: %v", err)
            // }
            //fmt.Printf("--- m dump:\n%s\n\n", string(d))

            //convert yam into json document
            d, err := yaml.Marshal(m)
            jData, err := yaml.YAMLToJSON([]byte(d))
            if err != nil {
                fmt.Printf("err: %v\n", err)
                return err
            }
            //strData := string(jData)
            //fmt.Println(strData)
            var jCatalog model.Catalog
            err = json.Unmarshal(jData, &jCatalog)

            //extract Services part from servicebroker configuration for catalog api to return
            var service []model.Service
            //b, err := json.Marshal(jCatalog.Services)
            b, err := json.Marshal(m.Services)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            //fmt.Println(string(b))
            //b1 := string(b)
            err = json.Unmarshal(b, &service)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            //fmt.Println(service)
            for _, s := range service {
                for _, t := range s.Tags {
                    fmt.Println(t)
                }
            }
            //extract PlanSettings
            var plansettings []model.PlanSetting
            p, err := json.Marshal(m.PlanSettings)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            //unmarshal plan settings into PlanSetting type
            err = json.Unmarshal(p, &plansettings)
            if err != nil {
                log.Fatalf("error: %v", err)
            }
            for _, ps := range plansettings {
                fmt.Println(ps.ID)
            }

            // var jSvcCat model.SvcCatalog
            // err = json.Unmarshal([]byte(b1), &jSvcCat)
            // if err != nil {
            // 	log.Fatalf("error: %v", err)
            // }
            // fmt.Println(jSvcCat.Services[0].Description)

            //svcCat = jCatalog.Services

            // j3 := ("{services:") + string(j2) + "}"
            // fmt.Println(string(j3))

            //json.dumps({'services': CONFIG['services']})
            //	return
            //var c Catalog
            // m := make(map[interface{}]interface{})
            // err = yaml.Unmarshal([]byte(data), &m)
            // if err != nil {
            // 	fmt.Printf("Error parsing YAML file: %s\n", err)
            // }
            // fmt.Printf("--- m:\n%v\n\n", m)

            // yaml, err := simpleyaml.NewYaml(data)
            // if err != nil {
            // 	log.Fatalln(err)
            // }
            // fmt.Println(yaml.Get("services").GetIndex(0).Get("plans").GetIndex(0))
            // tags := yaml.Get("services").GetIndex(0).Get("tags")
            // fmt.Println(tags)
            // fmt.Println(yaml.Get("plan_settings").GetIndex(0).Get("image_name"))
            // svc := yaml.Get("services")
            // fmt.Println(svc)

            // err = json.Unmarshal([]byte(data), &m)
            // fmt.Printf("--- m:\n%v\n\n", m)

            return nil
        }

## Catalog function

    //Catalog ...
func (c *Controller) Catalog(w http.ResponseWriter, r *http.Request) {

	fmt.Println("getting catalog for service")

	d, err := yaml.Marshal(c.Config.Services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	d1, err := json.Marshal(c.Config.Services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	d2 := ("{\"services\":") + string(d1) + "}"

	//fmt.Printf("--- m dump:\n%s\n\n", string(d))
	j2, err := yaml.YAMLToJSON([]byte(d))
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	j3 := ("{\"services\":") + string(j2) + "}"
	fmt.Println(string(j3))

	var svc model.SvcCatalog
	//err = json.Unmarshal([]byte(j3), &svc)
	err = json.Unmarshal([]byte(d2), &svc)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	WriteResponse(w, http.StatusOK, svc)
}