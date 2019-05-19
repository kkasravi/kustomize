// +build plugin

//go:generate go run sigs.k8s.io/kustomize/cmd/pluginator
package main

import (
	application "sigs.k8s.io/application/pkg/apis/app/v1beta1"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr         ifc.Loader
	rf          *resmap.Factory
	Application *application.Application
	types.GeneratorOptions
}

var KustomizePlugin plugin

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
