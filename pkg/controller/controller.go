package controller
  
import (
	"github.com/google/uuid"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"bytes"
        "encoding/json"
        "fmt"
        "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/apimachinery/pkg/fields"
        "k8s.io/client-go/kubernetes"
        restclient "k8s.io/client-go/rest"
        "k8s.io/client-go/tools/cache"
        "k8s.io/client-go/tools/clientcmd"
	//testapi "k8s.io/client-go/tools/cache/testing"
        "k8s.io/klog"
        "os"
        "path/filepath"
	"net/http"
	"strings"
)

var (
        kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
        config     *restclient.Config
        err        error
        podcount = 0
	endpointAPI = "http://10.221.42.109:30132/sdc/nitro/v1/config/app_event"
	//endpointAPI = "http://localhost:8080/api/v1/pods"
)
func NewController() *Controller {
	return new(Controller)
}
type Controller struct {
        Api  *KubernetesAPIServer 
        ServerList []string
	EventList []string
}
//This is interface for Kubernetes API Server
type KubernetesAPIServer struct {
        Suffix string
        Client kubernetes.Interface
}

func GenerateUUID() string {
	uuid := uuid.New()
	s := uuid.String()
	return s
}
// This function provides the kube client config file with the provided inputs
func getConfig(configFile, kubeURL, kubeServAcctToken string) (*restclient.Config, error) {
	if (configFile != "") {
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
		if err != nil {
			klog.Error("[ERROR] Did not find valid kube config info")
			return nil, err
		}
		return config, err
	}else {
		if configFile == "" && kubeURL == "" && kubeServAcctToken == "" {
			config, err = clientcmd.BuildConfigFromFlags("", "")
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
				Host: kubeURL,
				BearerToken: kubeServAcctToken,
				TLSClientConfig: restclient.TLSClientConfig{Insecure: true},
			}, nil
		}
	}
}

//This API creates Kubernetes client session. API requires config file. If file is not in default location, provide with path of the file.
func CreateK8sApiserverClient(configFile , kubeURL, kubeServAcctToken string) (*KubernetesAPIServer, error) {
        klog.Info("[INFO] Creating API Client", configFile)
        api := &KubernetesAPIServer{}
		config, err := getConfig(configFile, kubeURL, kubeServAcctToken)
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
        api.Client = client
        return api, nil
}

