package secret

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GenerateRoot(vaultAddress string, vaultPort string, vaultToken *string, commonName string) {

	payload := fmt.Sprintf(` {"common_name": "%s" }`, commonName)

	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/pki/root/generate/internal", vaultAddress, vaultPort), body)

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

type outputGenerateIntermediate struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Csr string `json:"csr"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

func GenerateIntermediate(vaultAddress string, vaultPort string, vaultToken *string, commonName, path string) string {

	var output outputGenerateIntermediate

	payload := fmt.Sprintf(` {"common_name": "%s" }`, commonName)

	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/%s/intermediate/generate/internal", vaultAddress, vaultPort, path), body)

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

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		log.Fatal(err)
	}

	return string(output.Data.Csr)

}

type PayloadSignIntermediate struct {
	Csr        string `json:"csr"`
	CommonName string `json:"common_name"`
	Format     string `json:"format"`
}

type outputSignIntermediate struct {
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Certificate  string   `json:"certificate"`
		IssuingCa    string   `json:"issuing_ca"`
		CaChain      []string `json:"ca_chain"`
		SerialNumber string   `json:"serial_number"`
	} `json:"data"`
	Auth interface{} `json:"auth"`
}

func SignIntermediate(vaultAddress string, vaultPort string, vaultToken *string, commonName, csr, path string) string {

	var SIdata PayloadSignIntermediate
	var output outputSignIntermediate

	SIdata.Csr = csr
	SIdata.CommonName = commonName
	SIdata.Format = "pem"

	dataJSON, err := json.Marshal(SIdata)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/%s/root/sign-intermediate", vaultAddress, vaultPort, path), bytes.NewBuffer(dataJSON))

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

	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		log.Fatal(err)
	}

	return output.Data.Certificate

}

func SetSignedIntermediate(vaultAddress string, vaultPort string, vaultToken *string, certInt, path string) {

	payload := fmt.Sprintf(` {"certificate": "%s" }`, certInt)

	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/%s/intermediate/set-signed", vaultAddress, vaultPort, path), body)

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

type setURLPayload struct {
	IssuingCertificates   []string `json:"issuing_certificates"`
	CrlDistributionPoints []string `json:"crl_distribution_points"`
}

func SetURLs(vaultAddress string, vaultPort string, vaultToken *string, CrlDistributionPoints, IssuingCertificates []string) {

	var urlData setURLPayload

	urlData.CrlDistributionPoints = CrlDistributionPoints
	urlData.IssuingCertificates = IssuingCertificates

	dataJSON, err := json.Marshal(urlData)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/v1/pki/config/urls", vaultAddress, vaultPort), bytes.NewBuffer(dataJSON))

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
