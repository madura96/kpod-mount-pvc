
.DEFAULT_GOAL := build

SHELL = /bin/bash # Force to use bash. Else, it uses sh and, for instance, time command is not available
CURRENT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null)

MAJOR := $(shell echo $(CURRENT_TAG) | cut -d'.' -f 1)
$(eval NEXT_MAJOR := $(shell echo $$(( $(MAJOR) + 1 ))))

MINOR := $(shell echo $(CURRENT_TAG) | cut -d'.' -f 2)
$(eval NEXT_MINOR := $(shell echo $$(( $(MINOR) +1 ))))

FIX := $(shell echo $(CURRENT_TAG) | cut -d'.' -f 3)
$(eval NEXT_FIX := $(shell echo $$(( $(FIX) + 1 ))))

ifeq ($(CURRENT_TAG),)
        MAJOR := "0"
        MINOR := "0"
        FIX := "0"
        NEXT_FIX = "1"
        NEXT_MINOR = "1"
        NEXT_MAJOR = "1"
endif

# https://stackoverflow.com/questions/10858261/how-to-abort-makefile-if-variable-not-set
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

.PHONY: major-tag minor-tag fix-tag clean build build-release
major-tag:
	git tag "$(NEXT_MAJOR).0.0"
	git push origin "$(NEXT_MAJOR).0.0"

minor-tag:
	git tag "$(MAJOR).$(NEXT_MINOR).0"
	git push origin "$(MAJOR).$(NEXT_MINOR).0"

fix-tag:
	git tag "$(MAJOR).$(MINOR).$(NEXT_FIX)"
	git push origin "$(MAJOR).$(MINOR).$(NEXT_FIX)"

build:
	@go mod tidy
	@go mod vendor
	@go build -mod vendor

build-release:
	@go mod tidy
	@go mod vendor
	@go build -mod vendor -ldflags="-X 'kpod-mount-pvc/cmd.Version=$(CURRENT_TAG)'"

clean:
	rm kpod-mount-pvc

tests:
	@go test ./...
