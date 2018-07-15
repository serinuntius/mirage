GIT_VER := $(shell git describe --tags)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)

clean:
	rm -rf pkg/*

binary: clean
	gox -osarch="linux/amd64 darwin/amd64 windows/amd64 windows/386" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-X main.version=${GIT_VER} -X main.buildDate=${DATE}"

package: binary
	cd ./pkg && find . -name "*${GIT_VER}*" -type f \
         -exec mkdir -p mirage/data \;  \
         -exec cp {} mirage/mirage \;   \
         -exec cp -r ../html ../config_sample.yml mirage/ \; \
         -exec zip -r {}.zip mirage \;     \
         -exec rm -rf mirage \;

