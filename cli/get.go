package cli

import (
	"bufio"
	"errors"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/assigner/etcdclient"
	"github.com/yansmallb/assigner/ping"
)

func get(containerid, app, component string, etcdpath string) error {
	log.Infoln("cli.get():Start get")
	log.Debugf("cli.get(): containerid=%+v", containerid)

	// connect to etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Errorf("cli.get():%+v\n", err)
		return err
	}
	//get config from etcd
	config, err := client.GetAbsoluteKey(etcdclient.ConfigPath)
	// check the config
	configs := strings.Split(config, ";")
	if len(configs) != 2 {
		err := errors.New("please start manage and set config first !!!")
		log.Errorf("cli.get():%+v\n", err)
		return err
	}
	bridge := configs[0]
	netmask := configs[1]

	// get all ips
	ips, err := client.GetAbsoluteDirIps(etcdclient.AppsPath + app + "/" + component + "/")
	if err != nil {
		log.Errorf("cli.get():%+v\n", err)
		return err
	}
	// check the len of ips , if len == 0 return err
	if len(ips) == 0 {
		log.Errorf("cli.get(): this app %+v:%+v dont have ip!!! and exit!!!", app, component)
		return errors.New("this app dont have ip!!!")
	}

	// rand a start index
	start := rand.Intn(100)
	var ip = ""
	var gateway = ""
	// try to get one ip
	for times := 0; times < 10; times++ {
		for i := 0; i < len(ips); i++ {
			index := (start + i) % len(ips)
			log.Debugf("cli.get(): try to lock ip=%+v \n", ips[index].Ip)
			con, err := client.CreateKey(ips[index].Ip, containerid)
			if err != nil {
				log.Errorf("cli.get():%+v\n", err)
				continue
			}
			if con != containerid {
				log.Infof("cli.get():this container %+v lock this ip %+v and continue\n", con, ips[index])
				continue
			}
			ip = ips[index].Ip
			gateway = ips[index].Gateway
			log.Infof("cli.get(): container %+v lock this ip %+v \n", containerid, ips[index])
			break
		}
		if ip != "" {
			break
		}
		log.Warnf("cli.get(): dont have empty ip , will start after 15 second again\n")
		time.Sleep(15 * time.Second)
	}

	// has get ip and try to set ip
	hasSet := false
	strip := ip + "/" + netmask + "@" + gateway
	for times := 3; times != 0; times-- {
		ExecPipework(bridge, containerid, strip)
		time.Sleep(5 * time.Second)
		// ping this ip and check it is useful
		alive := ping.Ping(ip, 3)
		log.Debugf("cli.get():ping  ip=%+v,alive=%+v", ip, alive)
		if alive != true {
			continue
		}

		hasSet = true
		log.Infof("cli.get():get ip success, container=%+v , ip=%+v\n", containerid, ip)
		break
	}
	// set ip failed
	if !hasSet {
		client.DeleteKey(ip)
		log.Errorf("cli.get(): set container %+v three times with ip %+v failed, unlock this ip and exit!!!", containerid, ip)
		return errors.New("set ip failed!!!")
	}
	// set containerid/ip  in etcd
	client.CreateAbsoluteKey(etcdclient.IdsPath+containerid, ip)
	log.Infoln("cli.get(): Get success")
	return nil
}

func ExecPipework(bridge, containerid, strip string) {
	cmd := exec.Command("pipework", bridge, containerid, strip)
	log.Debugf("cli.get(): args for pipework  %+v %+v %+v \n", bridge, containerid, strip)
	// get  pipe stdout
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("Error creating StdoutPipe for Cmd :%+v\n", err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		// watch the stdout for cmd
		for scanner.Scan() {
			log.Warnf("exec pipework out:%+v\n", scanner.Text())
		}
	}()
	// run cmd
	if err := cmd.Run(); err != nil {
		log.Errorf("cli.get(): cmd exec %+v\n", err)
	}
}
