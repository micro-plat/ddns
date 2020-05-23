package services

import (
	"fmt"
	"time"

	"github.com/micro-plat/hydra/registry"
)

type Domain struct {
	Domain string `form:"domain" json:"domain" valid:"dns,required"`
	IP     string `form:"ip" json:"ip" valid:"ip,required"`
}

func checkAndCreate(domain *Domain, r registry.IRegistry) error {
	root := registry.Join("/dns", domain.Domain)
	path := registry.Join(root, domain.IP)

	b, err := r.Exists(root)
	if err != nil {
		return err
	}
	if b {
		paths, _, err := r.GetChildren(root)
		if err != nil {
			return err
		}
		for _, pc := range paths {
			if err := r.Delete(registry.Join(root, pc)); err != nil {
				return err
			}
		}
	}
	if err := r.CreatePersistentNode(path,
		fmt.Sprintf(`{"time":%d}`, time.Now().Unix())); err != nil {
		return err
	}
	return nil
}
