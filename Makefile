COMMIT = $$(git describe --always)

generate:
	@go generate ./...

build: generate
	@echo "====> Build init"
	@sh -c ./build.sh
