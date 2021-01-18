# ENV
GOLANG_VERSION ?= 1.15
GOTAGS = gen
# Hooks
PRE_BUILD_TARGETS = packr
PRE_BUILD_RELEASE_TARGETS = packr
PRE_BUILD_DEBUG_TARGETS = packr
include main.mk
# Update main targets

GODOWNLOADER_VERSION = 0.1.0
PACKR_VERSION = 2.8.1

bin/packr: bin/packr-${PACKR_VERSION}
	@ln -sf packr-${PACKR_VERSION} bin/packr
bin/packr-${PACKR_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gobuffalo/packr/releases/download/v${PACKR_VERSION}/packr_${PACKR_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - packr2 > ./bin/packr-${PACKR_VERSION} && chmod +x ./bin/packr-${PACKR_VERSION}

.PHONY: packr
packr: bin/packr
	bin/packr

.PHONY: clean-packr
clean-packr: bin/packr
	bin/packr clean

.PHONY: update-link
update-link: bin/godownloader
	bin/godownloader --repo=vseinstrumentiru/lego > ./install2.sh

bin/godownloader: bin/godownloader-${GODOWNLOADER_VERSION}
	@ln -sf godownloader-${GODOWNLOADER_VERSION} bin/godownloader
bin/godownloader-${GODOWNLOADER_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/goreleaser/godownloader/releases/download/v${GODOWNLOADER_VERSION}/godownloader_${GODOWNLOADER_VERSION}_${OS}_x86_64.tar.gz | tar -zOxf - godownloader > ./bin/godownloader-${GODOWNLOADER_VERSION} && chmod +x ./bin/godownloader-${GODOWNLOADER_VERSION}
