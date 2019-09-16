package rest
  
import (
	"github.com/davecgh/go-spew/spew"
	"encoding/json"
        "github.com/gorilla/mux"
        "k8s.io/klog"
        "net/http"
        "time"
        "strings"
	"fmt"
	ctr "multicluster-ingress-controller/pkg/controller"
	"os"
)

func InitClientServer(){
	fmt.Println("Initializing Client Server Data Structure Which populate all k8s client and server details")
}

type createClientServer struct {
    	ClusterName string 
    	ConfigFileName string
    	KubeURL string
    	KubeServAcctToken string 
    	ServerURL []string 
    	WatchEvents []string 
}
var clientServerList []createClientServer
var routePath string = "/cic/nitro/v1/config/controller"
func CreateClientServerHandler(r *mux.Router) {

	// swagger:route GET /cic/nitro/v1/config/controller ClusterRegistration createClientServer
   	// ---
	// summary: This API lists all the configured entity details which incldue cluster names, config file path, watch events and the server list.
	// description: Returns 200 if it success. If there is no registeration, Error Not Found (404) will be returned.
   	// responses:
   	//   '200':
   	//     description: successful operation
   	//   '204':
   	//     description: successful operation, list is empty.
	r.HandleFunc(routePath, GetAllClientServer).Methods("GET")

	// swagger:operation POST /cic/nitro/v1/config/controller ClusterRegistration createClientServer
   	// ---
	// summary: This API adds cluster details which include cluster name, config path and list of servers.
	// description: Cluster Name can be any string. ConfigFileName must include relative path of kubernetes config file. ClusterName and ClusterFileName are mandatory argument.
	// parameters:
	// - name: RequestBody
	//   in: body
	//   description: Kubernetes API server URL
	//   schema:
        //     type: object
        //     required:
        //       - ClusterName
        //     properties:
        //       ClusterName:
        //          type: string
        //       ConfigFileName:
        //          type: string
        //       ServerURL:
        //          type: array
	//       WatchEvents:
	//          type: array
	//     example:
	//       ClusterName: ClusterABC
	//       ConfigFileName: /home/.kube/config
        //       ServerURL: [http://22.22.22.22, http://127.0.0.1:9000/]
	//       WatchEvents: [ingresses, endpoints]
	// responses:
	//  '200':
   	//     description: successful operation
	r.HandleFunc(routePath, PostClientServer).Methods("POST")
	
	// swagger:operation DELETE /cic/nitro/v1/config/controller ClusterRegistration createClientServer
   	// ---
	// summary: Delete the cluster details.
	// description: If there is no entity configured, Error Not Found (404) will be returned.
	// parameters:
	// - name: RequestBody
	//   in: body
	//   description: Kubernetes API server URL
	//   schema:
        //     type: object
        //     required:
        //       - ClusterName
        //       - ConfigFileName
        //     properties:
        //       ClusterName:
        //          type: string
        //       ConfigFileName:
        //          type: string
	//     example:
	//       ClusterName: ClusterABC
	//       ConfigFileName: /home/.kube/config
	// responses:
	//  '200':
   	//     description: successful operation
	//  '404':
   	//     description: entity did not find
	r.HandleFunc(routePath, DeleteClientServer).Methods("DELETE")
	
	// swagger:operation PUT /cic/nitro/v1/config/controller ClusterRegistration createClientServer
   	// ---
	// summary: This API can be used for updating the entities of a configured cluster.
	// description: If there is no matching entity, update operation cannot be performed. Error Not Found (404) will be returned.
	// parameters:
	// - name: RequestBody
	//   in: body
	//   description: Kubernetes API server URL
	//   schema:
        //     type: object
        //     required:
        //       - ClusterName
        //       - ConfigFileName
        //     properties:
        //       ClusterName:
        //          type: string
        //       ConfigFileName:
        //          type: string
	//     example:
	//       ClusterName: ClusterABC
	//       ConfigFileName: /home/.kube/config
	// responses:
	//  '200':
   	//     description: successful operation
	r.HandleFunc(routePath, UpdateClientServer).Methods("PUT")
	
}

