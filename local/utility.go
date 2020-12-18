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
func GetURL(proto, prefix, host, port string) string {
	if port == "" {
		return fmt.Sprintf("%s://%s%s", proto, prefix, host)
	}
	return fmt.Sprintf("%s://%s%s", proto, prefix, net.JoinHostPort(host, port))
}
