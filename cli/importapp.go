package cli

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/assigner/etcdclient"
)

func importapp(filepath string, etcdpath string) error {
	log.Infoln("cli.importapp():Start importapp")
	// yaml unmarshal
	app, err := UnmarshalConfig(filepath)
	if err != nil {
		log.Fatalf("cli.importapp():%+v\n", err)
		return err
	}
	log.Debugf("cli.importapp(): app=%+v\n", app)

	// connect to etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.importapp():%+v\n", err)
		return err
	}

	// create app dir to etcd
	path := etcdclient.AppsPath + app.App + "/" + app.Component + "/"
	err = client.CreateAbsoluteDir(path)
	if err != nil {
		log.Fatalf("cli.importapp():%+v\n", err)
		return err
	}
	// create ips to etcd
	for _, ip := range app.Ips {
		_, err = client.CreateAbsoluteKey(path+ip.Ip, ip.Gateway)
		if err != nil {
			log.Errorf("cli.importapp():%+v\n", err)
		}
	}
	log.Infoln("cli.importapp(): Importapp success")
	return nil
}

func UnmarshalConfig(path string) (*etcdclient.App, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	app := new(etcdclient.App)
	err = yaml.Unmarshal(in, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
