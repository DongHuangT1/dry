LDFLAGS=-trimpath -ldflags="-s -w"

define build
	CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build $(LDFLAGS) -o bin/dry_$1_$2$3
endef

all: linux darwin windows checksum

linux:
	$(call build,linux,386)
	$(call build,linux,arm)
	$(call build,linux,arm64)
	$(call build,linux,amd64)

darwin:
	$(call build,darwin,arm64)
	$(call build,darwin,amd64)

windows:
	$(call build,windows,386,.exe)
	$(call build,windows,arm64,.exe)
	$(call build,windows,amd64,.exe)

checksum:
	cd bin && sha256sum * > dry_checksums.txt

clean:
	rm -rf bin