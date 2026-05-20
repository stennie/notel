set shell := ["bash", "-cu"]

binary := "bin/notel"
dist_dir := "dist"
version := `git describe --tags --always --dirty 2>/dev/null || echo dev`

default:
    @just --list

fmt:
    find . -name '*.go' -not -path './vendor/*' -exec gofmt -w {} +

vet:
    go vet ./...

test:
    go test ./...

test-one TEST:
    go test ./... -run '{{TEST}}'

build:
    mkdir -p bin
    go build -ldflags "-X github.com/stennie/notel/cmd.Version={{version}}" -o {{binary}} .

release:
    bash ./scripts/release.sh "{{version}}" "{{dist_dir}}"

release-check:
    just check
    just release

run *args:
    go run . {{args}}

version:
    go run -ldflags "-X github.com/stennie/notel/cmd.Version={{version}}" . version

check:
    just fmt
    just vet
    go test ./...
    go build ./...

audit:
    go run . audit

audit-fix:
    go run . audit --fix

audit-verbose:
    go run . audit --verbose

list:
    go run . list

clean:
    rm -rf bin {{dist_dir}}
