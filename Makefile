
TAG := "help_ua_nl_bot"

all: build

build:
	docker build --tag "$(TAG)" .

clean:

run:
	docker run -e HELP_UA_NL_BOT_TOKEN "$(TAG)"
