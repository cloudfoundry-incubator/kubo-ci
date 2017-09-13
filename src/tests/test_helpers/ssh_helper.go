package test_helpers

import (
	"log"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/cloudfoundry/bosh-utils/errors"
)

func RunSSHCommand(server string, port int, username string, privateKey string, command string) (string, error) {
	parsedPrivateKey, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		log.Println(err)
		return "", err
	}

	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(parsedPrivateKey),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, port), config)
	if err != nil {
		return "", errors.WrapError(err, "Cannot dial")
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", errors.WrapError(err, "Cannot create session")
	}
	defer session.Close()

	var output bytes.Buffer

	session.Stdout = &output

	session.Run(command)

	return output.String(), nil
}

