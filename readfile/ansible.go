package readfile

import (
	"regexp"
)

var (
	regexpUser, regexpPass, regexpPort, RegexpGroupVars, RegexpGroupChildren, RegexpGroup, RegexpHost *regexp.Regexp
)

type AnsibleVars struct {
	SSHUser string
	SSHPass string
	SSHPort string
}

func init() {
	if regexpTmp, err := regexp.Compile(`ansible_ssh_user=([^\s]+)`); err != nil {
		panic(err)
	} else {
		regexpUser = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`ansible_ssh_pass=([^\s]+)`); err != nil {
		panic(err)
	} else {
		regexpPass = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`ansible_ssh_port=([0-9]+)`); err != nil {
		panic(err)
	} else {
		regexpPort = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`^\[([^\s]+):vars\]$`); err != nil {
		panic(err)
	} else {
		RegexpGroupVars = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`^\[([^\s]+):children\]$`); err != nil {
		panic(err)
	} else {
		RegexpGroupChildren = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`^\[([^:]+)\]$`); err != nil {
		panic(err)
	} else {
		RegexpGroup = regexpTmp
	}

	if regexpTmp, err := regexp.Compile(`^([^\s]+)`); err != nil {
		panic(err)
	} else {
		RegexpHost = regexpTmp
	}
}

func (vars *AnsibleVars) Update(line Line) bool {
	updated := false
	matchSub := regexpUser.FindStringSubmatch(line.Line)
	if len(matchSub) == 2 {
		vars.SSHUser = matchSub[1]
		updated = true
	}

	matchSub = regexpPass.FindStringSubmatch(line.Line)
	if len(matchSub) == 2 {
		vars.SSHPass = matchSub[1]
		updated = true
	}

	matchSub = regexpPort.FindStringSubmatch(line.Line)
	if len(matchSub) == 2 {
		vars.SSHPort = matchSub[1]
		updated = true
	}

	return updated
}

type AnsibleHost struct {
	Host string
	Vars AnsibleVars
}

type AnsibleHosts map[string]*AnsibleHost

func (hosts AnsibleHosts) AddOrUpdateHost(host *AnsibleHost) {
	hosts[host.Host] = host
}

type AnsibleGroup struct {
	GroupName string
	Vars      AnsibleVars
	Children  AnsibleGroups
	Hosts     AnsibleHosts
}

func NewAnsibleGroup(groupName string) *AnsibleGroup {
	return &AnsibleGroup{
		GroupName: groupName,
		Vars:      AnsibleVars{},
		Children:  AnsibleGroups{},
		Hosts:     AnsibleHosts{},
	}
}

type AnsibleGroups map[string]*AnsibleGroup

// 没有就添加并返回，有就返回
func (groups AnsibleGroups) AddOrGetGroup(groupName string) *AnsibleGroup {
	if _, ok := groups[groupName]; !ok {
		groups[groupName] = NewAnsibleGroup(groupName)
	}
	return groups[groupName]
}

func (groups AnsibleGroups) AddOrUpdateGroup(group *AnsibleGroup) {
	groups[group.GroupName] = group
}

func ParseAnsibleFile(filepath string) (AnsibleGroups, error) {
	lines, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	ansibleGroups := AnsibleGroups{}
	var group *AnsibleGroup
	var nextLineType LineType
	for _, line := range lines {
		matchSub := RegexpGroupVars.FindStringSubmatch(line.Line)
		if len(matchSub) == 2 {
			group = ansibleGroups.AddOrGetGroup(matchSub[1])
			nextLineType = LineTypeAnsibleGroupVars
			continue
		}

		matchSub = RegexpGroupChildren.FindStringSubmatch(line.Line)
		if len(matchSub) == 2 {
			group = ansibleGroups.AddOrGetGroup(matchSub[1])
			nextLineType = LineTypeAnsibleGroupChildren
			continue
		}

		matchSub = RegexpGroup.FindStringSubmatch(line.Line)
		if len(matchSub) == 2 {
			group = ansibleGroups.AddOrGetGroup(matchSub[1])
			nextLineType = LineTypeAnsibleHost
			continue
		}

		if nextLineType == LineTypeAnsibleGroupVars {
			group.Vars.Update(line)
			ansibleGroups.AddOrUpdateGroup(group)
			continue
		}

		if nextLineType == LineTypeAnsibleGroupChildren {
			group.Children.AddOrGetGroup(line.Line)
			ansibleGroups.AddOrUpdateGroup(group)
			continue
		}

		if nextLineType == LineTypeAnsibleHost {
			matchSub = RegexpHost.FindStringSubmatch(line.Line)
			if len(matchSub) == 2 {
				host := &AnsibleHost{
					Host: matchSub[1],
					Vars: AnsibleVars{},
				}
				host.Vars.Update(line)
				group.Hosts.AddOrUpdateHost(host)
				ansibleGroups.AddOrUpdateGroup(group)
			}
		}
	}

	// 数据变量
	for _, group := range ansibleGroups {
		for _, host := range group.Hosts {
			if host.Vars.SSHUser == "" {
				host.Vars.SSHUser = group.Vars.SSHUser
			}
			if host.Vars.SSHPass == "" {
				host.Vars.SSHPass = group.Vars.SSHPass
			}
			if host.Vars.SSHPort == "" {
				host.Vars.SSHPort = group.Vars.SSHPort
			}
			group.Hosts.AddOrUpdateHost(host)
		}
		ansibleGroups.AddOrUpdateGroup(group)
	}

	for _, group := range ansibleGroups {
		for _, child := range group.Children {
			g := ansibleGroups.AddOrGetGroup(child.GroupName)
			for _, host := range g.Hosts {
				if host.Vars.SSHUser == "" {
					host.Vars.SSHUser = group.Vars.SSHUser
				}
				if host.Vars.SSHPass == "" {
					host.Vars.SSHPass = group.Vars.SSHPass
				}
				if host.Vars.SSHPort == "" {
					host.Vars.SSHPort = group.Vars.SSHPort
				}
				group.Hosts.AddOrUpdateHost(host)
			}
		}
		ansibleGroups.AddOrUpdateGroup(group)
	}

	return ansibleGroups, nil
}