type endpoints struct {
    	ClusterName string 
    	ConfigFileName string 
    	ServerURL []string 
    	WatchEvents []string 
}
var k8sClientPath string = "/cic/nitro/v1/config"
var k8sEndpointPath string = k8sClientPath + "/endpoints"
var k8sServicePath string = k8sClientPath + "/services"
var k8sIngressPath string = k8sClientPath + "/ingresses"
var k8sPodPath string = k8sClientPath + "/pods"
var k8sSecretPath string = k8sClientPath + "/secrets"
var baseURL string = "/cic/nitro/v1/config"
func KubernetesEventsHandler(r *mux.Router){
	
	// swagger:operation  GET /cic/nitro/v1/config/endpoints/{clustername} K8sEndpoints repoList
	// ---
	// summary: This API lists all the endpoints from a given cluster name
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sEndpointPath+"/{clustername}", GetEndpointsAll).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/endpoints/{clustername}/{namespace} K8sEndpoints repoList
	// ---
	// summary: This API lists all the endpoints from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which endpoints are required to be collected
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sEndpointPath+"/{clustername}/{namespace}", GetEndpointsNamespace).Methods("GET")
	// swagger:operation  GET /cic/nitro/v1/config/endpoints/{clustername}/{namespace}/{name} K8sEndpoints repoList
	// ---
	// summary: This API lists all the endpoints from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which endpoints are required to be collected
	//   type: string
	//   required: true	
	// - name: name
	//   in: path
	//   description: name of the endpoint
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sEndpointPath+"/{clustername}/{namespace}/{name}", GetEndpointsName).Methods("GET")
// swagger:operation  GET /cic/nitro/v1/config/services/{clustername} K8sServices repoList
	// ---
	// summary: This API lists all the services from a given cluster name
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sServicePath+"/{clustername}", GetServicesAll).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/services/{clustername}/{namespace} K8sServices repoList
	// ---
	// summary: This API lists all the services from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which services are required to be collected
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sServicePath+"/{clustername}/{namespace}", GetServicesNamespace).Methods("GET")
	// swagger:operation  GET /cic/nitro/v1/config/services/{clustername}/{namespace}/{name} K8sServices repoList
	// ---
	// summary: This API lists all the services from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which services are required to be collected
	//   type: string
	//   required: true	
	// - name: name
	//   in: path
	//   description: name of the endpoint
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sServicePath+"/{clustername}/{namespace}/{name}", GetServicesName).Methods("GET")
		// swagger:operation  GET /cic/nitro/v1/config/pod/{clustername} K8sPods repoList
	// ---
	// suPodmmary: This API lists all the pod from a given cluster name
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sPodPath+"/{clustername}", GetPodsAll).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/pod/{clustername}/{namespace} K8sPods repoList
	// ---
	// summary: This API lists all the pod from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which pod are required to be collected
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sPodPath+"/{clustername}/{namespace}", GetPodsNamespace).Methods("GET")
	// swagger:operation  GET /cic/nitro/v1/config/pod/{clustername}/{namespace}/{name} K8sPods repoList
	// ---
	// summary: This API lists all the pod from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which pod are required to be collected
	//   type: string
	//   required: true	
	// - name: name
	//   in: path
	//   description: name of the endpoint
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sPodPath+"/{clustername}/{namespace}/{name}", GetPodsName).Methods("GET")
	
		// swagger:operation  GET /cic/nitro/v1/config/secret/{clustername} K8sSecrets repoList
	// ---
	// suSecretmmary: This API lists all the secret from a given cluster name
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sSecretPath+"/{clustername}", GetSecretsAll).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/secret/{clustername}/{namespace} K8sSecrets repoList
	// ---
	// summary: This API lists all the secret from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which secret are required to be collected
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sSecretPath+"/{clustername}/{namespace}", GetSecretsNamespace).Methods("GET")
	// swagger:operation  GET /cic/nitro/v1/config/secret/{clustername}/{namespace}/{name} K8sSecrets repoList
	// ---
	// summary: This API lists all the secret from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which secret are required to be collected
	//   type: string
	//   required: true	
	// - name: name
	//   in: path
	//   description: name of the endpoint
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sSecretPath+"/{clustername}/{namespace}/{name}", GetSecretsName).Methods("GET")	
		// swagger:operation  GET /cic/nitro/v1/config/ingress/{clustername} K8sIngresses repoList
	// ---
	// suSecretmmary: This API lists all the ingress from a given cluster name
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sIngressPath+"/{clustername}", GetIngressesAll).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/ingress/{clustername}/{namespace} K8sIngresses repoList
	// ---
	// summary: This API lists all the ingress from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which ingress are required to be collected
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sIngressPath+"/{clustername}/{namespace}", GetIngressesNamespace).Methods("GET")
	// swagger:operation  GET /cic/nitro/v1/config/ingress/{clustername}/{namespace}/{name} K8sIngresses repoList
	// ---
	// summary: This API lists all the ingress from a given cluster name and a namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: username of cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: namespace for which ingress are required to be collected
	//   type: string
	//   required: true	
	// - name: name
	//   in: path
	//   description: name of the endpoint
	//   type: string
	//   required: true	
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(k8sIngressPath+"/{clustername}/{namespace}/{name}", GetIngressesName).Methods("GET")

	
	// swagger:operation  GET /cic/nitro/v1/config/cluster/{clustername}/namespace/{namespace}/service/{service}/ K8sService repoList
	// ---
	// summary: Get the Kubernetes Service details by service name in a specific namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: Name of the cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: provides details in namespace
	//   type: string
	//   required: true
	// - name: service
	//   in: path
	//   description: Provides service details
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(baseURL+"/cluster/{clustername}/namespace/{namespace}/service/{service}", GetService).Methods("GET")
		// swagger:operation  GET /cic/nitro/v1/config/cluster/{clustername}/service/{service}/ K8sService repoList
	// ---
	// summary: Get the Kubernetes Service details by service name in a specific namespace
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: Name of the cluster
	//   type: string
	//   required: true
	// - name: service
	//   in: path
	//   description: Provides service details
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(baseURL+"/cluster/{clustername}/service/{service}", GetService).Methods("GET")
	
	// swagger:operation  GET /cic/nitro/v1/config/cluster/{clustername}/namespace/{namespace} K8sNamespace repoList
	// ---
	// summary: Get the namespace details.
	// description: Test Returns 200 if it success. If there is no cluster registered, Error Not Found (404) will be returned.
	// parameters:
	// - name: clustername
	//   in: path
	//   description: Name of the cluster
	//   type: string
	//   required: true
	// - name: namespace
	//   in: path
	//   description: Provides namespace details
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: successful operation
	//   '204':
	//     description: successful operation, list is empty.
	r.HandleFunc(baseURL+"/cluster/{clustername}/namespace/{namespace}", GetNamespace).Methods("GET")
	
}

