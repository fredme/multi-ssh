package main

import (
	"errors"
	"flag"
	"log"
	"multi-ssh/readfile"
	"multi-ssh/sshclient"
	"os"
	"os/user"
	"path"
	"sync"
)

var (
	privateKeyFile *string
	hostFile       *string
	groupName      *string
	command        *string
	debug          *bool
	tryRun         *bool
)

func init() {
	defaultPrivateKeyFile, err := userPrimaryKeyPath()
	if err != nil {
		log.Printf("get default userPrimaryKeyPath error: %\n", err)
	}
	privateKeyFile = flag.String("i", defaultPrivateKeyFile, "PrimaryKey File")

	defaultHostFile, err := hostsPath()
	if err != nil {
		log.Printf("get default hostsPath error: %\n", err)
	}
	hostFile = flag.String("H", defaultHostFile, "Ansible Inventory File")

	groupName = flag.String("g", "", "Hosts In Which Group")
	command = flag.String("cmd", "", "Command To Execute")
	tmp := false
	debug = &tmp
	tryRun = &tmp
	debug = flag.Bool("debug", false, "Show Debug Info")
	tryRun = flag.Bool("try", false, "Try Run, Just Show Command")
	flag.Parse()

	if *command == "" || *groupName == "" || *hostFile == "" {
		if *debug {
			log.Printf("OS Args: %v\n", os.Args)
			log.Printf("PrimaryKey File: '%s'\n", *privateKeyFile)
			log.Printf("Ansible Inventory File: '%s'\n", *hostFile)
			log.Printf("Hosts In Which Group: '%s'\n", *groupName)
			log.Printf("Command To Execute: '%s'\n", *command)
		}
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	groups, err := readfile.ReadAnsbileInventoryFile(*hostFile)
	if err != nil {
		log.Printf("readfile.ReadAnsbileInventoryFile error: %v", err)
		return
	}

	group := readfile.Group{}
	for _, groupItem := range groups {
		if groupItem.Name == *groupName {
			group = groupItem
			break
		}
	}
	if group.Name == "" {
		log.Printf("group '%s' not found in inventory file (hosts)", *groupName)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(group.Hosts))
	for indexInHosts, host := range group.Hosts {
		if *tryRun {
			log.Printf("[TRY]\t%s\t[%d/%d]\t%s\n", *groupName, indexInHosts, len(group.Hosts), host.Host)
			wg.Done()
			continue
		}
		go func(indexInHosts int, host readfile.Host) {
			defer wg.Done()
			stdout, err := sshclient.SSHCommand(*privateKeyFile, host.Host, host.SSHPort, host.SSHUser, host.SSHPass, *command)
			if err != nil {
				log.Printf("[FAILED]\t%s\t[%d/%d]\t%s\n%s\n%s\n", *groupName, indexInHosts, len(group.Hosts), host.Host, err, stdout)
			} else {
				log.Printf("[SUCCESS]\t%s\t[%d/%d]\t%s\n%s\n", *groupName, indexInHosts, len(group.Hosts), host.Host, stdout)
			}
		}(indexInHosts, host)

	}

	wg.Wait()
}

// 获取当前用户的默认primarykey路径
func userPrimaryKeyPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(u.HomeDir, ".ssh", "id_rsa"), nil
}

// 获取默认hosts文件，按照当前目录、用户家目录、/etc/ansible/hosts的顺序查找
func hostsPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	paths := []string{
		"hosts",
		path.Join(u.HomeDir, "ansible", "hosts"),
		"/etc/ansible/hosts",
	}

	hostFile := ""
	for _, p := range paths {
		info, err := os.Stat(p)
		if os.IsNotExist(err) {
			continue
		} else if info.IsDir() {
			continue
		}
		hostFile = p
		break
	}

	if hostFile == "" {
		return "", errors.New("hosts file not found")
	}
	return hostFile, nil
}
