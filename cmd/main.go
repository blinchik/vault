package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	cf "github.com/blinchik/aws/config"
	ec2 "github.com/blinchik/aws/services/ec2"
	secret "github.com/blinchik/aws/services/secretmanager"
	sys "github.com/blinchik/vault/sys"
)

const (
	vaultPort = "8200"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	cfg, region := cf.CreateConfigFromEC2Role()

	init := flag.Bool("init", false, "init vault")
	test := flag.Bool("test", false, "init vault")

	flag.Parse()

	if *test {
		q := ec2.DescribeInstanceByTag("Function", "consul_client", "Group", "brainstorming", cfg, region)

		fmt.Printf(q)
	}
	if *init {

		vaultAddress := os.Args[2]
		vaultPort := os.Args[3]
		tokenName := os.Args[4]

		output := sys.VaultInit(vaultAddress, vaultPort, 1, 1)

		secret.CreateSecret(fmt.Sprintf("brain-vault_root_token_%s", tokenName), output.RootToken, "vault root token (should be not saved in future)", cfg, region)

		for idx, key := range output.Keys {

			secret.CreateSecret(fmt.Sprintf("brain-vault_backup_key_%d_%s", idx, tokenName), key, "vault backup key", cfg, region)
			sys.VaultUnseal(vaultAddress, vaultPort, key)

		}

	}

}
