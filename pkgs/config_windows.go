package pkgs

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	HOST_FILE string
	NAME_FILE string

<<<<<<< HEAD:pkgs/config_windows.go
=======
	registrykey=`SYSTEM\CurrentControlSet\services\Tcpip\Parameters\Interfaces`
>>>>>>> dev1.2-ljy:dns/pkgs/config_windows.go
)

func init() {
	root := filepath.VolumeName(os.Getenv("SYSTEMROOT"))
	HOST_FILE = fmt.Sprintf(`%s:\Windows\System32\drivers\etc\hosts*`, root)
	NAME_FILE = fmt.Sprintf(`%s:\Windows\System32\drivers\etc\names.conf`, root)
}
