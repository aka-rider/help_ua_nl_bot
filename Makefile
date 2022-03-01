
TAG := "help_ua_nl_bot"

all: build

build:
	docker build --tag "$(TAG)" .

clean:

run:
	docker run -e HELP_UA_NL_BOT_TOKEN "$(TAG)"

deploy: build
	echo "${SCALEWAY_SECRET_KEY}" | docker login "${SCALEWAY_CONTAINER_REGISTRY}" -u nologin --password-stdin
	docker tag "$(TAG):latest" "${SCALEWAY_CONTAINER_REGISTRY}/$(TAG):latest"
	docker push "${SCALEWAY_CONTAINER_REGISTRY}/$(TAG):latest"
