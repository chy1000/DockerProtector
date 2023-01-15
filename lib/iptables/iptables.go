package iptables

import (
	"DockerProtector/lib/utils"
	"strconv"
	"strings"
)

type Iptables map[string][]Item

type Item struct{
	Num int
	Target string
	Prot string
	Opt string
	Source string
	Destination string
	Port int
}

func Create() (Iptables, error){
	iptables := Iptables{}
	result, err := utils.ExecCmd("iptables", "--line-numbers", "-n", "-L", "DOCKER")
	if err != nil { return nil, err }
	if !strings.Contains(result, "Chain DOCKER") { return nil, nil }
	lines := strings.Split(result, "\n")
	for i, v := range lines {
		if i < 2 { continue }
		slice := strings.Fields(v)
		num, _ := strconv.Atoi(slice[0])
		dest := slice[5]
		portSlice := strings.Split(slice[7], ":")
		port, _ := strconv.Atoi(portSlice[1])
		item := Item{num, slice[1], slice[2], slice[3], slice[4], dest, port}
		iptables[dest] = append(iptables[dest], item)
	}
	return iptables, nil
}

func (i Iptables) Accept(destination, source string, port int) error {
	dockerIptables := i[destination]
	// 检测切片看是否有重复的记录
	for _, item := range dockerIptables {
		if source == item.Source && port == item.Port {
			return nil
		}
	}
	_, err := utils.ExecCmd(
		"iptables",
		"-A",
		"DOCKER",
		"-s", source,
		"-p", "tcp",
		"-d", destination,
		"--dport", strconv.Itoa(port),
		"-j", "ACCEPT",
	)
	if err != nil { return err }
	num := 1
	if len(dockerIptables) > 1 {
		num = dockerIptables[len(dockerIptables)-1].Num + 1
	}
	item := Item{num, "ACCEPT", "tcp", "--", source, destination, port}
	i[destination] = append(i[destination], item)
	return nil
}

func (i Iptables) Delete(destination, source string, port int) error {
	items := i[destination]
	length := len(items) - 1
	var tempItems []Item
	for j := length; j >= 0 ; j-- {
		item := items[j]
		if source == item.Source && port == item.Port {
			_, err := utils.ExecCmd("iptables", "-D", "DOCKER", strconv.Itoa(item.Num))
			if err != nil { return err }
			continue
		}
		tempItems = append(tempItems, item)
	}
	for j, item := range tempItems {
		item.Num = j + 1
	}
	i[destination] = tempItems
	return nil
}
