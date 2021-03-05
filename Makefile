GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build: fmt
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/vault-plugin-auth-tfe cmd/vault-plugin-auth-tfe/main.go

start:
	# run this first -->>> eval $(doormat aws --account se_demos_dev)
	vault server -dev -dev-root-token-id=root \
	-dev-plugin-dir=./vault/plugins -log-level=debug \
	-dev-listen-address=192.168.178.40:8200

enable:
# fad4d28b6f57ca6a1acd49b948e0a279d805280c461bb29fcb8781e57c1c3562
	# vault plugin register -sha256=fad4d28b6f57ca6a1acd49b948e0a279d805280c461bb29fcb8781e57c1c3562 auth vault-plugin-auth-tfe
	vault auth enable -path=tfe-auth vault-plugin-auth-tfe

clean:
	rm -f ./vault/plugins/vault-plugin-auth-tfe

fmt:
	go fmt $$(go list ./...)

test:
	# vault write -force auth/tfe-auth/config
	vault write auth/tfe-auth/config organization=org2 terraform_host=tfe.ric.gcp.hashidemos.io
	vault read auth/tfe-auth/config
	vault write auth/tfe-auth/role/role1 workspaces=123,456
	vault write auth/tfe-auth/role/role2 workspaces=* policies=default
	vault list auth/tfe-auth/role
	vault read auth/tfe-auth/role/role1

	# vault write auth/tfe-auth/login role=role2 workspace=aaa run-id=aa atlas-token=aa
	vault write auth/tfe-auth/login role=role2 \
		workspace=tfe-gcp-test-network \
		run-id=run-U7VpRnrDSGhyk8Ff \
		atlas-token=eJ1fkmbxGLtbNg.atlasv1.GGpvS5FwHsYTBLze9S4Pqsx2ahPc67Ypv8d5XlgHptWQ06dwHRrtnXWb2tyTzIp0860


.PHONY: build clean fmt start enable test
