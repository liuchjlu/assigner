package cli

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	containerid := ""
	app := "xiaoben"
	component := "test"
	etcdpath := "http://192.168.11.51:2379"
	err := get(containerid, app, component, etcdpath)
	fmt.Println(err)
}
