include Makefile.include

.PHONY: all 
all: \
	validation/validation.pb.go \
	control/ansible/ansible.pb.go \
	control/bash/bash.pb.go \
	control/jupyter/jupyter.pb.go \
	log/log.pb.go \
	metadata/id/id.pb.go \
	metadata/file/file.pb.go \
	metadata/network/network.pb.go \
	metadata/os/os.pb.go 

validation/validation.pb.go: validation/validation.proto
	$(protoc-build)

control/ansible/ansible.pb.go: control/ansible/ansible.proto
	$(protoc-build)

control/bash/bash.pb.go: control/bash/bash.proto
	$(protoc-build)

control/jupyter/jupyter.pb.go: control/jupyter/jupyter.proto
	$(protoc-build)

log/log.pb.go: log/log.proto
	$(protoc-build)

metadata/id/id.pb.go: metadata/id/id.proto
	$(protoc-build)

metadata/file/file.pb.go: metadata/file/file.proto
	$(protoc-build)

metadata/network/network.pb.go: metadata/network/network.proto
	$(protoc-build)

metadata/os/os.pb.go: metadata/os/os.proto
	$(protoc-build)

.PHONY: clean
clean:
	find . -name "*.pb.go" -delete
