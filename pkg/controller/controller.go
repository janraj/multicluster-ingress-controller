package controller
  
import (
	"github.com/google/uuid"
	v1beta1 "k8s.io/api/extensions/v1beta1"
        "fmt"
        "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" 
        restclient "k8s.io/client-go/rest"
	//testapi "k8s.io/client-go/tools/cache/testing"
        "os"
        "path/filepath"
	"net/http"
	"strings"
	"errors"
	"multicluster-ingress-controller/pkg/events"
	"multicluster-ingress-controller/pkg/swaggerToKubeInterface"
	"multicluster-ingress-controller/pkg/kubernetesAPIServer"
	"multicluster-ingress-controller/pkg/swaggerInterface"
)

const (
	endpoints string = events.Endpoints
	services string = events.Services
	pods string = events.Pods
	secrets string = events.Secrets
	ingresses string = events.Ingresses
)

var (
        kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
        config     *restclient.Config
        err        error
        podcount = 0
	endpointAPI = "http://10.221.42.109:30132/sdc/nitro/v1/config/app_event"
	//endpointAPI = "http://localhost:8080/api/v1/pods"
)



func GenerateUUID() string {
	uuid := uuid.New()
	s := uuid.String()
	return s
}

func GetKubeEndpointsAll(api *kubernetesAPIServer.KubernetesAPIServer) *v1.EndpointsList {
	fmt.Println("Endpoints GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Endpoints(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeEndpointsNamespace(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1.EndpointsList {
	fmt.Println("Endpoints GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Endpoints(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeEndpointsName(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Endpoints {
    fmt.Println("Endpoints GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Endpoints(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}
func EndpointGet(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1.EndpointsList {
        fmt.Println("ENDPOINT GET API: Calling kubernetes API server")
	endpointslist, err := api.Client.Core().Endpoints(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	/*
	fmt.Println(endpointslist)
	obj, err := json.Marshal(endpointslist)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original endpoint object: %v", err)
        }
        var objJson v1.EndpointsList
	if err = json.Unmarshal(obj, &objJson); err != nil {
        	klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
        return string(message)
	*/
	return endpointslist
}

func NamespaceGet(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Namespace {
        fmt.Println("NAMESPACE Name GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func GetKubePodsAll(api *kubernetesAPIServer.KubernetesAPIServer) *v1.PodList {
	fmt.Println("POD GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubePodsNamespace(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1.PodList {
	fmt.Println("POD GET NAMESPACE API: Calling kubernetes API server")
	obj, err := api.Client.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubePodsName(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Pod {
    fmt.Println("POD GET Name API: Calling kubernetes API server")
	obj, err := api.Client.Core().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func GetKubeSecretsAll(api *kubernetesAPIServer.KubernetesAPIServer) *v1.SecretList {
	fmt.Println("SECRET GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeSecretsNamespace(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1.SecretList {
	fmt.Println("SECRET Namespace GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeSecretsName(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Secret {
    fmt.Println("SECRET Name GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func GetKubeIngressesAll(api *kubernetesAPIServer.KubernetesAPIServer) *v1beta1.IngressList {
	fmt.Println("INGRESS GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.ExtensionsV1beta1().Ingresses(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeIngressesNamespace(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1beta1.IngressList {
	fmt.Println("INGRESS Namespace GET API: Calling kubernetes API server")
	obj, err := api.Client.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}


func GetKubeIngressesName(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1beta1.Ingress {
    fmt.Println("INGRESS Name GET API: Calling kubernetes API server")
	obj, err := api.Client.ExtensionsV1beta1().Ingresses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}
func GetKubeServicesAll(api *kubernetesAPIServer.KubernetesAPIServer) *v1.ServiceList {
	fmt.Println("SERVICE GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Services(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}
func GetKubeServicesNamespace(api *kubernetesAPIServer.KubernetesAPIServer, namespace string) *v1.ServiceList {
	fmt.Println("SERVICE Namespace GET ALL API: Calling kubernetes API server")
	obj, err := api.Client.Core().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func GetKubeServicesName(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Service {
	return ServiceGet(api, namespace, name)
}

func ServiceGet(api *kubernetesAPIServer.KubernetesAPIServer, namespace string, name string) *v1.Service {
    fmt.Println("SERVICE Name GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func StartController(swg2KubeIntfPtr *swaggerToKubeInterface.SwaggerToKubeInterface) (int, string) {
	
     err := swg2KubeIntfPtr.LinkAndStartWatchEvents()
     if err != nil {
     	return http.StatusBadRequest, err.Error()
     }
     
     return http.StatusOK, "Controller Added"
}
func GetKubeEvents(configFile, kubeURL, kubeServAcctToken string, event_params ...string) (interface{}, error) {
	cS := swaggerInterface.New(configFile, kubeURL, kubeServAcctToken)
    api, err := cS.CreateK8sApiserverClient() 
	if (err != nil) {
		fmt.Println("GetKubeEvents: Error while starting client API session")
		return "",err
	}
	fmt.Println("GetKubeEvents: event_params ", event_params)
	entity := event_params[0]
	switch entity {
		case pods:
			return getKubePodEvents(api, event_params)
		case ingresses:
			return getKubeIngressEvents(api, event_params)
		case services:
			return getKubeServiceEvents(api, event_params)
		case secrets:
			return getKubeSecretEvents(api, event_params)
		case endpoints:
			return getKubeEndpointEvents(api, event_params)
		default:
			fmt.Println("Wrong type of Event_params", event_params[0])
			return "", nil
	}
}

func getKubePodEvents(api *kubernetesAPIServer.KubernetesAPIServer, event_params []string) (interface{}, error) {
	var message interface{}
	switch len(event_params) {
		case 1: {
			message = GetKubePodsAll(api)
		}
		case 2: {
			namespace := event_params[1]
			message = GetKubePodsNamespace(api, namespace)
		}
		case 3: {
			namespace := event_params[1]
			name := event_params[2]
			message = GetKubePodsName(api, namespace, name)
		}
		default: {
			fmt.Println("Wrong number of event_params", event_params)
			return nil, errors.New("Wrong number of event_params" + strings.Join(event_params, " "))
		}
	}
	return message, nil
}
func getKubeServiceEvents(api *kubernetesAPIServer.KubernetesAPIServer, event_params []string) (interface{}, error) {
	var message interface{}
	switch len(event_params) {
		case 1: {
			message = GetKubeServicesAll(api)
		}
		case 2: {
			namespace := event_params[1]
			message = GetKubeServicesNamespace(api, namespace)
		}
		case 3: {
			namespace := event_params[1]
			name := event_params[2]
			message = GetKubeServicesName(api, namespace, name)
		}
		default: {
			fmt.Println("Wrong number of event_params", event_params)
			return nil, errors.New("Wrong number of event_params" + strings.Join(event_params, " "))
		}
	}
	return message, nil
}
func getKubeEndpointEvents(api *kubernetesAPIServer.KubernetesAPIServer, event_params []string) (interface{}, error) {
	var message interface{}
	switch len(event_params) {
		case 1: {
			message = GetKubeEndpointsAll(api)
		}
		case 2: {
			namespace := event_params[1]
			message = GetKubeEndpointsNamespace(api, namespace)
		}
		case 3: {
			namespace := event_params[1]
			name := event_params[2]
			message = GetKubeEndpointsName(api, namespace, name)
		}
		default: {
			fmt.Println("Wrong number of event_params", event_params)
			return nil, errors.New("Wrong number of event_params" + strings.Join(event_params, " "))
		}
	}
	return message, nil
}
func getKubeSecretEvents(api *kubernetesAPIServer.KubernetesAPIServer, event_params []string) (interface{}, error) {
	var message interface{}
	switch len(event_params) {
		case 1: {
			message = GetKubeSecretsAll(api)
		}
		case 2: {
			namespace := event_params[1]
			message = GetKubeSecretsNamespace(api, namespace)
		}
		case 3: {
			namespace := event_params[1]
			name := event_params[2]
			message = GetKubeSecretsName(api, namespace, name)
		}
		default: {
			fmt.Println("Wrong number of event_params", event_params)
			return nil, errors.New("Wrong number of event_params" + strings.Join(event_params, " "))
		}
	}
	return message, nil
}
func getKubeIngressEvents(api *kubernetesAPIServer.KubernetesAPIServer, event_params []string) (interface{}, error) {
	var message interface{}
	switch len(event_params) {
		case 1: {
			message = GetKubeIngressesAll(api)
		}
		case 2: {
			namespace := event_params[1]
			message = GetKubeIngressesNamespace(api, namespace)
		}
		case 3: {
			namespace := event_params[1]
			name := event_params[2]
			message = GetKubeIngressesName(api, namespace, name)
		}
		default: {
			fmt.Println("Wrong number of event_params", event_params)
			return nil, errors.New("Wrong number of event_params" + strings.Join(event_params, " "))
		}
	}
	return message, nil
}

func GetK8sEvents(configFile, kubeURL, kubeServAcctToken, event, namespace, name string) (interface{}, error){
	 cS := swaggerInterface.New(configFile, kubeURL, kubeServAcctToken)
     api, err := cS.CreateK8sApiserverClient() 
     var message interface{}
     if (err != nil){
		fmt.Println("Error while starting client API session")
		return "",err
     }
     fmt.Println("GetK8sEvents: ", event, namespace, name)
     if (strings.ToLower(event) == "endpoints"){
		message = EndpointGet(api, namespace)
		fmt.Printf("ENDPOINT API: List of endpoints retrieved %s", message)
     }
     if (strings.ToLower(event) == "namespace"){
		message = NamespaceGet(api, "", name)
		fmt.Printf("NAMESPACE API: List of all namespace retrieved %s", message)
     }
     if (strings.ToLower(event) == "service"){
		message = ServiceGet(api, namespace, name)
		fmt.Printf("SERVICE API: service retrieved %s", message)
     }
     return message, err
}