func GetEndpoints(resp http.ResponseWriter, req *http.Request){
	fmt.Println("Get Endpoints from the cluster")

	if (clientServerList == nil) {
		fmt.Println("ENDPOINT GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return 
	} 
    	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, "/cic/nitro/v1/config/endpoints/")
        fmt.Println("ENDPOINT: cluster Name input:", clusterName)
	for _, v  := range clientServerList {
       		if (v.ClusterName == clusterName){
			fmt.Println("ENDPOINT GET API: Valid Cluster") 
			message, err := ctr.GetK8sEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "endpoints", "", "")
			if (err != nil) {
				resp.WriteHeader(http.StatusInternalServerError)
				return 
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return 
		} 
    }
	fmt.Println("ENDPOINT GET API: There is no cluster registered with the name",clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}

func GetNamespace(resp http.ResponseWriter, req *http.Request){
	fmt.Println("Get the details of a given namespace")

	if (clientServerList == nil) {
		fmt.Println("NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return 
	} 
    	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, "/cic/nitro/v1/config/cluster/")
	clusterName := strings.Split(urlArgs, "/")[0]
	namespacename := strings.Split(urlArgs, "/namespace/")[1]
        fmt.Println("NAMESPACE: Namespace and CLuster Name:", namespacename, clusterName)
	for _, v := range clientServerList {
       		if (v.ClusterName == clusterName){
			fmt.Println("NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetK8sEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "namespace", "", namespacename)
			if (err != nil) {
				resp.WriteHeader(http.StatusInternalServerError)
				return 
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return 
		} 
    	}
	fmt.Println("NAMESPACE GET API: There is no cluster registered with the name",clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}

func GetService(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get the details of a given service")

	if clientServerList == nil {
		fmt.Println("SERVICE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, "/cic/nitro/v1/config/cluster/")
	var clusterName, namespace, serviceName string
	if strings.Contains(urlArgs, "namespace") {
		urlSplit := strings.Split(urlArgs, "/")
		clusterName = urlSplit[0]
		namespace = urlSplit[2]
		serviceName = urlSplit[4]
	} else {
		clusterName = strings.Split(urlArgs, "/")[0]
		serviceName = strings.Split(urlArgs, "/service/")[1]
		namespace = "default"
	}
	fmt.Println("Service: %s/%s and Cluster Name %s:", namespace, serviceName, clusterName)
	for _,v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("SERVICE GET API: Valid Cluster")
			message, err := ctr.GetK8sEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "service", namespace, serviceName)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("SERVICE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetEndpointsAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get All Endpoints from the cluster")

	if clientServerList == nil {
		fmt.Println("ENDPOINT ALL GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, k8sEndpointPath + "/")
	fmt.Println("ENDPOINT: cluster Name input:", clusterName)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("ENDPOINT ALL GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "endpoints")
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("ENDPOINT ALL GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetEndpointsNamespace(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Endpoints for Namespace from the cluster")

	if clientServerList == nil {
		fmt.Println("ENDPOINT NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sEndpointPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	fmt.Println("ENDPOINT: cluster Name, Namespace input:", clusterName, namespace)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("ENDPOINT NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "endpoints", namespace)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("ENDPOINT NAMESPACE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetEndpointsName(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Endpoints for Name from the cluster")

	if clientServerList == nil {
		fmt.Println("ENDPOINT NAME GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sEndpointPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	name := urlArgsList[2]
	fmt.Println("ENDPOINT: cluster Name, Namespace, Name input:", clusterName, namespace, name)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("ENDPOINT NAME GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "endpoints", namespace, name)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("ENDPOINT NAME GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}

func GetPodsAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get All Pods from the cluster")

	if clientServerList == nil {
		fmt.Println("Pod ALL GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, k8sPodPath + "/")
	fmt.Println("Pod: cluster Name input:", clusterName)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Pod ALL GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "pods")
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Pod ALL GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetPodsNamespace(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Pods for Namespace from the cluster")

	if clientServerList == nil {
		fmt.Println("Pod NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sPodPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	fmt.Println("Pod: cluster Name, Namespace input:", clusterName, namespace)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Pod NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "pods", namespace)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Pod NAMESPACE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetPodsName(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Pods for Name from the cluster")

	if clientServerList == nil {
		fmt.Println("Pod NAME GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sPodPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	name := urlArgsList[2]
	fmt.Println("Pod: cluster Name, Namespace, Name input:", clusterName, namespace, name)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Pod NAME GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "pods", namespace, name)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Pod NAME GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetSecretsAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get All Secrets from the cluster")

	if clientServerList == nil {
		fmt.Println("Secret ALL GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, k8sSecretPath + "/")
	fmt.Println("Secret: cluster Name input:", clusterName)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Secret ALL GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "secrets")
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Secret ALL GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetSecretsNamespace(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Secrets for Namespace from the cluster")

	if clientServerList == nil {
		fmt.Println("Secret NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sSecretPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	fmt.Println("Secret: cluster Name, Namespace input:", clusterName, namespace)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Secret NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "secrets", namespace)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Secret NAMESPACE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetSecretsName(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Secrets for Name from the cluster")

	if clientServerList == nil {
		fmt.Println("Secret NAME GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sSecretPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	name := urlArgsList[2]
	fmt.Println("Secret: cluster Name, Namespace, Name input:", clusterName, namespace, name)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Secret NAME GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "secrets", namespace, name)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Secret NAME GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetServicesAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get All Services from the cluster")

	if clientServerList == nil {
		fmt.Println("Service ALL GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, k8sServicePath + "/")
	fmt.Println("Service: cluster Name input:", clusterName)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Service ALL GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "services")
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Service ALL GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetServicesNamespace(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Services for Namespace from the cluster")

	if clientServerList == nil {
		fmt.Println("Service NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sServicePath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	fmt.Println("Service: cluster Name, Namespace input:", clusterName, namespace)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Service NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "services", namespace)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Service NAMESPACE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetServicesName(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Services for Name from the cluster")

	if clientServerList == nil {
		fmt.Println("Service NAME GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sServicePath+ "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	name := urlArgsList[2]
	fmt.Println("Service: cluster Name, Namespace, Name input:", k8sServicePath, urlArgs, urlArgsList, clusterName, namespace, name)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Service NAME GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "services", namespace, name)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Service NAME GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetIngressesAll(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get All Ingresses from the cluster")

	if clientServerList == nil {
		fmt.Println("Ingress ALL GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	clusterName := strings.TrimPrefix(req.URL.Path, k8sIngressPath + "/")
	fmt.Println("Ingress: cluster Name input:", clusterName)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Ingress ALL GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "ingresses")
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Ingress ALL GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetIngressesNamespace(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Ingresses for Namespace from the cluster")

	if clientServerList == nil {
		fmt.Println("Ingress NAMESPACE GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sIngressPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	fmt.Println("Ingress: cluster Name, Namespace input:", clusterName, namespace)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Ingress NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "ingresses", namespace)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Ingress NAMESPACE GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}
func GetIngressesName(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Get Ingresses for Name from the cluster")

	if clientServerList == nil {
		fmt.Println("Ingress NAME GET API: There is no cluster registered")
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	urlArgs := strings.TrimPrefix(req.URL.Path, k8sIngressPath + "/")
	urlArgsList := strings.Split(urlArgs, "/")
	clusterName := urlArgsList[0]
	namespace := urlArgsList[1]
	name := urlArgsList[2]
	fmt.Println("Ingress: cluster Name, Namespace, Name input:", clusterName, namespace, name)
	for _, v := range clientServerList {
		if v.ClusterName == clusterName {
			fmt.Println("Ingress NAME GET API: Valid Cluster")
			message, err := ctr.GetKubeEvents(v.ConfigFileName, v.KubeURL, v.KubeServAcctToken, "ingresses", namespace, name)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(message)
			return
		}
	}
	fmt.Println("Ingress NAME GET API: There is no cluster registered with the name", clusterName)
	resp.WriteHeader(http.StatusNoContent)
	return
}

func GetAllClientServer(resp http.ResponseWriter, req *http.Request){
	fmt.Println("Get All Client Server")

	if (clientServerList == nil) {
		fmt.Println("List is Empty")
		resp.WriteHeader(http.StatusNoContent)
		return 
	} 
	resp.WriteHeader(http.StatusOK)
    resp.Header().Set("Content-Type", "application/json")
    json.NewEncoder(resp).Encode(clientServerList)
	return

}

func UpdateClientServer(resp http.ResponseWriter, req *http.Request){
	fmt.Println("UPDATE  Client Server")
	if (clientServerList == nil) {
		fmt.Println("Update Operation is only for configured entities")
		resp.WriteHeader(http.StatusNotFound)
		return 
	} 
	newdata := createClientServer{}
	err := json.NewDecoder(req.Body).Decode(&newdata)
	if err != nil {
                fmt.Println("READING ERROR", err)
        }
	for id := range clientServerList {
       		if (clientServerList[id].ClusterName == newdata.ClusterName){
			fmt.Println("Entity is Exist, Updating the entity")
			clientServerList[id] = newdata
			resp.WriteHeader(http.StatusOK)
			return 
		} 
    	}
	fmt.Println("Update Operation is only for configured entities")
	resp.WriteHeader(http.StatusNotFound)
	return 
}

func DeleteClientServer(resp http.ResponseWriter, req *http.Request){
	fmt.Println("Delete Client Server\n", len(clientServerList))
	if (clientServerList == nil) {
		fmt.Println("Entity is not configured")
		resp.WriteHeader(http.StatusNotFound)
		return 
	} 
	newdata := createClientServer{}
	err := json.NewDecoder(req.Body).Decode(&newdata)
	if err != nil {
                fmt.Println("READING ERROR", err)
        }
	for id := range clientServerList {
       		if (clientServerList[id].ClusterName == newdata.ClusterName){
			fmt.Println("Entity is Exist, Removing the etity", newdata.ClusterName)
			clientServerList = append(clientServerList[:id], clientServerList[id+1:]...)
			resp.WriteHeader(http.StatusOK)
			return 
		} 
    	}
	fmt.Printf("Dump Complete List")
	spew.Dump(clientServerList)
	fmt.Println("Entity is not configured")
	resp.WriteHeader(http.StatusNotFound)
	return 
}


func PostClientServer(resp http.ResponseWriter, req *http.Request){
	fmt.Println("Post Client Server")
	newdata := createClientServer{}
	err := json.NewDecoder(req.Body).Decode(&newdata)
	if err != nil {
                fmt.Println("READING ERROR", err)
        }
    	fmt.Printf("DECODER Results: %v\n", newdata)
	for id := range clientServerList {
       		if (clientServerList[id].ClusterName == newdata.ClusterName){
			response := "Entity is Exist, please use update API for updating Endpoint/Server details"
			fmt.Println(response)
        		http.Error(resp, response, http.StatusNotModified)
			return 
		} 
    	}
	
	fmt.Printf("Dump Complete Struture=%v", clientServerList)
	spew.Dump(clientServerList)
	status, statusString := ctr.StartController(newdata.ConfigFileName, newdata.KubeURL, newdata.KubeServAcctToken,
		newdata.ServerURL, newdata.WatchEvents)
	if status == http.StatusOK {
		clientServerList = append(clientServerList, newdata)
	}
	resp.WriteHeader(status)
	resp.Write([]byte(statusString))
	return
}


	

func StartRestServer() (*http.Server){
	fmt.Println("Staring the REST Server")
        // Create Server and Route Handlers
        r := mux.NewRouter()
	dir, err := os.Getwd()
	if err != nil {
		klog.Error(err)
	}

	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(dir+"/swaggerui"))))
	

	CreateClientServerHandler(r)
        KubernetesEventsHandler(r)	
       	 
	srv := &http.Server{
                Handler:      r,
                Addr:         "localhost:8000",
                ReadTimeout:  10 * time.Second,
                WriteTimeout: 10 * time.Second,
        }

        // Start Server
        go func() {
                klog.Info("Starting MultiCluster Ingress Rest Interface at", srv.Addr)
                if err := srv.ListenAndServe(); err != nil {
                        klog.Fatal(err)
                }
        }()
        // Graceful Shutdown
	return srv
}
