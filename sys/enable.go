package sys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type EnablePayload struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func EnableSecretEngine(vaultAddress string, vaultPort string, vaultToken *string, secretEngine, path string) {

	var enableData EnablePayload

	enableData.Path = path
	enableData.Type = secretEngine

	dataJSON, err := json.Marshal(enableData)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/sys/mounts/%s", vaultAddress, vaultPort, path), bytes.NewBuffer(dataJSON))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Vault-Token", *vaultToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bodyBytes))

}

// type mountConfig struct {
// 	MaxLeaseTTL     int `json:"max_lease_ttl"`
// }

// func TuneSecretEngine(vaultAddress string, vaultPort string, vaultToken *string, secretEngine string) {

// 	var mount mountConfig

// 	mount.MaxLeaseTTL =

// 	payload := fmt.Sprintf(` {"type": "%s" }`, secretEngine)

// 	body := strings.NewReader(payload)
// 	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/sys/mounts/%s/tune", vaultAddress, vaultPort, secretEngine), body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req.Header.Set("X-Vault-Token", *vaultToken)
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println(string(bodyBytes))

// }
