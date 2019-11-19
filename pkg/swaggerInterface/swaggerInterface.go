package swaggerInterface

import (
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"multicluster-ingress-controller/pkg/kubernetesAPIServer"
)

type CreateClientServer struct {
	ClusterName       string
	ConfigFileName    string
	KubeURL           string
	KubeServAcctToken string
	Namespaces        []string
	ServerURL         []string
	WatchEvents       []string
}

func New(configFile, kubeURL, kubeServAcctToken string) *CreateClientServer {
	return &CreateClientServer{
		ConfigFileName:    configFile,
		KubeURL:           kubeURL,
		KubeServAcctToken: kubeServAcctToken,
	}
}

func (cS *CreateClientServer) ValidateKubeClusterFields() error {

	configFile := cS.ConfigFileName
	kubeURL := cS.KubeURL
	kubeServAcctToken := cS.KubeServAcctToken

	if configFile == "" && kubeURL == "" && kubeServAcctToken == "" {
		fmt.Println("Service Account details are collected from POD")
		return nil
	}
	if configFile == "" && kubeURL != "" && kubeServAcctToken != "" {
		fmt.Println("Using kubeURL and kubeServAcctToken field for Kube Client configuration")
		return nil
	}
	if configFile != "" {
		fmt.Println("Using configFile for configuration")
		return nil
	}
	if configFile == "" && kubeURL != "" && kubeServAcctToken == "" {
		errString := "Kubernetes Cluster Service Account Token Not Present: Connection Parameters invalid"
		fmt.Println(errString)
		return errors.New(errString)
	}
	if configFile == "" && kubeURL == "" && kubeServAcctToken != "" {
		errString := "Kubernetes Cluster URL Not Present: Connection Parameters invalid"
		fmt.Println(errString)
		return errors.New(errString)
	}
	return errors.New("Invalid Input Kubernetes connection Parameters")
}

/* This API creates Kubernetes client session. API requires config file or Kubernetes URL and Kubernetes Service Account Token.
 * If file is not in default location, provide with path of the file.
 */
func (cS *CreateClientServer) CreateK8sApiserverClient() (*kubernetesAPIServer.KubernetesAPIServer, error) {
	configFile := cS.ConfigFileName

	klog.Info("[INFO] Creating API Client", configFile)
	api := kubernetesAPIServer.New()
	config, err := cS.getConfig()
	if err != nil {
		klog.Error("[ERROR] Failed to obtain Kubernetes Config")
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error("[ERROR] Failed to establish connection")
		klog.Fatal(err)
	}
	klog.Info("[INFO] Kubernetes Client is created")
	api.SetClient(client)
	return api, nil
}

// This function provides the kube client config file with the provided inputs
func (cS *CreateClientServer) getConfig() (*restclient.Config, error) {
	configFile := cS.ConfigFileName
	kubeURL := cS.KubeURL
	kubeServAcctToken := cS.KubeServAcctToken

	if configFile != "" {
		config, err := clientcmd.BuildConfigFromFlags("", configFile)
		if err != nil {
			klog.Error("[ERROR] Did not find valid kube config info")
			return nil, err
		}
		return config, err
	} else {
		if configFile == "" && kubeURL == "" && kubeServAcctToken == "" {
			config, err := clientcmd.BuildConfigFromFlags("", "")
			if err != nil {
				klog.Error("[WARNING] Citrix Node Controller Failed to create a Client")
				return nil, err
			}
			return config, err
		} else {
			/* A valid KubeURL and Token scenario as the validation for the integrity
			   of input Kube parameters are done at validateKubeClusterFields()
			*/
			return &restclient.Config{
				Host:            kubeURL,
				BearerToken:     kubeServAcctToken,
				TLSClientConfig: restclient.TLSClientConfig{Insecure: true},
			}, nil
		}
	}
}
