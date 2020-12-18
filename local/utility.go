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
	return strings.HasPrefix(d, "wwww.")
}
func GetURL(proto string, host string, port string) string {
	if port == "" {
		return fmt.Sprintf("%s://%s", proto, host)
	}
	return fmt.Sprintf("%s://%s", proto, net.JoinHostPort(host, port))
}
