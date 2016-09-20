package cli

import (
	"errors"
	"time"

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
	// delete from etcd
	err = deleteIp(client, ip)
	if err != nil {
		log.Errorf("cli.delete():%+v\n", err)
		return err
	}
	log.Infoln("cli.delete():Delete success")
	return nil
}

func deleteByContainerid(containerid string, etcdPath string) error {
	log.Infoln("cli.deleteByContainerid():Start DeleteByContainerid")
	client, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		log.Errorf("cli.deleteByContainerid():%+v\n", err)
		return err
	}

	//query from etcd
	ip, err := client.QueryContainerid(containerid)
	if err != nil {
		log.Errorf("cli.deleteByContainerid():%+v\n", err)
		return err
	}
	log.Infof("cli.deleteByContainerid(): result ip=%+v", ip)
	//delete from etcd
	err = deleteIp(client, ip)
	if err != nil {
		log.Errorf("cli.deleteByContainerid():%+v\n", err)
		return err
	}
	log.Infoln("cli.deleteByContainerid():Delete success")
	return nil
}

func deleteIp(client *etcdclient.Etcd, ip string) error {
	//delete from etcd
	containerid := ""
	for i := 1; i <= 3; i++ {
		tempid, err := client.DeleteKey(ip)
		if err == nil {
			// check delete
			containerid = tempid
			_, err := client.GetKey(ip)
			if err != nil {
				// deleted  and break
				log.Infof("cli.deleteIp():try to delete id %+v \n", etcdclient.IdsPath+containerid)
				client.DeleteAbsoluteKey(etcdclient.IdsPath + containerid)
				return nil
			}
		}
		log.Errorf("cli.deleteIp():%+v, had try %+v times \n", err, i)
		time.Sleep(15 * time.Second)
	}
	return errors.New("delete three times and failed!!!")
}
