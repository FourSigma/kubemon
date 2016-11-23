package service

import (
	"io"

	"github.com/FourSigma/kubemon/models"

	v1core "k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	"k8s.io/client-go/1.4/pkg/api"
)

type podService struct {
	v1core.PodInterface
}

func (p *podService) DeletePod(podName string) error {
	return p.Delete(podName, api.ListOptions{})
}

func (p *podService) ListPods() (rs []models.Pod, err error) {
	pList, err := p.List(api.ListOptions{})
	if err != nil {
		return
	}
	for _, v := range pList.Items {

	}
	return
}

func (p *podService) PodLog(podName string) (r io.Reader, err error) {
	return
}
