package cli

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Run() {
	if len(os.Args) == 1 {
		help()
		return
	}

	var err error
	command := os.Args[1]
	log.Debugf("cli.Run(): cli args:%+v\n", os.Args)
	if command == "import" {
		if len(os.Args) != 4 {
			importErr := "the `import` command takes two arguments. See help"
			log.Errorln(importErr)
			return
		}
		filePath := os.Args[2]
		etcdPath := os.Args[3]
		err = importapp(filePath, etcdPath)
	}
	if command == "delete" {
		if len(os.Args) != 4 {
			deleteErr := "the `delete` command takes two arguments. See help"
			log.Errorln(deleteErr)
			return
		}
		ip := os.Args[2]
		etcdPath := os.Args[3]
		err = delete(ip, etcdPath)
	}
	if command == "get" {
		if len(os.Args) != 5 {
			getErr := "the `get` command takes four arguments. See help"
			log.Errorln(getErr)
			return
		}
		containerId := os.Args[2]
		app := strings.Split(os.Args[3], ":") // include app, component
		if len(app) != 2 {
			getErr := "the `get` command App fomat not match. See help"
			log.Errorln(getErr)
			return
		}
		etcdPath := os.Args[4]
		err = get(containerId, app[0], app[1], etcdPath)
	}
	if command == "manage" {
		if len(os.Args) != 5 {
			manageErr := "the `manage` command takes five arguments. See help"
			log.Errorln(manageErr)
			return
		}
		bridge := os.Args[2]
		netmask := os.Args[3]
		etcdPath := os.Args[4]
		err = manage(bridge, netmask, etcdPath)
	}
	if command == "help" {
		help()
	}
	if err != nil {
		log.Errorf("cli.Run():%+v\n", err)
		return
	}
}
