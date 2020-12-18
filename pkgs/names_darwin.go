package pkgs

import (
	"bufio"
	"net"
	"os"
	"time"

	"sort"
	"strings"
)

func WatchNameFile(closeCh chan struct{}, nameCh chan []string) {
	//检查周期
	tickersec := time.Minute / time.Second
	fw := NewFileWatcher(int(tickersec))
	fw.Change = func(string) error {
		names, err := GetNameServers()
		if err != nil {
			return err
		}
		nameCh <- names
		return nil
	}
	fw.Deleted = func(string) error {
		names, err := GetNameServers()
		if err != nil {
			return err
		}
		nameCh <- names
		return nil
	}
	fw.Add(NAME_FILE)
	select {
	case <-closeCh:
		return
	}
}

func GetNameServers() (nameserver []string, err error) {
	buf, err := os.Open(NAME_FILE)
	if err != nil {
		return []string{}, nil
	}
	defer buf.Close()

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := PrepareLine(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		ip := net.ParseIP(line)
		if ip == nil {
			continue //ip格式错误
		}
		nameserver = append(nameserver, line)

	}
	nameserver = Distinct(nameserver)
	sort.Strings(nameserver)
 
	return nameserver, nil

}
