package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/assigner/etcdclient"
)

func manage(bridge, netmask string, etcdpath string) error {
	log.Infoln("cli.manage():Start manage")

	// connect to etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.manage():%+v\n", err)
		fmt.Printf("[error]cli.manage():%+v\n", err)
		return err
	}
	//create dir on etcd
	client.CreateAbsoluteDir(etcdclient.AppsPath)
	client.CreateAbsoluteDir(etcdclient.IpsPath)
	client.CreateAbsoluteDir(etcdclient.IdsPath)

	//create config on etcd
	config := bridge + ";" + netmask
	log.Debugln("cli.manage() config{bridge;netmask} :%+v", config)
	_, err = client.CreateAbsoluteKey(etcdclient.ConfigPath, config)
	if err != nil {
		log.Errorf("cli.createConfig():%+v\n", err)
	}
	log.Infoln("cli.manage(): Manage success")
	return nil
}
