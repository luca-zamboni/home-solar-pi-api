
main_package_path = ./cmd
binary_name = home-solar


.PHONY: build
build:
	go build -o=./bin/${binary_name} ${main_package_path}

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "postgres-data" \
		--build.include_ext "go, env" \
		--misc.clean_on_exit "true"
