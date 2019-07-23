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

//This API creates Kubernetes client session. API requires config file. If file is not in default location, provide with path of the file.
func CreateK8sApiserverClient(configFile string) (*KubernetesAPIServer, error) {
        klog.Info("[INFO] Creating API Client", configFile)
        api := &KubernetesAPIServer{}
	if (configFile != "") {
                config, err = clientcmd.BuildConfigFromFlags("", configFile)
                if err != nil {
                        klog.Error("[ERROR] Did not find valid kube config info")
			return nil, err
                }
	}else {
        	config, err = clientcmd.BuildConfigFromFlags("", "")
        	if err != nil {
                	klog.Error("[WARNING] Citrix Node Controller Failed to create a Client")
			return nil, err
        	}
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
func StartController(configFile string, servers []string, events [] string){
     api, err := CreateK8sApiserverClient(configFile) 
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

func GetK8sEvents(configFile string, event string, namespace string, name string) (interface{}, error){
     api, err := CreateK8sApiserverClient(configFile) 
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
     return message, err
}
