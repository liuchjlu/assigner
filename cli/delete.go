package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/assigner/etcdclient"
)

func delete(ip string, etcdPath string) error {
	log.Infoln("cli.delete():Start Delete")
	client, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		log.Errorf("cli.delete():%+v\n", err)
		return err
	}

	//delete from etcd
	err = client.DeleteKey(ip)
	if err != nil {
		log.Errorf("cli.delete():%+v\n", err)
	}
	log.Infoln("cli.delete():Delete success")
	return nil
}
