package pkgs

import (
	"strings"
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
