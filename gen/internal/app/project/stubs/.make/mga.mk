# Dependency versions
MGA_VERSION = 0.4.2

.PHONY: generate
generate: bin/mga ## Generate code
	go generate -x ./...
	bin/mga generate kit endpoint ./internal/app/../...
	bin/mga generate event handler --output subpkg:suffix=gen ./internal/app/../...
	bin/mga generate event dispatcher --output subpkg:suffix=gen ./internal/app/../...
	bin/mga generate testify mock ./internal/app/../...

bin/mga-${MGA_VERSION}:
	@mkdir -p bin
	curl -sfL https://git.io/mgatool | bash -s v${MGA_VERSION}
	@mv bin/mga $@

bin/mga: bin/mga-${MGA_VERSION}
	@ln -sf mga-${MGA_VERSION} bin/mga
