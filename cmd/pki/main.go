package main

import (
	"flag"
	"log"
	"os"
	"strings"

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

	setURL := flag.Bool("setURL", false, "setURL")

	flag.Parse()

	if *enable {

		secretEngine := os.Args[5]
		path := os.Args[6]

		sys.EnableSecretEngine(vaultAddress, vaultPort, &vaultToken, secretEngine, path)

	}

	if *gr {

		commonName := os.Args[5]

		secret.GenerateRoot(vaultAddress, vaultPort, &vaultToken, commonName)

	}

	if *gi {

		commonName := os.Args[5]
		path := os.Args[6]

		cert := secret.GenerateIntermediate(vaultAddress, vaultPort, &vaultToken, commonName, path)

		signedCert := secret.SignIntermediate(vaultAddress, vaultPort, &vaultToken, commonName, cert, "pki_root1")

		q := strings.Split(signedCert, "\n")

		input := strings.Join(q, "\\n")
		// fmt.Println(input)

		secret.SetSignedIntermediate(vaultAddress, vaultPort, &vaultToken, input, path)

	}

	if *setURL {

		crlDistributionPoints := strings.Split(os.Args[5], ",")
		IssuingCertificates := strings.Split(os.Args[6], ",")

		secret.SetURLs(vaultAddress, vaultPort, &vaultToken, crlDistributionPoints, IssuingCertificates)

	}

}
