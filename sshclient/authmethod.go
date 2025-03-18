package sshclient

import (
	"multi-ssh/readfile"

	"golang.org/x/crypto/ssh"
)

// AuthMethod 私钥 认证方法
func sshPublicKeyAuthMethod(privateKeyFile string) (ssh.AuthMethod, error) {
	key, err := readfile.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	sshKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(sshKey), nil
}

// AuthMethod 密码 认证方法
func sshPasswordAuthMethod(password string) (ssh.AuthMethod, error) {
	return ssh.Password(password), nil
}
