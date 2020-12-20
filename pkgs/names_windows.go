package pkgs

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	registrykey = `SYSTEM\CurrentControlSet\services\Tcpip\Parameters\Interfaces`
)

func WatchNameFile(closeCh chan struct{}, nameCh chan []string) {
	last, err := GetNameServers()
	if err != nil {
		panic(err)
	}
	period := time.Minute
	ticker := time.NewTicker(period)
	for {
		select {
		case <-closeCh:
			return
		case <-ticker.C:
			ticker.Stop()
			news, _ := GetNameServers()

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
	rootkey, err := registry.OpenKey(registry.LOCAL_MACHINE, registrykey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s失败：%w`, registrykey, err)
		return
	}
	defer rootkey.Close()
	info, err := rootkey.Stat()
	subKeys, err := rootkey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s 子节点失败：%w`, registrykey, err)
		return
	}

	for _, sk := range subKeys {
		regsubkey, err1 := registry.OpenKey(registry.LOCAL_MACHINE, registrykey+`\`+sk, registry.QUERY_VALUE)
		if err1 != nil {
			err = fmt.Errorf(`读取HKEY_LOCAL_MACHINE\%s 子节点失败：%w`, registrykey, err1)
			return
		}

		val, _, err := regsubkey.GetStringValue("NameServer")
		if err != nil {
			continue
		}
		val = strings.TrimSpace(val)
		if val != "" {
			nameserver = append(nameserver, strings.Split(val, ",")...)
			fmt.Println("sk:", val, nameserver)
		}
	}
	nameserver = Distinct(nameserver)
	sort.Strings(nameserver)
	return
}
