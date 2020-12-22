package local

import (
	"fmt"
	"net"
	"strings"
)

//TrimDomain trim前面的www及前后的".
func TrimDomain(d string) string {
	return strings.Trim(strings.TrimPrefix(d, "www"), ".")
}
func HasWWW(d string) bool {
	return strings.HasPrefix(d, "www.")
}
func GetURL(proto, host, port string) string {
	if port == "" || port == "80" {
		return fmt.Sprintf("%s://%s", proto, host)
	}
	return fmt.Sprintf("%s://%s", proto, net.JoinHostPort(host, port))
}
