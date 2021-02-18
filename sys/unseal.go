package sys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type vaultUnsealPayload struct {
	Key string `json:"key"`
}

type vaultUnsealResp struct {
	Sealed      bool   `json:"sealed"`
	T           int    `json:"t"`
	N           int    `json:"n"`
	Progress    int    `json:"progress"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

//VaultUnseal
func VaultUnseal(vaultAddress, vaultPort, MasterKey string) vaultUnsealResp {

	var data vaultUnsealPayload
	var output vaultUnsealResp

	data.Key = MasterKey

	dataJSON, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%s/v1/sys/unseal", vaultAddress, vaultPort), bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, &output)

	if err != nil {
		log.Fatal(err)
	}

	return output

}
