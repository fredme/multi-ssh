package sshclient

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func SSHCommand(privateKeyFile, ip, port, username, password, command string) (string, error) {
	allAuthMethods := []ssh.AuthMethod{}

	publicKeyAuthMethod, err := sshPublicKeyAuthMethod(privateKeyFile)
	if err != nil {
		log.Printf("SSHPublicKeyAuthMethod error: %v", err)
	}
	allAuthMethods = append(allAuthMethods, publicKeyAuthMethod)

	if password != "" {
		passwordAuthMethod, err := sshPasswordAuthMethod(password)
		if err != nil {
			log.Printf("SSHPasswordAuthMethod error: %v", err)
		}
		allAuthMethods = append(allAuthMethods, passwordAuthMethod)
	}

	if len(allAuthMethods) == 0 {
		log.Printf("no auth methods provided")
		return "", nil
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ip, port), &ssh.ClientConfig{
		User:            username,
		Auth:            allAuthMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
