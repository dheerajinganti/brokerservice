package kubeclient

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetkClient() (*rest.Config, error) {
	fmt.Println("inside getkClient")
	var configPath string

	// see if we are running inside k8s cluster
	_, kHost := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	_, kPort := os.LookupEnv("KUBERNETES_SERVICE_PORT")

	if kHost && kPort {
		c, err := rest.InClusterConfig()
		if err == nil {
			fmt.Println("Using in-cluster kube config")
			return c, nil
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}

	usr, err := user.Current()
	fmt.Println("current user is: " + usr.Username)
	configPath = filepath.Join(usr.HomeDir, ".kube", "config")

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c, err := clientcmd.DefaultClientConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s does not exist, using default kube config", configPath)
		return c, nil
	}

	//clientcmd.RESTConfigFromKubeConfig(b)
	c, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Using kube config '%s'", configPath)
	return c, nil

}
