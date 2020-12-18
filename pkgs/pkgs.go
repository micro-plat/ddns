package pkgs

import (
	"strings"

	"github.com/micro-plat/hydra/global"
)

func GetSyncData(syncChan chan string) (files []string) {
	for {
		select {
		case p := <-syncChan:
			files = append(files, p)
		default:
			return
		}
	}
}

func Distinct(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	list := make(map[string]string)
	for _, a := range arr {
		if _, ok := list[a]; !ok {
			list[a] = a
			newArr = append(newArr, a)
		}
	}
	return newArr
}

//Filte 过滤掉本机IP，0.0.0.0,127.0.0.1
func Filte(hosts ...string) []string {
	nhosts := make([]string, 0, len(hosts))
	localIP := global.LocalIP()
	for _, host := range hosts {
		if strings.HasPrefix(host, localIP) ||
			strings.HasPrefix(host, "0.0.0.0") ||
			strings.HasPrefix(host, "127.0.0.1") {
			continue
		}
		nhosts = append(nhosts, host)
	}
	return nhosts
}
func PrepareLine(txt string) string {
	txt = strings.TrimSpace(txt)
	txt = strings.Replace(txt, "\t", " ", -1)

	for {
		nbf := len(txt)
		txt = strings.Replace(txt, "  ", " ", -1)
		if nbf == len(txt) {
			break
		}
	}
	return txt
}
