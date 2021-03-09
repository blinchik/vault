package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"strconv"

	secret "github.com/blinchik/vault/secret"
	sys "github.com/blinchik/vault/sys"
)

const (
	vaultPort = "8200"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	vaultAddress := os.Args[2]
	vaultPort := os.Args[3]
	vaultToken := os.Args[4]

	enable := flag.Bool("enable", false, "enable secret engine")
	gr := flag.Bool("gr", false, "generate root CA")
	gi := flag.Bool("gi", false, "generate intermidate CA")
	cr := flag.Bool("cr", false, "create pki role")

	setURL := flag.Bool("setURL", false, "setURL")

	flag.Parse()

	if *enable {

		secretEngine := os.Args[5]
		path := os.Args[6]
		defaultLeaseTtl := os.Args[7]
		maxLeaseTtl := os.Args[8]


		sys.EnableSecretEngine(vaultAddress, vaultPort, &vaultToken, secretEngine, path, defaultLeaseTtl, maxLeaseTtl)

	}

	if *gr {

		commonName := os.Args[5]
		ttl := os.Args[6]

		secret.GenerateRoot(vaultAddress, vaultPort, &vaultToken, commonName, ttl)

	}

	if *gi {

		commonName := os.Args[5]
		path := os.Args[6]
		ttl := os.Args[7]


		cert := secret.GenerateIntermediate(vaultAddress, vaultPort, &vaultToken, commonName, path)

		signedCert := secret.SignIntermediate(vaultAddress, vaultPort, &vaultToken, commonName, cert, "pki", ttl)

		q := strings.Split(signedCert, "\n")

		input := strings.Join(q, "\\n")


		secret.SetSignedIntermediate(vaultAddress, vaultPort, &vaultToken, input, path)

	}

	if *setURL {

		crlDistributionPoints := strings.Split(os.Args[5], ",")
		IssuingCertificates := strings.Split(os.Args[6], ",")

		secret.SetURLs(vaultAddress, vaultPort, &vaultToken, crlDistributionPoints, IssuingCertificates)

	}


	
	if *cr {

		path := os.Args[5]
		roleName := os.Args[6]


		allow_subdomains, err := strconv.ParseBool(os.Args[7])

		if err != nil {
			log.Fatal(err)
		}

		allowed_domains := strings.Split(os.Args[8], ",")

		maxTtl := os.Args[9]
		ttl := os.Args[10]


		secret.CreateRole(vaultAddress, vaultPort, &vaultToken, path, roleName, allow_subdomains, allowed_domains, maxTtl, ttl)

	}

}
