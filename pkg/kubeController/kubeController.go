package kubeController

import (
	"errors"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"multicluster-ingress-controller/pkg/events"
	"multicluster-ingress-controller/pkg/kubeControls"
	"multicluster-ingress-controller/pkg/kubernetesAPIServer"
	"multicluster-ingress-controller/pkg/namespaces"
	"multicluster-ingress-controller/pkg/swaggerInterface"
	"net/http"
)

func NewkubeController() *KubeController {
	return &KubeController{}
}

type KubeController struct {
	Api        *kubernetesAPIServer.KubernetesAPIServer
	ServerList []string
	kubeState  kubeControls.KubeControl
}

func (kC *KubeController) DeleteAll() {
	kC.kubeState.DeleteAll()
	kC.ServerList = nil
	kC.Api = nil
}

func (kC *KubeController) PopulateKubeController(api *kubernetesAPIServer.KubernetesAPIServer, swgClientServer *swaggerInterface.CreateClientServer) error {
	kC.Api = api
	kC.ServerList = swgClientServer.ServerURL
	kCtls := kubeControls.New()
	err := kCtls.PopulateKubeControls(*swgClientServer)
	if err != nil {
		return err
	}
	kC.kubeState = *kCtls
	return nil
}

func (kC *KubeController) RunWatchEvents() error {
	kS := kC.kubeState
	var err error
	for ns := 0; ns < kS.Namespaceslength(); ns++ {
		var namespacePtr *namespaces.Namespace
		namespacePtr, err = kS.GetNamespace(ns)
		if err != nil {
			break
		}
		namespace := *namespacePtr
		namespaceName := namespace.GetNamespaceName()

		for ev := 0; ev < namespace.EventsLength(); ev++ {
			var eventPtr *events.Event
			eventPtr, err = namespace.GetEvent(ev)
			if err != nil {
				break
			}
			event := *eventPtr
			eventName := event.GetEventName()
			eventStopChannel := event.GetEventStopChannel()
			err = kC.watchEvent(namespaceName, eventName, eventStopChannel)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		kC.DeleteAll()
		return err
	}
	return nil
}

func (kC *KubeController) watchEvent(namespace string, eventName string, eventStopChannel chan struct{}) error {
	switch eventName {
	case events.Endpoints:
		kC.EndpointWatchPerNamespace(namespace, eventStopChannel)
	case events.Pods:
		kC.PodWatchPerNamespace(namespace, eventStopChannel)
	case events.Services:
		kC.ServiceWatchPerNamespace(namespace, eventStopChannel)
	case events.Ingresses:
		kC.IngressWatchPerNamespace(namespace, eventStopChannel)
	case events.Secrets:
		kC.SecretWatchPerNamespace(namespace, eventStopChannel)
	case events.Namespaces:
		// No watch event for Namespaces Entity
	default:
		return errors.New("INVALID event Name: " + eventName)

	}
	return nil
}

func (kC *KubeController) IngressWatchPerNamespace(namespace string, stop chan struct{}) {
	api := kC.Api
	klog.Info("[INFO] Watch Ingress for namespace: ", namespaces.GetNamespaceString(namespace))

	ingressListWatcher := cache.NewListWatchFromClient(api.Client.ExtensionsV1beta1().RESTClient(), "ingresses", namespace, fields.Everything())
	_, controller := cache.NewInformer(ingressListWatcher, &v1beta1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			kC.ingressEventParseAndSendData(obj, "ADDED")
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			kC.ingressEventParseAndSendData(newobj, "MODIFIED")
		},
		DeleteFunc: func(obj interface{}) {
			kC.ingressEventParseAndSendData(obj, "DELETED")
		},
	},
	)
	go controller.Run(stop)

	return
}

func (kC *KubeController) SecretWatchPerNamespace(namespace string, stop chan struct{}) {
	api := kC.Api
	klog.Info("[INFO] Watch Secret events for namespace: ", namespaces.GetNamespaceString(namespace))

	secretListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "secrets", namespace, fields.Everything())
	_, controller := cache.NewInformer(secretListWatcher, &v1.Secret{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			kC.secretEventParseAndSendData(obj, "ADDED")
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			kC.secretEventParseAndSendData(newobj, "MODIFIED")
		},
		DeleteFunc: func(obj interface{}) {
			kC.secretEventParseAndSendData(obj, "DELETED")
		},
	},
	)
	go controller.Run(stop)
	return
}

