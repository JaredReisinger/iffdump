.DEFAULT_GOAL := build

IMAGE_NAME := jaredreisinger/xyzzybot
IMAGE_VERSION := 0.1

build:
	go build -v .
.PHONY: build

install: build
	go install -v .
.PHONY: install

image:
	docker build \
		-t ${IMAGE_NAME}:${IMAGE_VERSION} \
		-t ${IMAGE_NAME}:latest \
		.
.PHONY: image

acquire-external-tools:
	go get -u github.com/kardianos/govendor
.PHONY: acquire-external-tools

update-dependencies:
	govendor fetch +outside
.PHONY: update-dependencies

try: build
	./iffdump -verbose ../xyzzybot/_working-root/D647607DY/savegame.qut
	@# ./iffdump ../xyzzybot/_working-root/D647607DY/savegame.qut
.PHONY: try

# shell:
# 	docker run --rm \
# 		--tty \
# 		--interactive \
# 		--volume ${PWD}/config:/usr/local/etc/xyzzybot \
# 		--volume ${PWD}/sample-games:/usr/local/games \
# 		${IMAGE_NAME}:${IMAGE_VERSION}

# lint:
# 	go lint -v ./...
# .PHONY: lint
