package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	application "sigs.k8s.io/application/pkg/apis/app/v1beta1"
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	kvplugin "sigs.k8s.io/kustomize/k8sdeps/kv/plugin"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

var KustomizePlugin plugin

func main() {
	args := os.Args[1:]
	targetPath := ""
	if len(args) > 0 {
		targetPath = args[0]
	}
	targetDir := filepath.Dir(targetPath)
	fsys := fs.MakeRealFS()
	_loader, loaderErr := loader.NewLoader(loader.RestrictionNone, targetDir, fsys)
	if loaderErr != nil {
		log.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	buf, bufErr := _loader.Load(targetPath)
	if bufErr != nil {
		log.Fatalf("could not load file: %v Error %v", targetPath, bufErr)
	}
	_resmapF := resmap.NewFactory(resource.NewFactory(
		kunstruct.NewKunstructuredFactoryWithGeneratorArgs(
			&types.GeneratorMetaArgs{
				PluginConfig: kvplugin.ActivePluginConfig(),
			})))

	err := KustomizePlugin.Config(_loader, _resmapF, buf)
	if err != nil {
		log.Fatalf("KustomizePlugin.Config returned %v", err)
	}
	_resmap, resmapErr := KustomizePlugin.Generate()
	if resmapErr != nil {
		log.Fatalf("KustomizePlugin.Generate returned %v", err)
	}
	out, outErr := _resmap.EncodeAsYaml()
	if outErr != nil {
		log.Fatalf("could not generate yaml %v", outErr)
	}
	os.Stdout.Write(out)
}

type plugin struct {
	ldr         ifc.Loader
	rf          *resmap.Factory
	Application *application.Application
	types.GeneratorOptions
}

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, buf []byte) error {
	p.Application = &application.Application{}
	p.GeneratorOptions = types.GeneratorOptions{}
	p.ldr = ldr
	p.rf = rf
	return yaml.Unmarshal(buf, p.Application)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	buf, err := yaml.Marshal(p.Application)
	if err != nil {
		return nil, err
	}
	return p.rf.NewResMapFromBytes(buf)
}
