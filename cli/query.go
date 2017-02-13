package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/assigner/etcdclient"
)

func query(containerid string, etcdPath string) error {
	log.Infoln("cli.query():Start query")
	client, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		log.Errorf("cli.query():%+v\n", err)
		return err
	}

	//query from etcd
	ip, err := client.QueryContainerid(containerid)
	if err != nil {
		log.Errorf("cli.query():%+v\n", err)
	} else {
		log.Infof("cli.query(): result ip=%+v", ip)
	}
	log.Infoln("cli.query():Query success")
	return nil
}