func (kC *KubeController) EndpointWatchPerNamespace(namespace string, stop chan struct{}) {
	api := kC.Api
	klog.Info("[INFO] Watch Endpoint for namespace: ", namespaces.GetNamespaceString(namespace))

	endpointListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "endpoints", namespace, fields.Everything())
	_, controller := cache.NewInformer(endpointListWatcher, &v1.Endpoints{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			kC.endpointEventParseAndSendData(obj, "ADDED")
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			kC.endpointEventParseAndSendData(newobj, "MODIFIED")
		},
		DeleteFunc: func(obj interface{}) {
			kC.endpointEventParseAndSendData(obj, "DELETED")
		},
	},
	)

	go controller.Run(stop)
	return
}

func (kC *KubeController) ServiceWatchPerNamespace(namespace string, stop chan struct{}) {
	api := kC.Api
	klog.Info("[INFO] Watch Service events for namespace: ", namespaces.GetNamespaceString(namespace))
	serviceListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourceServices), namespace, fields.Everything())
	_, controller := cache.NewInformer(serviceListWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			kC.serviceEventParseAndSendData(obj, "ADDED")
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			kC.serviceEventParseAndSendData(newobj, "MODIFIED")
		},
		DeleteFunc: func(obj interface{}) {
			kC.serviceEventParseAndSendData(obj, "DELETED")
		},
	},
	)

	go controller.Run(stop)
	return
}

func (kC *KubeController) PodWatchPerNamespace(namespace string, stop chan struct{}) {
	api := kC.Api

	klog.Info("[INFO] Watch Pods for namespace: ", namespaces.GetNamespaceString(namespace))
	PodListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourcePods), namespace, fields.Everything())
	_, controller := cache.NewInformer(PodListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			kC.podEventParseAndSendData(obj, "ADDED")
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			kC.podEventParseAndSendData(newobj, "MODIFIED")
		},
		DeleteFunc: func(obj interface{}) {
			kC.podEventParseAndSendData(obj, "DELETED")
		},
	},
	)
	go controller.Run(stop)
	return
}

func (kC *KubeController) endpointEventParseAndSendData(obj interface{}, eventType string) {

	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var objJson v1.Endpoints
	if err = json.Unmarshal(objByte, &objJson); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	if objJson.ObjectMeta.Namespace == "kube-system" {
		return
	}
	kC.parseAndSendData(string(message), objJson.ObjectMeta, objJson.TypeMeta, "Endpoints", eventType)
}
func (kC *KubeController) secretEventParseAndSendData(obj interface{}, eventType string) {

	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var objJson v1.Secret
	if err = json.Unmarshal(objByte, &objJson); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	if objJson.ObjectMeta.Namespace == "kube-system" {
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	kC.parseAndSendData(string(message), objJson.ObjectMeta, obj.(*v1.Secret).TypeMeta, "Secret", eventType)

}

func (kC *KubeController) ingressEventParseAndSendData(obj interface{}, eventType string) {

	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var objJson v1beta1.Ingress
	if err = json.Unmarshal(objByte, &objJson); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	if objJson.ObjectMeta.Namespace == "kube-system" {
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	kC.parseAndSendData(string(message), objJson.ObjectMeta, obj.(*v1beta1.Ingress).TypeMeta, "Ingress", eventType)

}

func (kC *KubeController) podEventParseAndSendData(obj interface{}, eventType string) {
	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var objJson v1.Pod
	if err = json.Unmarshal(objByte, &objJson); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	if objJson.ObjectMeta.Namespace == "kube-system" {
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	kC.parseAndSendData(string(message), objJson.ObjectMeta, objJson.TypeMeta, "Pod", eventType)
}

func (kC *KubeController) serviceEventParseAndSendData(obj interface{}, eventType string) {
	objByte, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var objJson v1.Service
	if err = json.Unmarshal(objByte, &objJson); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	if objJson.ObjectMeta.Namespace == "kube-system" {
		return
	}
	message, err := json.MarshalIndent(objJson, "", "  ")
	kC.parseAndSendData(string(message), objJson.ObjectMeta, objJson.TypeMeta, "Service", eventType)
}

func GenerateUUID() string {
	uuid := uuid.New()
	s := uuid.String()
	return s
}

func (kC *KubeController) parseAndSendData(obj string, metaData metav1.ObjectMeta, metaHeader metav1.TypeMeta, kind string, objtype string) {
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
	kC.sendData(bytes.NewBuffer(respJson))
}

func (kC *KubeController) sendData(data *bytes.Buffer) {
	contr := *kC
	//servers.mux.Lock()
	fmt.Print("[INFO] Sending the data")
	for _, v := range contr.ServerList {
		tmp_data := *data
		result, err := http.Post(v, "application/json", &tmp_data)
		if err != nil {
			fmt.Print("[INFO]", result, err)
		}

	}
	//servers.mux.Unlock()
}
