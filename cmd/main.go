package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	cf "github.com/blinchik/aws/config"
	secret "github.com/blinchik/aws/services/secretmanager"
	sys "github.com/blinchik/vault/sys"
)

var consulVault string

const (
	vaultPort = "8200"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	cfg, region := cf.CreateConfigFromEC2Role()

	init := flag.Bool("init", false, "init vault")
	enable := flag.Bool("enable", false, "enable secret engine")

	flag.Parse()

	if *init {

		vaultAddress := os.Args[2]
		vaultPort := os.Args[3]

		output := sys.VaultInit(vaultAddress, vaultPort, 5, 3)

		secret.CreateSecret("brain-vault_root_token", output.RootToken, "vault root token (should be not saved in future)", cfg, region)
		fmt.Println(output)

		for idx, key := range output.Keys {

			secret.CreateSecret(fmt.Sprintf("brain-vault_backup_key_%d", idx), key, "vault backup key", cfg, region)
			sys.VaultUnseal(vaultAddress, vaultPort, key)

		}

		if *enable {

			secretEngines := strings.Split(os.Args[2], ",")

			for _, se := range secretEngines {

				sys.EnableSecretEngine(consulVault, vaultPort, &output.RootToken, se)

			}

		}

	}

}
