package pkgs

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"sort"
	"strings"
)

func WatchNameFile(closeCh chan struct{}, nameCh chan []string) {
	fw := NewFileWatcher(5)
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
	fmt.Println("darwin.nameserver:", nameserver)

	return nameserver, nil

}
