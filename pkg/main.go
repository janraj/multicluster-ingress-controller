// Package main Multicluster Ingress Controller API.
//
// Multicluster Ingress Controller is capable of proccessing ingress events from multicluster.
// It can send kubernetes events to multiple servers and have control of event selection.
// 
//
// 
//
//     Schemes: http
//     Host: localhost:8000
//     BasePath: /
//     Version: 1.0.0
//     Contact: Janraj CJ<janrajcj@gmail.com> 
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main
  
import (
        "log"
	intr "multicluster-ingress-controller/pkg/interface"
)

func InitCitrixControlPlane() error {
        log.Println("Initializing MultiCluster Ingress Controller ....")
        return nil
}
func main() {
        InitCitrixControlPlane()
	intr.InitClientServer()
	intr.StartRestServer()
	select {}
}
