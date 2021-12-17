package sys

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
func VaultUnseal(vaultAddress, vaultPort, MasterKey, caChainFile, clientCertFile, clientKeyFile string) vaultUnsealResp {

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)

	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(caChainFile)

	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}

	var data vaultUnsealPayload
	var output vaultUnsealResp

	data.Key = MasterKey

	dataJSON, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://%s:%s/v1/sys/unseal", vaultAddress, vaultPort), bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
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
