CR=ghcr.io
IMAGE=ROUKIEN/ram
TAG?=latest

build:
	docker build -t $(CR)/$(IMAGE):$(TAG) .
