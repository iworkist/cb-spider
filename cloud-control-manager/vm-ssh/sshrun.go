// Package for VM's SSH and SCP of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2019.10.

package sshrun

import (
	"github.com/cloud-barista/cb-store/config"
        "github.com/sirupsen/logrus"


	"io"
	"time"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

var cblog *logrus.Logger

func init() {
        cblog = config.Cblogger
}

//====================================================================
type SSHInfo struct {
        UserName        string  // ex) "root"
        PrivateKey      []byte  // ex)   []byte(`-----BEGIN RSA PRIVATE KEY-----
                                //              MIIEoQIBAAKCAQEArVNOLwMIp5VmZ4VPZotcoCHdEzimKalAsz+ccLfvAA1Y2ELH
                                //              ...`)
        ServerPort      string  // ex) "node12:22"
	Timeout		int	// ex) 3 (= 3sec)
}
//====================================================================

func Connect(sshInfo SSHInfo) (scp.Client, error) {
	cblog.Info("call Connect()")

        clientConfig, _ := getClientConfig(sshInfo.UserName, sshInfo.PrivateKey, ssh.InsecureIgnoreHostKey(), sshInfo.Timeout)
        client := scp.NewClient(sshInfo.ServerPort, &clientConfig)
        err := client.Connect()
        return client, err
}

//====================================================================
type SSHKeyPathInfo struct {
        UserName        string  // ex) "root"
        KeyPath     	string 	// ex) "/root/.ssh/id_rsa // You should use the full path.
        ServerPort      string  // ex) "node12:22"
}
//====================================================================

func ConnectKeyPath(sshKeyPathInfo SSHKeyPathInfo) (scp.Client, error) {
        cblog.Info("call ConnectKeyPath()")

	clientConfig, _ := auth.PrivateKey(sshKeyPathInfo.UserName, sshKeyPathInfo.KeyPath, ssh.InsecureIgnoreHostKey())
        client := scp.NewClient(sshKeyPathInfo.ServerPort, &clientConfig)
        err := client.Connect()
        return client, err
}

func getClientConfig(username string, privateKey []byte, keyCallBack ssh.HostKeyCallback, timeout int) (ssh.ClientConfig, error) {

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	clientConfig := ssh.ClientConfig{
                User: username,
                Auth: []ssh.AuthMethod{
                        ssh.PublicKeys(signer),
                },
                HostKeyCallback: keyCallBack,
        }
	if timeout > 0 {
		clientConfig.Timeout = time.Second * time.Duration(timeout)
	}

	return clientConfig, nil
}

func Close(client scp.Client){
	cblog.Info("call Close()")

	client.Close()	
}

func RunCommand(client scp.Client, cmd string) (string, error) {
	cblog.Info("call RunCommand()")

	session := client.Session
	sshOut, err := session.StdoutPipe()
	session.Stderr = os.Stderr

	err = session.Run(cmd)
	//err = session.Start(cmd)

	return stdoutToString(sshOut), err
}

func stdoutToString(sshOut io.Reader) string {
	buf := make([]byte, 1000)
	num, err := sshOut.Read(buf)
	outStr := ""
	if err == nil {
		outStr = string(buf[:num])
	}
	for err == nil {
		num, err = sshOut.Read(buf)
		outStr += string(buf[:num])
		if err != nil {
			if err.Error() != "EOF" {
				cblog.Error(err)
			}
		}

	}
	return strings.Trim(outStr, "\n")
}

func Copy(client scp.Client, sourcePath string, remotePath string) error {
	cblog.Info("call Copy()")

        // Open a file
        file, _ := os.Open(sourcePath)
        defer file.Close()
        return client.CopyFile(file, remotePath, "0755")
}

//=============== for One Call Service
func SSHRun(sshInfo SSHInfo, cmd string) (string, error) {
        cblog.Info("call SSHRun()")


        sshCli, err := Connect(sshInfo)
        if err != nil {
                return "", err
        }
	defer Close(sshCli)

	return RunCommand(sshCli, cmd)
}

func SSHRunByKeyPath(sshInfo SSHKeyPathInfo, cmd string) (string, error) {
        cblog.Info("call SSHRunKeyPath()")


        sshCli, err := ConnectKeyPath(sshInfo)
        if err != nil {
                return "", err
        }
	defer Close(sshCli)

	return RunCommand(sshCli, cmd)
}

func SSHCopy(sshInfo SSHInfo, sourcePath string, remotePath string) error {
        cblog.Info("call SSHCopy()")


        sshCli, err := Connect(sshInfo)
        if err != nil {
                return err
        }
	defer Close(sshCli)

	return Copy(sshCli, sourcePath, remotePath)
}

func SSHCopyByKeyPath(sshInfo SSHKeyPathInfo, sourcePath string, remotePath string) error {
        cblog.Info("call SSHCopyByKeyPath()")


        sshCli, err := ConnectKeyPath(sshInfo)
        if err != nil {
                return err
        }
	defer Close(sshCli)

        return Copy(sshCli, sourcePath, remotePath)
}

