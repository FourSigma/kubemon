package models

import (
	"encoding/json"
	"fmt"

	"k8s.io/client-go/1.4/pkg/api/v1"
)

type PodList []Pod

func ToPodList(list ...v1.Pod) (rs PodList) {
	for _, v := range list {
		tmp := &v
		rs = append(rs, ToPod(tmp))
	}
	return
}

func ToPod(p *v1.Pod) Pod {
	return Pod{
		Pod: p,
	}
}

type Pod struct {
	*v1.Pod
}

func (p *Pod) MarshalJSON() (b []byte, err error) {

	s := struct {
		Name      string      `json:"name"`
		Status    v1.PodPhase `json:"status"`
		CreatedAt string      `json:"createdAt"`
	}{
		Name:      p.GetName(),
		Status:    p.Status.Phase,
		CreatedAt: p.GetCreationTimestamp().String(),
	}

	fmt.Println(p.GetName())

	return json.Marshal(s)
}
