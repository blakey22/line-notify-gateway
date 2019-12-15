TAG=$(shell git describe --abbrev=0 --tags)

.PHONY: image
image:
	docker build -t blakey22/line-notify-gateway:$(TAG) .