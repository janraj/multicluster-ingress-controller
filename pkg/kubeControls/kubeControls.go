/* kubeControls.go

   This file contains kubernetes server details datastructure
*/
package kubeControls

import (
	"errors"
	"k8s.io/api/core/v1"
	"multicluster-ingress-controller/pkg/namespaces"
	"multicluster-ingress-controller/pkg/swaggerInterface"
	"strconv"
)

type KubeControl struct {
	controlList []namespaces.Namespace
}

func New() *KubeControl {
	return &KubeControl{}
}

func (kC *KubeControl) Namespaceslength() int {
	return len(kC.controlList)
}

func (kC *KubeControl) GetNamespace(index int) (*namespaces.Namespace, error) {
	length := len(kC.controlList)
	if index >= length {
		return nil, errors.New("KubeControl: Index out of bounds Index: " + strconv.Itoa(index) + " Length: " + strconv.Itoa(length))
	}
	return &kC.controlList[index], nil
}

func (kC *KubeControl) PopulateKubeControls(cS swaggerInterface.CreateClientServer) error {
	err := kC.AddNamespaces(cS)
	if err != nil {
		return err
	}
	return nil
}

func (kC *KubeControl) AddNamespaces(cS swaggerInterface.CreateClientServer) error {
	var controlList []namespaces.Namespace
	namespaceList := cS.Namespaces
	if len(namespaceList) == 0 {
		namespaceList = append(namespaceList, v1.NamespaceAll)
	}
	for _, v := range namespaceList {
		nS := namespaces.New(v)
		err := nS.AddEvents(cS)
		if err != nil {
			return err
		}
		controlList = append(controlList, *nS)
	}
	kC.controlList = controlList
	return nil

}
func (kC *KubeControl) DeleteAll() {
	for _, v := range kC.controlList {
		v.DeleteAll()
	}
	kC.controlList = nil
}
