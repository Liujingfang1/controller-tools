/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package project

import (
	"sigs.k8s.io/controller-tools/pkg/scaffold/input"
)

var _ input.File = &Makefile{}

// Makefile scaffolds the Makefile
type Makefile struct {
	input.Input
	// Image is controller manager image name
	Image string
}

// GetInput implements input.File
func (c *Makefile) GetInput() (input.Input, error) {
	if c.Path == "" {
		c.Path = "Makefile"
	}
	c.TemplateBody = makefileTemplate
	c.Input.IfExistsAction = input.Error
	return c.Input, nil
}

var makefileTemplate = `
all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager {{ .Repo }}/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/rbac
	kubectl apply -f config/crds
	kubectl apply -f config/manager

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go build -o /tmp/controller-gen sigs.k8s.io/controller-tools/cmd/controller-gen && /tmp/controller-gen all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t {{ .Image }}

# Push the docker image
docker-push:
	docker push {{ .Image }}
`
