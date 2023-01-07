package iptables

import (
	"DockerProtector/lib/utils"
	"strconv"
	"strings"
)

type Iptables map[string][]*Item

type Item struct{
	Num int
	Target string
	Prot string
	Opt string
	Source string
	Destination string
	Port string
}

func Create() (*Iptables, error){
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
		item := &Item{num, slice[1], slice[2], slice[3], slice[4], dest, portSlice[1]}
		iptables[dest] = append(iptables[dest], item)
	}
	return &iptables, nil
}

func (i *Iptables) Accept() {

}

func (i *Iptables) Drop() {

}

func (i *Iptables) Delete() {

}