func ingressEventParseAndSendData(obj interface{}, eventType string, contr Controller) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1beta1.Ingress
	if err = json.Unmarshal(objByte, &objJson); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	if (objJson.ObjectMeta.Namespace == "kube-system"){
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta,  obj.(*v1beta1.Ingress).TypeMeta, "Ingress", eventType, contr)

}
func secretEventParseAndSendData(obj interface{}, eventType string, contr Controller) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Secret
	if err = json.Unmarshal(objByte, &objJson); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	if (objJson.ObjectMeta.Namespace == "kube-system"){
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta,  obj.(*v1.Secret).TypeMeta, "Secret", eventType, contr)

}
func IngressWatcher(api *KubernetesAPIServer, contr Controller) {
        ingressListWatcher := cache.NewListWatchFromClient(api.Client.ExtensionsV1beta1().RESTClient(), "ingresses", v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(ingressListWatcher, &v1beta1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			ingressEventParseAndSendData(obj, "ADDED", contr)
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			ingressEventParseAndSendData(newobj, "MODIFIED", contr)
                },
                DeleteFunc: func(obj interface{}) {
			ingressEventParseAndSendData(obj, "DELETED", contr)
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}
func SecretWatcher(api *KubernetesAPIServer, contr Controller) {
        secretListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "secrets", v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(secretListWatcher, &v1.Secret{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			secretEventParseAndSendData(obj, "ADDED", contr)
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			secretEventParseAndSendData(newobj, "MODIFIED", contr)
                },
                DeleteFunc: func(obj interface{}) {
			secretEventParseAndSendData(obj, "DELETED", contr)
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}
func endpointEventParseAndSendData(obj interface{}, eventType string, contr Controller) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Endpoints
	if err = json.Unmarshal(objByte, &objJson); err != nil {
        	klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
	if (objJson.ObjectMeta.Namespace == "kube-system"){
		return
	}
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Endpoints", eventType, contr)
}
func EndpointWatcher(api *KubernetesAPIServer, contr Controller) {
        endpointListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "endpoints", v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(endpointListWatcher, &v1.Endpoints{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			endpointEventParseAndSendData(obj, "ADDED", contr) 
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			endpointEventParseAndSendData(newobj, "MODIFIED", contr)
                },
                DeleteFunc: func(obj interface{}) {
			endpointEventParseAndSendData(obj, "DELETED", contr)
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}
func EndpointGet(api *KubernetesAPIServer, namespace string) *v1.EndpointsList {
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

func NamespaceGet(api *KubernetesAPIServer, namespace string, name string) *v1.Namespace {
        fmt.Println("NAMESPACE GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func ServiceGet(api *KubernetesAPIServer, namespace string, name string) *v1.Service {
    fmt.Println("SERVICE GET API: Calling kubernetes API server")
	obj, err := api.Client.Core().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return obj
}
func serviceEventParseAndSendData(obj interface{}, eventType string, contr Controller) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Service
	if err = json.Unmarshal(objByte, &objJson); err != nil {
        	klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	if (objJson.ObjectMeta.Namespace == "kube-system"){
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Service", eventType, contr)
}

func ServiceWatcher(api *KubernetesAPIServer, contr Controller) {
        serviceListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourceServices), v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(serviceListWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			serviceEventParseAndSendData(obj, "ADDED", contr)
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			serviceEventParseAndSendData(newobj, "MODIFIED", contr)
                },
                DeleteFunc: func(obj interface{}) {
			serviceEventParseAndSendData(obj, "DELETED", contr)
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}


func podEventParseAndSendData(obj interface{}, eventType string, contr Controller){
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Pod
	if err = json.Unmarshal(objByte, &objJson); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	if (objJson.ObjectMeta.Namespace == "kube-system"){
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Pod", eventType, contr)
}

func PodWatcher(api *KubernetesAPIServer, contr Controller) {
        PodListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourcePods), v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(PodListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			podEventParseAndSendData(obj, "ADDED", contr)
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			podEventParseAndSendData(newobj, "MODIFIED", contr)
                },
                DeleteFunc: func(obj interface{}) {
			podEventParseAndSendData(obj, "DELETED", contr)
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}

func parseAndSendData(obj string, metaData  metav1.ObjectMeta, metaHeader metav1.TypeMeta, kind string, objtype string, contr Controller) {
	resp := make(map[string]string)
	resp["app_event_id"] = GenerateUUID() 
	resp["resource_type"] = kind 
	resp["resource_name"] = metaData.Name
	resp["resource_generations"] = string(metaData.Generation)
	resp["resource_id"] = string(metaData.UID)
	resp["app_environment_id"] = "joan-test"
	resp["app_environment_type"] = "Kubernetes"
	resp["type"] = objtype
	resp["trig_time"] = ""
	resp["server_group_id"] = metaData.Namespace
	resp["resource_version"] = metaData.ResourceVersion
	resp["message"] = obj
	respJson, err := json.Marshal(resp)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        	fmt.Printf("[ERROR] Failed to Marshal original object: %v", err)
	}
	sendData(bytes.NewBuffer(respJson), contr)
}


func sendData(data *bytes.Buffer, contr Controller){
	//servers.mux.Lock()
        fmt.Print("[INFO] Sending the data")
	for _,v := range contr.ServerList {	
		tmp_data := *data
		result, err := http.Post(v, "application/json", &tmp_data)
		if (err != nil){
			fmt.Print("[INFO]", result, err)
		}
		
	}
	//servers.mux.Unlock()
}
func validateKubeClusterFields(configFile, kubeURL, kubeServAcctToken string) bool {
	if configFile == "" && kubeURL == "" && kubeServAcctToken == "" {
		fmt.Println("Using inbuilt POD values field for Kube client configuration")
		return true
	}
	if configFile == "" && kubeURL != "" && kubeServAcctToken != "" {
		fmt.Println("Using kubeURL and kubeServAcctToken field for Kube Client configuration")
		return true
	}
	if configFile != "" {
		fmt.Println("Using configFile for configuration")
		return true
	}
	if configFile == "" && kubeURL != "" && kubeServAcctToken == "" {
		fmt.Println("Kubernetes Cluster Service Account Token Not Present: Connection Parameters invalid")
		return false
	}
	if configFile == "" && kubeURL == "" && kubeServAcctToken != "" {
		fmt.Println("Kubernetes Cluster URL Not Present: Connection Parameters invalid")
		return false
	}
	return false
}
func StartController(configFile, kubeURL, kubeServAcctToken string, servers []string, events [] string){
	 if !validateKubeClusterFields(configFile, kubeURL, kubeServAcctToken) {
	 	 return
	 }
     api, err := CreateK8sApiserverClient(configFile, kubeURL, kubeServAcctToken) 
     if (err != nil){
		 fmt.Println("Error while starting client API session")
		 return
     }
     contr := Controller{}
     contr.Api = api
     contr.ServerList = servers
     contr.EventList = events 
     for _, event := range events {
	 if (strings.ToLower(event) == "ingresses"){
		IngressWatcher(api, contr)
         }
	 if (strings.ToLower(event) == "endpoints"){
		EndpointWatcher(api, contr)
         }
	 if (strings.ToLower(event) == "pods"){
		PodWatcher(api, contr)
         }
	 if (strings.ToLower(event) == "services"){
		ServiceWatcher(api, contr)
         }
	 if (strings.ToLower(event) == "secrets"){
		SecretWatcher(api, contr)
         }
     }
}

func GetK8sEvents(configFile, kubeURL, kubeServAcctToken, event, namespace, name string) (interface{}, error){
     api, err := CreateK8sApiserverClient(configFile, kubeURL, kubeServAcctToken) 
     var message interface{}
     if (err != nil){
		fmt.Println("Error while starting client API session")
		return "",err
     }
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
