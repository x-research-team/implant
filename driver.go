package implant

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"plugin"

	"github.com/x-research-team/contract"
)

const (
	so = "so"

	all  = "**"
	tmpl = "%s/*.%s"
)

var (
	plugins = make(map[string]*plugin.Plugin)
	modules contract.KernelModules
)

func Init(ss ...string) {
	root, err := filepath.Abs("../")
	if err != nil {
		log.Printf("[ERR] [PLUGIN] Error: %v\n", err)
	}

	for _, s := range ss {
		libs, err := filepath.Glob(fmt.Sprintf(tmpl, path.Join(root, s, all), so))
		if err != nil {
			log.Printf("[ERR] [PLUGIN] Error: %v\n", err)
		}

		for _, lib := range libs {
			plugins[lib], err = plugin.Open(lib)
			if err != nil {
				log.Printf("[ERR] [PLUGIN] Error: %v\n", err)
				continue
			}

			component, err := plugins[lib].Lookup("Init")
			if err != nil {
				log.Printf("[ERR] [PLUGIN] Error: %v\n", err)
				continue
			}

			c := component.(func() contract.KernelModule)()
			modules = append(modules, c)
		}
	}
}

func Modules() contract.KernelModules {
	return modules
}
