package main

import (
	"log"
	"os"

	mongodb "git.drom.ru/security/vault-extended-plugin-database-mongodb"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func main() {
	if err := Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run starts the RPC server for the plugin.
func Run() error {
	dbplugin.ServeMultiplex(mongodb.New)
	return nil
}
