package services

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/micro-plat/hydra/registry"
)

type Domain struct {
	Domain string `form:"domain" json:"domain" valid:"dns,required"`
	IP     string `form:"ip" json:"ip" valid:"ip,required"`
}

func checkAndCreate(domain *Domain, registry registry.IRegistry) interface{} {

	root := filepath.Join("/dns", domain.Domain)
	path := filepath.Join(root, domain.IP)

	b, err := registry.Exists(root)
	if err != nil {
		return err
	}
	if b {
		paths, _, err := registry.GetChildren(root)
		if err != nil {
			return err
		}
		for _, pc := range paths {
			if err := registry.Delete(filepath.Join(root, pc)); err != nil {
				return err
			}
		}
	}
	if err := registry.CreatePersistentNode(path,
		fmt.Sprintf(`{"time":%d}`, time.Now().Unix())); err != nil {
		return err
	}
	return "success"
}
