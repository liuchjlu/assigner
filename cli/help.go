package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: assigner COMMAND [args...]
  Version: 0.1.0
  Author:yansmallb
  Email:yanxiaoben@iie.ac.cn

  Commands:
      import  [localyaml path] [etcd path]                 import apps' ips from yaml 
      delete  [ip] [etcd path]                             try to recycling one ip and delete the key from etcd
      get     [container id] [app:component] [etcdpath]    try to get one ip for container and will exit after get one ip
      manage  [bridge] [netmask] [etcd path]               create the config in etcd 
      help
  `
	fmt.Println(helpString)
}
