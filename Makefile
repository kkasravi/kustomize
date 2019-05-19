# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
GCLOUD_PROJECT ?= kubeflow-images-public
GOLANG_VERSION ?= 1.12
GOPATH ?= $(HOME)/go
VERBOSE ?= 
PLUGIN_DIR ?= plugin/app.k8s.io/v1beta1
KIND ?= Application
export GO111MODULE = on
export GO = go

all: build

# Run go fmt against code
fmt:
	@$(GO) fmt ./...

# Run go vet against code
vet:
	@$(GO) vet ./...

$(GOPATH)/bin/deepcopy-gen:
	GO111MODULE=on $(GO) get k8s.io/code-generator/cmd/deepcopy-gen


build: fmt vet
	GO111MODULE=on $(GO) build -i -gcflags 'all=-N -l' -o bin/kustomize kustomize.go

install: build
	@echo copying bin/kustomize to /usr/local/bin
	@cp bin/kustomize /usr/local/bin

build-plugin:
	GO111MODULE=on $(GO) build -i -gcflags 'all=-N -l' -o $(PLUGIN_DIR)/$(KIND).so \
		-buildmode plugin -tags=plugin $(PLUGIN_DIR)/plugin.go
	cp ./plugin/app.k8s.io/v1beta1/Application.so /Users/kdkasrav/go/src/github.com/kubeflow/manifests/plugins/kustomize/plugin/app.k8s.io/v1beta1/Application.so

build-main:
	GO111MODULE=on $(GO) build -i -gcflags 'all=-N -l' -o $(PLUGIN_DIR)/test_app $(PLUGIN_DIR)/main.go
