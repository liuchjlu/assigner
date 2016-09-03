package cli

import (
	"testing"
)

func TestManage(t *testing.T) {
	bridge := "br0"
	netmask := "24"
	etcdpath := "http://192.168.11.51:2379"
	manage(bridge, netmask, etcdpath)
}
