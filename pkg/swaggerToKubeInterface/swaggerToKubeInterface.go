package swaggerToKubeInterface

import (
	"errors"
	"fmt"
	"multicluster-ingress-controller/pkg/kubeController"
	"multicluster-ingress-controller/pkg/swaggerInterface"
)

type SwaggerToKubeInterface struct {
	SwaggerClientServer *swaggerInterface.CreateClientServer
	KController         *kubeController.KubeController
}

func (sTKI *SwaggerToKubeInterface) DeleteAll() {
	sTKI.KController.DeleteAll()
	sTKI.KController = nil
	sTKI.SwaggerClientServer = nil
}
func (sTKI *SwaggerToKubeInterface) createLinkswgIntfToKubeController() error {

	swgClientServer := *sTKI.SwaggerClientServer

	err := swgClientServer.ValidateKubeClusterFields()
	if err != nil {
		return err
	}

	api, err := swgClientServer.CreateK8sApiserverClient()
	if err != nil {
		fmt.Println("Error while starting client API session")
		return errors.New("Error while starting client API session: error: " + err.Error())
	}

	kC := kubeController.NewkubeController()
	err = kC.PopulateKubeController(api, &swgClientServer)
	if err != nil {
		return errors.New("FAILURE: Kube Controller Add Fail: Error: " + err.Error())
	}

	(*sTKI).KController = kC
	return nil

}

func (sTKI *SwaggerToKubeInterface) LinkAndStartWatchEvents() error {
	err := sTKI.createLinkswgIntfToKubeController()

	if err != nil {
		return err
	}
	kubeController := *sTKI.KController
	err = kubeController.RunWatchEvents()

	if err != nil {
		return errors.New("Invalid Controller Input: Error :" + err.Error())
	}
	return nil

}
