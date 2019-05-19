package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"os"
	application "sigs.k8s.io/application/pkg/apis/app/v1beta1"
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	kvplugin "sigs.k8s.io/kustomize/k8sdeps/kv/plugin"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
)

var KustomizePlugin plugin

func main() {
	args := os.Args[1:]
	targetPath := ""
	if len(args) > 0 {
		targetPath = args[0]
	}
	fsys := fs.MakeRealFS()
	_loader, loaderErr := loader.NewLoader(loader.RestrictionRootOnly, targetPath, fsys)
	if loaderErr != nil {
		log.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	buf, bufErr := _loader.Load(targetPath)
	if bufErr != nil {
		log.Fatalf("could not load file: %v Error %v", targetPath, bufErr)
	}
	err := yaml.Unmarshal(buf, KustomizePlugin.Application)
	if err != nil {
		log.Fatalf("could not unmarshal file: %v Error %v", targetPath, err)

	}
	_resmap := resmap.NewFactory(resource.NewFactory(
		kunstruct.NewKunstructuredFactoryWithGeneratorArgs(
			&types.GeneratorMetaArgs{
				PluginConfig: kvplugin.ActivePluginConfig(),
			})))

	KustomizePlugin.Config(_loader, _resmap, nil)
}

type plugin struct {
	Application *application.Application
	options     types.GeneratorOptions
}

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, k ifc.Kunstructured) error {
	var buf []byte
	var err error
	buf, err = k.MarshalJSON()
	if err != nil {
		return fmt.Errorf("cannot marshal Kunstructured %v", err)
	}
	p.Application = &application.Application{}
	specErr := yaml.Unmarshal(buf, p.Application)
	if specErr != nil {
		return fmt.Errorf("cannot unmarshal Kunstructured %v", specErr)
	}
	return nil
}

func (p *plugin) Generate() (resmap.ResMap, error) {

	return nil, nil
}
