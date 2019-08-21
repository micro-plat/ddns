package services

import "net"
import "strings"

type ICache interface {
	Lookup(string) []net.IP
	Save(string) error
}
type Cache struct {
}

func NewCache() *Cache {
	return &Cache{}
}
func (c *Cache) Lookup(name string) []net.IP {
	if strings.HasPrefix(name, "github.com") {
		return []net.IP{
			net.ParseIP("140.82.114.3").To4(),
		}
	}

	return nil

}
func (c *Cache) Save(string) error {
	return nil
}
