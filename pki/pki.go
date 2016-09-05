package pki

import (
	"fmt"
	"os"
	"os/exec"
)

// Init initializes a new pki
func Init(dir string) error {
	script := fmt.Sprintf(`
wget -O /tmp/easy-rsa.tgz https://github.com/OpenVPN/easy-rsa/releases/download/3.0.1/EasyRSA-3.0.1.tgz;
tar xfvz /tmp/easy-rsa.tgz;
mv EasyRSA-3.0.1 %v;
cd %v;
./easyrsa init-pki;
echo default | ./easyrsa build-ca nopass;
  `, dir, dir)
	return execScript(script)
}

// AddServer creates a new server cert and signs it with the specified ca
func AddServer(dir, serverID string) error {
	script := fmt.Sprintf(`cd %v && ./easyrsa build-server-full %v nopass`, dir, serverID)
	return execScript(script)
}

// AddClient creates a new client cert and signs it with the specified ca
func AddClient(dir, clientID string) error {
	script := fmt.Sprintf(`cd %v && ./easyrsa build-client-full %v nopass`, dir, clientID)
	return execScript(script)
}

// CreateDH creates diffi hellman parameters
func CreateDH(dir string) error {
	script := fmt.Sprintf(`cd %v && ./easyrsa gen-dh`, dir)
	return execScript(script)
}

func execScript(script string) error {
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
