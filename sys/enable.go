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
	Type   string `json:"type"`
	Path   string `json:"path"`
	Config struct {
		DefaultLeaseTtl string `json:"default_lease_ttl"`
		MaxLeaseTtl     string `json:"max_lease_ttl"`
	} `json:"config"`
}

func EnableSecretEngine(vaultAddress string, vaultPort string, vaultToken *string, secretEngine, path, defaultLeaseTtl, maxLeaseTtl string) {

	var enableData EnablePayload

	enableData.Path = path
	enableData.Type = secretEngine
	enableData.Config.DefaultLeaseTtl = defaultLeaseTtl
	enableData.Config.MaxLeaseTtl = maxLeaseTtl

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
