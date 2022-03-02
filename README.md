# Телеграм-бот для організації волонтерської допомоги Україні в Нідерландах

[https://t.me/help_ua_nl_bot](https://t.me/help_ua_nl_bot)

## How To

Have Golang installed and ready
### Debug mode

Create new Telegram bot [https://t.me/BotFather](https://t.me/BotFather) and save the secret <token>

~~~bash
    export TEST_HELP_UA_NL_BOT_TOKEN=<token>
    go mod download
    go run main.go
~~~

### Deploy

~~~bash
    make deploy
~~~
