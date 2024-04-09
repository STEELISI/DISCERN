include Makefile.include

.PHONY: all 
all: build/FusionCore

build-deps := $(shell find . -name "*.go")
build/FusionCore: $(build-deps)
	$(call go-build,$@)

.PHONY: clean
clean:
	rm -f build/FusionCore

.PHONY: install
install: install-fusioncore install-fusioncore-service install-fusioncore-config

.PHONY: install-fusioncore
install-fusioncore: build/FusionCore
	install -D $< /usr/local/bin/FusionCore

.PHONY: install-fusioncore-service
install-fusioncore-service: files/FusionCore.service
	install -D $< /etc/systemd/system/FusionCore.service
	systemctl daemon-reload

.PHONY: install-fusioncore-config
install-fusioncore-config: files/CoreConfig.yaml
	install -D $< /etc/discern/CoreConfig.yaml 

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/FusionCore
	rm -f /etc/systemd/system/FusinCore.service
	rm -f /etc/discern/CoreConfig.yaml
