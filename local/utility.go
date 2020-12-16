package local

import "strings"

//TrimDomain trim前面的www及前后的".
func TrimDomain(d string) string {
	return strings.Trim(strings.TrimPrefix(d, "www"), ".")

}
