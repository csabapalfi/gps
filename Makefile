SHELL = /bin/bash

release:
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	git tag $(version)
	git push origin --tags
	mkdir -p releases/$(version)
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/gps-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/gps-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/gps-windows-amd64 .
ifndef RELEASE
	go get -u github.com/aktau/github-release
endif
	github-release release --user csabapalfi --repo gps --tag $(version) || true
	github-release upload --user csabapalfi --repo gps --tag $(version) --name gps-linux-amd64 --file releases/$(version)/gps-linux-amd64 || true
	github-release upload --user csabapalfi --repo gps --tag $(version) --name gps-darwin-amd64 --file releases/$(version)/gps-darwin-amd64 || true
	github-release upload --user csabapalfi --repo gps --tag $(version) --name gps-windows-amd64 --file releases/$(version)/gps-windows-amd64 || true
