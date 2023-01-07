package main

import (
	"DockerProtector/lib/docker"
	"fmt"
	"github.com/docker/docker/api/types/events"
)

func main() {
	docker.Events(func(eventsMsg events.Message){
		if eventsMsg.Type != "container" { return }
		if eventsMsg.Status != "start" { return }
		name := eventsMsg.Actor.Attributes["name"]
		fmt.Println(name)
		/*
		// 不在规则表的不处理
		if _, ok := rules.Rules[name]; !ok { return }

		rule := rules.Rules[name]
		result, _ := rule.GetIptables()
		fmt.Println(result)

		ip = ips[name]
		# 删除所有规则
		delRules(ip)
		# 根据规则表重新添加所有规则
		addRules(ip, rules[name])
		# 显示最新的规则
		print( showRules(ip) )
		 */
	})
}

