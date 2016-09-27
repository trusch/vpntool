package main

import (
	"flag"
	"log"
	"strings"

	"github.com/trusch/vpntool/openvpn"
)

var pkiDir = flag.String("pki", "pki", "pki directory")
var dir = flag.String("out", ".", "ovpn directory")
var initVPN = flag.Bool("init", false, "init vpn and create server")
var addClient = flag.String("clients", "", "add client(s) to vpn (accepts comma separated list)")
var deploy = flag.String("deploy", "", "deploy this entity to --url")
var url = flag.String("url", "", "url to use")
var peerToPeer = flag.Bool("peer-to-peer", true, "enable client to client communication")
var revoke = flag.String("revoke", "", "revoke this client")

func main() {
	flag.Parse()
	if *initVPN {
		if err := openvpn.Init(*pkiDir, *dir, *peerToPeer); err != nil {
			log.Fatal(err)
		}
	}
	if *addClient != "" {
		if *url == "" {
			log.Fatal("specify --url to point to your VPN server")
		}
		if strings.Index(*addClient, ",") > -1 {
			clients := strings.Split(*addClient, ",")
			for _, client := range clients {
				if err := openvpn.CreateClient(*pkiDir, client, *url, *dir); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			if err := openvpn.CreateClient(*pkiDir, *addClient, *url, *dir); err != nil {
				log.Fatal(err)
			}
		}
	}
	if *deploy != "" {
		if *url == "" {
			log.Fatal("specify --url to point to your deploy target")
		}
		if err := openvpn.Deploy(*dir, *pkiDir, *deploy, *url); err != nil {
			log.Fatal(err)
		}
	}
	if *revoke != "" {
		if err := openvpn.Revoke(*dir, *pkiDir, *revoke, *url); err != nil {
			log.Fatal(err)
		}
	}
}
