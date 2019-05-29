package rest
import( 
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"bytes"
        "encoding/json"
	"fmt"
	"os"
        "path/filepath"
)
func TestInitRestInterface(t *testing.T) {
	Convey("Initializing Citrix Control Plane", t, func(){
		InitClientServer()
	})
}
func TestRestInterface(t *testing.T) {
	StartRestServer()
	data := createClientServer{}
	data.ClusterName = "Cluster-1" 
	data.ConfigFileName = "/home/.kube/config"
	data.ServerURL = append(data.ServerURL, "22.22.22.22")
	data.ServerURL = append(data.ServerURL, "32.32.32.32")
	data.WatchEvents = append(data.WatchEvents, "ingresses")
	data.WatchEvents = append(data.WatchEvents,  "endpoints")
	dataJson, _ := json.Marshal(data)

	Convey("Citrix Control Plane Rest Interface", t, func(){
		Convey("When Citrix Control Plane is running", func() {

			Convey("invalid Get Rest API should retun 404, Server Not Found Error", func() {
				resp, err := http.Get("http://localhost:8000/"+routePath+"test")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
				}
			})
			Convey("Get Rest API should retun 204, if there is no registeration yet", func() {
				resp, err := http.Get("http://localhost:8000/"+routePath)
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 204)
				}
			})
			Convey("Put Rest API should retun 404, Status Not Found, if there is no entity", func() {
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "PUT")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
					fmt.Print(resp)
				}
			})
			Convey("Delete Rest API should retun 404, Status Not Found, if there is no entity", func() {
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "DELETE")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
					fmt.Print(resp)
				}
			})
			Convey("invalid Post Rest API should retun 404, Server Not Found Error", func() {
				resp, err := http.Post("http://localhost:8000/"+routePath+"test", "application/json", bytes.NewBuffer(dataJson))
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
				}
			})
			Convey("Post Rest API should retun 200, Success", func() {
				resp, err := http.Post("http://localhost:8000"+routePath, "application/json", bytes.NewBuffer(dataJson))
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
			})
			Convey("Get Rest API should retun 200, Success, after a POST", func() {
				resp, err := http.Get("http://localhost:8000"+routePath)
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
			})
			Convey("Post Rest API should retun 304, StatusNotModified, when trying to add already existing entity", func() {
				resp, err := http.Post("http://localhost:8000"+routePath, "application/json", bytes.NewBuffer(dataJson))
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 304)
				}
			})
			Convey("Put Rest API should retun 200, Success, when updating the entity", func() {
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "PUT")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
			})
			Convey("Put Rest API should retun 404, StatusNotFound, when updating non existing entity", func() {
				data.ClusterName = "Cluster-2"
				dataJson, _ = json.Marshal(data)
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "PUT")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
					fmt.Print(resp)
				}
				data.ClusterName = "Cluster-1"
				dataJson, _ = json.Marshal(data)
			})
			Convey("Delete Rest API should retun 404, StatusNotFound, when try to delete non existing entity", func() {
				data.ClusterName = "Cluster-2"
				dataJson, _ = json.Marshal(data)
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "DELETE")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 404)
					fmt.Print(resp)
				}
				data.ClusterName = "Cluster-1"
				dataJson, _ = json.Marshal(data)
			})
			Convey("Delete Rest API should retun 200, Success, entity Must removed", func() {
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "DELETE")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
			})
			Convey("Another Input: Post Rest API should retun 200, Success", func() {
				data.ClusterName = "Cluster-2"
				data.ConfigFileName = filepath.Join(os.Getenv("HOME"), ".kube", "config")
				data.WatchEvents = append(data.WatchEvents, "pods")		
				data.WatchEvents = append(data.WatchEvents, "services")		
				data.WatchEvents = append(data.WatchEvents, "endpoints")		
				dataJson, _ = json.Marshal(data)
				resp, err := http.Post("http://localhost:8000"+routePath, "application/json", bytes.NewBuffer(dataJson))
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
				data.ClusterName = "Cluster-1"
				data.ConfigFileName = "/home/.kube/config"
				data.WatchEvents = data.WatchEvents[:1]
				dataJson, _ = json.Marshal(data)
			})
			Convey("Another Input: Delete Rest API should retun 200, Success, entity Must removed", func() {
				data.ClusterName = "Cluster-2"
				data.ConfigFileName = filepath.Join(os.Getenv("HOME"), ".kube", "config")
				data.WatchEvents = append(data.WatchEvents, "pods")		
				data.WatchEvents = append(data.WatchEvents, "services")		
				data.WatchEvents = append(data.WatchEvents, "endpoints")		
				dataJson, _ = json.Marshal(data)
				resp, err := httpHelper("http://localhost:8000"+routePath, bytes.NewBuffer(dataJson), "DELETE")
				if (err == nil) {
					So(resp.StatusCode, ShouldEqual, 200)
					fmt.Print(resp)
				}
				data.ClusterName = "Cluster-1"
				data.ConfigFileName = "/home/.kube/config"
				data.WatchEvents = data.WatchEvents[:1]
				dataJson, _ = json.Marshal(data)
			})
		})
	
	})
}

func httpHelper(uri string,  body *bytes.Buffer, op string)(*http.Response, error){
	client := &http.Client{}
	var err error
	var req *http.Request
	if (op == "PUT") {
		req, err = http.NewRequest(http.MethodPut, uri, body)
	}else if (op == "DELETE"){
		req, err = http.NewRequest(http.MethodDelete, uri, body)
		
	}else {
		req, err = http.NewRequest(http.MethodGet, uri, body)
	}
	
	if err != nil {
		fmt.Printf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do() failed with '%s'\n", err)
	}
	return resp, err
}
