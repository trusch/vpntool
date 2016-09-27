package openvpn

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/trusch/vpntool/pki"
)

// Init initializes an new VPN
func Init(pkiDir, outputDirectory string, peerToPeer bool) error {
	err := pki.Init(pkiDir)
	if err != nil {
		return err
	}
	return createServer(pkiDir, outputDirectory, peerToPeer)
}

// CreateClient generates new keys and a openvpn configfile
func CreateClient(pkiDir, id, url, outputDirectory string) error {
	err := pki.AddClient(pkiDir, id)
	if err != nil {
		return err
	}
	return createClientConfig(pkiDir, id, url, outputDirectory)
}

// Deploy deploys a vpn config
func Deploy(configDir, pkiDir, id, target string) error {
	copyScript := fmt.Sprintf(`scp %v/%v.conf %v/pki/crl.pem %v:/tmp/`, configDir, id, pkiDir, target)
	installScript := fmt.Sprintf(`
    if ! dpkg-query -W openvpn; then
      sudo apt-get -y install openvpn
    fi
    sudo mv /tmp/%v.conf /etc/openvpn/
		sudo mv /tmp/crl.pem /etc/openvpn/
    sudo systemctl enable openvpn@%v
    sudo systemctl start openvpn@%v
  `, id, id, id)
	copyCmd := exec.Command("bash", "-c", copyScript)
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr
	copyCmd.Stdin = os.Stdin
	if err := execScript(copyScript); err != nil {
		return err
	}
	return execScriptRemote(installScript, target)
}

// Revoke revokes a VPN clients certificate, and deploys the CRL to a server if last argument is non-empty
func Revoke(configDir, pkiDir, id, serverTarget string) error {
	if err := pki.Revoke(pkiDir, id); err != nil {
		return err
	}
	deleteScript := fmt.Sprintf(`rm	%v/%v.conf`, configDir, id)
	if err := execScript(deleteScript); err != nil {
		return err
	}
	if serverTarget != "" {
		copyScript := fmt.Sprintf(`scp %v/pki/crl.pem %v:/tmp/crl.pem`, pkiDir, serverTarget)
		if err := execScript(copyScript); err != nil {
			return err
		}
		installScript := fmt.Sprintf(`sudo mv /tmp/crl.pem /etc/openvpn/crl.pem`)
		if err := execScriptRemote(installScript, serverTarget); err != nil {
			return err
		}
	}
	return nil
}

func execScript(script string) error {
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func execScriptRemote(script, target string) error {
	cmd := exec.Command("ssh", "-t", target, script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func createServer(pkiDir, outputDirectory string, peerToPeer bool) error {
	err := pki.AddServer(pkiDir, "vpnserver")
	if err != nil {
		return err
	}
	err = pki.CreateDH(pkiDir)
	if err != nil {
		return err
	}
	return createServerConfig(pkiDir, outputDirectory, peerToPeer)
}

// createServerConfig generates a server config
func createServerConfig(pkiDir, outputDirectory string, peerToPeer bool) error {
	ca, err := getCA(pkiDir)
	if err != nil {
		return err
	}
	key, err := getKey(pkiDir, "vpnserver")
	if err != nil {
		return err
	}
	cert, err := getCert(pkiDir, "vpnserver")
	if err != nil {
		return err
	}
	dh, err := getDH(pkiDir)
	if err != nil {
		return err
	}
	opts := templateOptions{
		CA:   ca,
		Key:  key,
		Cert: cert,
		DH:   dh,
	}
	if peerToPeer {
		opts.Additional = []string{"client-to-client"}
	}
	f, err := os.Create(outputDirectory + "/server.conf")
	if err != nil {
		return err
	}
	defer f.Close()
	return serverConfigTemplate.Execute(f, opts)
}

// createClientConfig generates a client config
func createClientConfig(pkiDir, id, url, outputDirectory string) error {
	ca, err := getCA(pkiDir)
	if err != nil {
		return err
	}
	key, err := getKey(pkiDir, id)
	if err != nil {
		return err
	}
	cert, err := getCert(pkiDir, id)
	if err != nil {
		return err
	}
	opts := templateOptions{
		CA:   ca,
		Key:  key,
		Cert: cert,
		URL:  url,
	}
	f, err := os.Create(fmt.Sprintf("%v/%v.conf", outputDirectory, id))
	if err != nil {
		return err
	}
	defer f.Close()
	return clientConfigTemplate.Execute(f, opts)
}

var serverConfigTemplateString = `
port 1194
proto tcp
dev tun
server 10.8.0.0 255.255.255.0
ifconfig-pool-persist ipp.txt
keepalive 10 120
comp-lzo
persist-key
persist-tun
status openvpn-status.log
crl-verify /etc/openvpn/crl.pem
verb 3
{{range .Additional}}
{{.}}
{{end}}
<ca>
{{.CA}}
</ca>
<cert>
{{.Cert}}
</cert>
<key>
{{.Key}}
</key>
<dh>
{{.DH}}
</dh>
`
var serverConfigTemplate = template.Must(template.New("serverConfig").Parse(serverConfigTemplateString))

var clientConfigTemplateString = `
client
dev tun
proto tcp
remote {{.URL}} 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
comp-lzo
verb 3
<ca>
{{.CA}}
</ca>
<cert>
{{.Cert}}
</cert>
<key>
{{.Key}}
</key>
`
var clientConfigTemplate = template.Must(template.New("clientConfig").Parse(clientConfigTemplateString))

type templateOptions struct {
	CA         string
	Cert       string
	Key        string
	DH         string
	URL        string
	Additional []string
}

func getFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func getCA(pkiDir string) (string, error) {
	return getFile(pkiDir + "/pki/ca.crt")
}
func getDH(pkiDir string) (string, error) {
	return getFile(pkiDir + "/pki/dh.pem")
}

func getKey(pkiDir, id string) (string, error) {
	return getFile(fmt.Sprintf("%v/pki/private/%v.key", pkiDir, id))
}

func getCert(pkiDir, id string) (string, error) {
	return getFile(fmt.Sprintf("%v/pki/issued/%v.crt", pkiDir, id))
}
