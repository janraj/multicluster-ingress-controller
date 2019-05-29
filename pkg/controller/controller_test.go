package controller
import(
	"github.com/gorilla/mux"
	"time"
	"io/ioutil"
	"encoding/json"
	"net/http"
	//"log"
	 
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
        //"k8s.io/client-go/kubernetes/fake"
        //cachefake "k8s.io/client-go/tools/cache/testing"
)
/*
func getK8sClient() (*KubernetesAPIServer) {
        fake := fake.NewSimpleClientset()
        api := &KubernetesAPIServer{
                Suffix: "Test",
                Client: fake,
        }
        return api
}
*/

func ServerHandler(w http.ResponseWriter, r *http.Request) {
        b, err := ioutil.ReadAll(r.Body)
        defer r.Body.Close()
        if err != nil {
                http.Error(w, err.Error(), 500)
                return
        }
        buf := make(map[string]string)
        err = json.Unmarshal(b, &buf)
        if err != nil {
                http.Error(w, err.Error(), 500)
                return
        }
        fmt.Printf("Test Server: Resource Type = %s\n", buf["resource_type"])
}



func StartDummyServer(){
	r := mux.NewRouter()
        r.HandleFunc("/", ServerHandler)
        srv := &http.Server{
                Handler:      r,
                Addr:         "127.0.0.1:9000",
                ReadTimeout:  10 * time.Second,
                WriteTimeout: 10 * time.Second,
        }

        // Start Server
        go func() {
                //log.Info("Starting Server")
                if err := srv.ListenAndServe(); err != nil {
                       //log.Fatal(err)
                }
        }()
}

func TestGenerateUUID(t *testing.T) {
	Convey("Citrix Control Plane, UUID Generation", t, func(){
		uuid1 := GenerateUUID()
		uuid2 := GenerateUUID()
		So(uuid1, ShouldNotEqual, uuid2)
		uuid3 := GenerateUUID()
		uuid4 := GenerateUUID()
		So(uuid3, ShouldNotEqual, uuid4)
	})
}
func TestController(t *testing.T) {
	//api := getK8sClient()
	StartDummyServer()
	//StartDummyK8sClient()
	Convey("Controller Package Validation", t, func(){
		Convey("Starting Controller for event ingresses which should send to localhost:9000", func(){
			var serverList []string
			serverList = append(serverList, "http://127.0.0.1:9000/")
			var eventList []string 
			eventList = append(eventList, "ingresses")
			eventList = append(eventList, "endpoints")
			eventList = append(eventList, "pods")
			eventList = append(eventList, "services")
			StartController("/Users/janraj/.kube/config", serverList, eventList)		
			time.Sleep(3000 * time.Millisecond)	
		})
	})
}
