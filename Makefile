main_package_path = ./cmd
binary_name = home-solar-pi-server

env_file = ../.env


.PHONY: build
build:
	CGO_ENABLED=0 go build -o=./dist/${binary_name} ${main_package_path}
	cp ${env_file} dist/.env

.PHONY: build-dev
build-dev:
	go build -o=./dist/${binary_name} ${main_package_path}
	cp ${env_file} dist/.env

.PHONY: run-dev
run-dev:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build-dev" --build.bin "./dist/${binary_name} ./dist/.env" --build.delay "100" \
		--build.exclude_dir "postgres-data" \
		--build.include_ext "go" \
		--misc.clean_on_exit "true"

.PHONY: build-docker
build-docker: build
	docker build -t home-solar-pi-server .
