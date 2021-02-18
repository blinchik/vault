package sys

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func EnableSecretEngine(vaultAddress string, vaultPort string, vaultToken *string, secretEngine string) {

	payload := fmt.Sprintf(` {"type": "%s" }`, secretEngine)

	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/sys/mounts/%s", vaultAddress, vaultPort, secretEngine), body)

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
