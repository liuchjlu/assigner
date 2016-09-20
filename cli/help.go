package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: assigner COMMAND [args...]
  Version: 0.3.1
  Author:yansmallb
  Email:yanxiaoben@iie.ac.cn

  Commands:
      import  [localyaml path] [etcd path]                 import apps' ips from yaml 
      delete  [ip]/[containerid] [etcd path]               try to recycling one ip and delete the key from etcd
      get     [container id] [app:component] [etcdpath]    try to get one ip for container and will exit after get one ip
      manage  [bridge] [netmask] [etcd path]               create the config in etcd 
      query   [containerid] [etcd path]                    query ip by containerid 
      help
  `
	fmt.Println(helpString)
}
