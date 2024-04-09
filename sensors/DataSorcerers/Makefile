include Makefile.include


.PHONY: all 
all: \
	build/discern-validation-sorcerer \
	build/discern-ansible-sorcerer \
	build/discern-bash-sorcerer \
	build/discern-jupyter-sorcerer \
	build/discern-logs-sorcerer \
	build/discern-id-sorcerer \
	build/discern-id-server-sorcerer \
	build/discern-file-sorcerer \
	build/discern-network-sorcerer \
	build/discern-os-sorcerer

build/discern-validation-sorcerer: validation/*.go
	$(call go-build,validation,$(realpath $(dir .)))

build/discern-ansible-sorcerer: control/ansible/*.go
	$(call go-build,control/ansible,$(realpath $(dir .)))

build/discern-bash-sorcerer: control/bash/*.go
	$(call go-build,control/bash,$(realpath $(dir .)))

build/discern-jupyter-sorcerer: control/jupyter/*.go
	$(call go-build,control/jupyter,$(realpath $(dir .)))

build/discern-logs-sorcerer: logs/*.go
	$(call go-build,logs,$(realpath $(dir .)))

build/discern-id-sorcerer: metadata/id/*.go
	$(call go-build,metadata/id,$(realpath $(dir .)))

build/discern-id-server-sorcerer: metadata/id/server/*.go
	$(call go-build,metadata/id/server,$(realpath $(dir .)))

build/discern-file-sorcerer: metadata/file/*.go
	$(call go-build,metadata/file,$(realpath $(dir .)))

build/discern-network-sorcerer: metadata/network/*.go
	$(call cgo-build,metadata/network,$(realpath $(dir .)))

build/discern-os-sorcerer: metadata/os/*.go
	$(call go-build,metadata/os,$(realpath $(dir .)))

.PHONY: clean
clean:
	rm -f build/*
