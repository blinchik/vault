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
	"os"
	"strings"
	"time"
)

type vaultInitPayload struct {
	SecretShares      int `json:"secret_shares"`
	SecretThreshold   int `json:"secret_threshold"`
	RecoveryShares    int `json:"recovery_shares"`
	RecoveryThreshold int `json:"recovery_threshold"`
}

type vaultInitResp struct {
	Keys               []string `json:"keys"`
	KeysBase64         []string `json:"keys_base64"`
	RecoveryKeys       []string `json:"recovery_keys"`
	RecoveryKeysBase64 []string `json:"recovery_keys_base64"`
	RootToken          string   `json:"root_token"`
}

// VaultInit This endpoint initializes a new Vault. The Vault must not have been previously initialized.
// The recovery options, as well as the stored shares option, are only available when using Vault HSM.
// https://www.vaultproject.io/api-docs/system/init
func VaultInit(vaultAddress, vaultPort string, shares, threshold int, caChainFile, clientCertFile, clientKeyFile string) vaultInitResp {

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

	var data vaultInitPayload
	var output vaultInitResp

	/*
		Auto-Unseal uses the recovery key options for initialization.
		Therefore, "recovery_shares" and "recovery_threshold" should be used
	*/

	data.SecretShares = 5
	data.SecretThreshold = 3

	data.RecoveryShares = shares
	data.RecoveryThreshold = threshold

	dataJSON, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://%s:%s/v1/sys/init", vaultAddress, vaultPort), bytes.NewBuffer(dataJSON))
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

	if strings.Contains(string(bodyBytes), "Vault is already initialized") {

		os.Exit(1)
	}

	return output

}
