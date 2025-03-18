package readfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// ansible inventory文件的变量
type Vars struct {
	SSHPort string
	SSHUser string
	SSHPass string
}

type Host struct {
	Host string
	Vars
}

type Group struct {
	Name  string
	Hosts []Host
	Vars
}

// 读取ansible inventory文件
func ReadAnsbileInventoryFile(hosts string) ([]Group, error) {
	f, err := os.Open(hosts)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	blankLine, err := regexp.Compile(`^\s*$`)
	if err != nil {
		log.Printf("regexp.Compile error: %v\n", err)
		return nil, err
	}

	reader := bufio.NewReader(f)
	lines := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if blankLine.Match([]byte(line)) {
			continue
		}
		lines = append(lines, line)
	}

	return parseGroups(lines)
}

func parseGroups(lines []string) ([]Group, error) {
	groups := []Group{}

	groupRegexp, err := regexp.Compile(`^\s*\[(.+)\]\s*$`)
	if err != nil {
		log.Printf("regexp.Compile error: %v\n", err)
	}

	groupName := ""
	groupVars := false
	for _, line := range lines {
		subMatch := groupRegexp.FindStringSubmatch(line)

		// group的行  [www]  [www:vars]
		if len(subMatch) == 2 {
			tmpSplit := strings.Split(subMatch[1], ":")
			// [www]
			if len(tmpSplit) == 1 {
				groupName = tmpSplit[0]
				groupVars = false
				continue
			} else if len(tmpSplit) == 2 {
				// [www:vars]
				if tmpSplit[1] == "vars" {
					groupName = tmpSplit[0]
					groupVars = true
					continue
					// [www:unknown]
				} else {
					groupName = tmpSplit[0]
					log.Printf("groupName: %s, unknown tag: %s\n", groupName, tmpSplit[1])
					continue
				}
			} else {
				// [www:a:b]
				log.Printf("groupName: %s, line '%s' split with ':' get more than 2 items\n", groupName, line)
				continue
			}
		}

		// 不是group的行
		index := -1
		for i, groupItem := range groups {
			if groupItem.Name == groupName {
				index = i
				break
			}
		}
		if index == -1 {
			groups = append(groups, Group{
				Name:  groupName,
				Hosts: []Host{},
			})
			index = len(groups) - 1
		}

		if groupVars {
			parseGroupVars(line, &groups[index])
		} else {
			host, err := parseHost(line)
			if err != nil {
				log.Printf("parseHost error: %v\n", err)
				continue
			}

			inHosts := false
			for indexInHosts, hostItem := range groups[index].Hosts {
				if hostItem.Host == host.Host {
					groups[index].Hosts[indexInHosts] = host
					inHosts = true
					break
				}
			}
			if !inHosts {
				groups[index].Hosts = append(groups[index].Hosts, host)
			}
		}
	}

	// 从groups vars同步vars到hosts
	for indexInGroups, groupItem := range groups {
		for indexInHosts, hostItem := range groupItem.Hosts {
			if hostItem.SSHPort == "" {
				groups[indexInGroups].Hosts[indexInHosts].SSHPort = groupItem.SSHPort
			} else if hostItem.SSHUser == "" {
				groups[indexInGroups].Hosts[indexInHosts].SSHUser = groupItem.SSHUser
			} else if hostItem.SSHPass == "" {
				groups[indexInGroups].Hosts[indexInHosts].SSHPass = groupItem.SSHPass
			}
		}
	}
	return groups, nil
}

// 解析group vars
func parseGroupVars(line string, group *Group) error {
	tmpSplit := strings.Split(line, "=")
	if len(tmpSplit) == 2 {
		if tmpSplit[0] == "ansible_ssh_user" {
			group.SSHUser = tmpSplit[1]
		} else if tmpSplit[0] == "ansible_ssh_pass" {
			group.SSHPass = tmpSplit[1]
		} else if tmpSplit[0] == "ansible_ssh_port" {
			group.SSHPort = tmpSplit[1]
		} else {
			return fmt.Errorf("unknown group vars: %s", tmpSplit[0])
		}
		return nil
	}
	return fmt.Errorf("line '%s' split with '=' get not 2 items", line)
}

// 解析host
func parseHost(line string) (host Host, err error) {
	fields := strings.Fields(line)
	if len(fields) >= 1 {
		host.Host = fields[0]
	} else {
		return host, errors.New("line split with ' ' get less than 1 items")
	}

	for _, field := range fields[1:] {
		if strings.Contains(field, "=") {
			tmpSplit := strings.Split(field, "=")
			if len(tmpSplit) == 2 {
				if tmpSplit[0] == "ansible_ssh_user" {
					host.SSHUser = tmpSplit[1]
				} else if tmpSplit[0] == "ansible_ssh_pass" {
					host.SSHPass = tmpSplit[1]
				} else if tmpSplit[0] == "ansible_ssh_port" {
					host.SSHPort = tmpSplit[1]
				} else {
					// log.Printf("unknown host vars: %s\n", tmpSplit[0])
				}
			}
		} else {
			log.Printf("host field '%s' format error\n", field)
		}
	}

	return host, nil
}
