package main

import (
	"flag"
	"log"
	"strings"

	"github.com/trusch/vpntool/pki"
)

var pkiDir = flag.String("pki", "pki", "pki directory to operate in")
var initCA = flag.Bool("init", false, "init a new CA in --pki")
var addClient = flag.String("add-client", "", "client to create")
var addServer = flag.String("add-server", "", "server to create")
var createDH = flag.Bool("create-dh", false, "create Diffi Hellman parameters")
var revoke = flag.String("revoke", "", "clients/servers to revoke")

func main() {
	flag.Parse()
	if *initCA {
		if err := pki.Init(*pkiDir); err != nil {
			log.Fatal(err)
		}
	}
	if *addClient != "" {
		if strings.Index(*addClient, ",") > -1 {
			clients := strings.Split(*addClient, ",")
			for _, client := range clients {

				if err := pki.AddClient(*pkiDir, client); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			if err := pki.AddClient(*pkiDir, *addClient); err != nil {
				log.Fatal(err)
			}
		}
	}
	if *addServer != "" {
		if strings.Index(*addServer, ",") > -1 {
			servers := strings.Split(*addServer, ",")
			for _, server := range servers {
				if err := pki.AddServer(*pkiDir, server); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			if err := pki.AddServer(*pkiDir, *addServer); err != nil {
				log.Fatal(err)
			}
		}
	}
	if *createDH {
		if err := pki.CreateDH(*pkiDir); err != nil {
			log.Fatal(err)
		}
	}
	if *revoke != "" {
		if strings.Index(*revoke, ",") > -1 {
			entities := strings.Split(*revoke, ",")
			for _, entity := range entities {
				if err := pki.Revoke(*pkiDir, entity); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			if err := pki.Revoke(*pkiDir, *revoke); err != nil {
				log.Fatal(err)
			}
		}
	}
}
