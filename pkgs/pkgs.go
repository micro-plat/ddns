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
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
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
