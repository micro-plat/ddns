package pkgs

import (
	"fmt"
	"sort"
	"strings"
	"time" 

	"github.com/golang/sys/windows/registry"
)

const (
	registrykey = `SYSTEM\CurrentControlSet\services\Tcpip\Parameters\Interfaces`
)

func WatchNameFile(closeCh chan struct{}, nameCh chan []string) {
	last := GetNameServers()
	period := time.Second * 5
	ticker := time.NewTicker(period)
	for {
		select {
		case <-closeCh:
			return
		case <-ticker.C:
			ticker.Stop()
			news := GetNameServers()

			if len(news) != len(last) {
				last = news
				nameCh <- news
			}
			isMatch := true
			for i := range news {
				if !strings.EqualFold(news[i], last[i]) {
					isMatch = false
					break
				}
			}
			if !isMatch {
				last = news
				nameCh <- news
			}
			ticker.Reset(period)
		}

	}
}

func GetNameServers() (nameserver []string, err error) {
	rootkey, err := registry.OpenKey(registry.LOCAL_MACHINE, registrykey, registry.QUERY_VALUE)

	if err != nil {
		err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s失败：%w`, registrykey, err)
		return
	}
	defer rootkey.Close()

	subKeys, err := rootkey.ReadSubKeyNames(rootkey.SubKeyCount)
	if err != nil {
		err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s 子节点失败：%w`, registrykey, err)
		return
	}

	for i, sk := range subKeys {
		regsubkey, err := registry.OpenKey(registry.LOCAL_MACHINE, registrykey+`\`+sk, registry.QUERY_VALUE)
		if err != nil {
			err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s 子节点失败：%w`, registrykey, err)
			return
		}
		val, _, err := regsubkey.GetStringValue("NameServer")
		if err != nil {
			continue
		}
		nameserver = append(nameserver, val...)
	}
	nameserver = Distinct(nameserver)
	sort.Strings(nameserver)
	return
}
