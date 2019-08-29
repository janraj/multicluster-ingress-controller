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
var k8sEndpointPath string = "/cic/nitro/v1/config/endpoints"
var baseURL string = "/cic/nitro/v1/config"
func KubernetesEventsHandler(r *mux.Router){
	
	// swagger:operation  GET /cic/nitro/v1/config/endpoints/{clustername} K8sEndpoints repoList
	// ---
	// summary: This API lists all the endpoints from a given cluster name.
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
	r.HandleFunc(k8sEndpointPath+"/{clustername}", GetEndpoints).Methods("GET")
	
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
	for id := range clientServerList {
       		if (clientServerList[id].ClusterName == clusterName){
			fmt.Println("ENDPOINT GET API: Valid Cluster")
			message, err := ctr.GetK8sEvents(clientServerList[id].ConfigFileName, "endpoints", "", "")
			if (err != nil) {
				resp.WriteHeader(http.StatusInternalServerError)
				return 
			}
			json.NewEncoder(resp).Encode(message)
			resp.WriteHeader(http.StatusOK)
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
	for id := range clientServerList {
       		if (clientServerList[id].ClusterName == clusterName){
			fmt.Println("NAMESPACE GET API: Valid Cluster")
			message, err := ctr.GetK8sEvents(clientServerList[id].ConfigFileName, "namespace", "", namespacename)
			if (err != nil) {
				resp.WriteHeader(http.StatusInternalServerError)
				return 
			}
			json.NewEncoder(resp).Encode(message)
			resp.WriteHeader(http.StatusOK)
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
			message, err := ctr.GetK8sEvents(v.ConfigFileName, "service", namespace, serviceName)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(resp).Encode(message)
			resp.WriteHeader(http.StatusOK)
			return
		}
	}
	fmt.Println("SERVICE GET API: There is no cluster registered with the name", clusterName)
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
    	resp.Header().Set("Content-Type", "application/json")
    	json.NewEncoder(resp).Encode(clientServerList)
	resp.WriteHeader(http.StatusOK)
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
	clientServerList = append(clientServerList, newdata)
	fmt.Printf("Dump Complete Struture=%v", clientServerList)
	spew.Dump(clientServerList)
	resp.WriteHeader(http.StatusOK)
        ctr.StartController(newdata.ConfigFileName, newdata.ServerURL, newdata.WatchEvents)	
	return
}


	

func StartRestServer() (*http.Server){
	platform := os.Getenv("PLATFORM")
	address := "localhost"
	fmt.Println("Staring the REST Server")
        // Create Server and Route Handlers
        r := mux.NewRouter()
	dir, err := os.Getwd()
	if err != nil {
		klog.Error(err)
	}
	if (platform == "Container"){
		address = "0.0.0.0"
	}
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(dir+"/swaggerui"))))
	CreateClientServerHandler(r)
        KubernetesEventsHandler(r)	
	srv := &http.Server{
                Handler:      r,
                Addr:         address+":8000",
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
