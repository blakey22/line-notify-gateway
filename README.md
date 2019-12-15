# Line Notify Gateway
Notification Messages via [LINE Notify](https://notify-bot.line.me/en/)

Currently, this gateway supports the following webhooks:
* Prometheus
* Github
* Docker Hub

Supported localized messages:
* en_US
* zh_TW

Supported environment variables:
* `SECRET`: the access token of this gateway
* `TOKEN`: token of LINE notify
* `LOCALE`: locale of notification messages

## Authorization
By setting environment variable (`SECRET`) or command argument (`-s`), the line notify gateway will verify HTTP header: `Authorization: Bearer <secret>` to ensure basic access control.


## Prometheus
You can set a Prometheus Alert Manager Webhook to receive LINE notification of alerts.

URL: `<host>:<port>/prometheus`

You can refer to the demo setting: [alertmanager.yml](./demo/alertmanager.yml)

## Github
You can set a Github Webhook to receive LINE notification about repository changes.

URL: `<your_ip>:<port>/github`

**Note: While you're setting Github Webhook, please set `Content type` to `application/json`**

Supported Events:
* create
* delete
* push
* pull_request


## Docker Hub
You can set a Docker Hub Webhook to receive LINE notification about an event of new docker image is pushed.

URL: `<your_ip>:<port>/docker-hub`


## Tester
You can use it to send a test message via LINE notify

URL: `<host>:<port>/tester`

Example:
```
curl -d "test message" -X POST http://127.0.0.1:8080/tester
```


## Demo (Prometheus)
You will need `docker` and `docker-compose` to run this demo.

1. Please create a `docker-compose.override.yml` file to provide the LINE notify token and gateway secret.
```
version: '3.7'

services:
 line-notify-gateway:
 environment:
 - TOKEN=<your_line_notify_token>
 - SECRET=93944dfd-d476-446e-be73-7bb62c1e0446
```

2. Run `docker-compose up` and wait around 30 seconds, you should able to see the LINE notification. :)

## Usage
You can run `./line-notify-gateway -h` to get the help message

```
Usage:
  line-notify-gateway [OPTIONS]

Application Options:
  -h, --host=           Bind host (default: 0.0.0.0) [$HOST]
  -p, --port=           Bind port (default: 8080) [$PORT]
  -s, --secret=         Secret token for accessing the gateway [$SECRET]
  -t, --token=          LINE notify token [$TOKEN]
      --endpoint=       LINE notify endpoint (default:
                        https://notify-api.line.me/api/notify) [$ENDPOINT]
      --locale=         Locale of LINE message [$LOCALE]
      --templates-path= Path of message templates (default: ./templates)
                        [$TEMPLATES_PATH]

Help Options:
  -h, --help            Show this help message
```

## References:
* LINE Notify: https://notify-bot.line.me/en/
* Prometheus Alert Manager Webhook: https://prometheus.io/docs/alerting/configuration/#webhook_config
* Github Webhook: https://developer.github.com/webhooks/
* Docker Hub Webhook: https://docs.docker.com/docker-hub/webhooks/
