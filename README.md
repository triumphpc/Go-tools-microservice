### Go, RabbitMQ and gRPC [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) microservice 👋


#### 👨‍💻 Full list what has been used:
* [GRPC](https://grpc.io/) - gRPC
* [RabbitMQ](https://github.com/streadway/amqp) - RabbitMQ
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker
* [Prometheus](https://prometheus.io/) - Prometheus
* [Grafana](https://grafana.com/) - Grafana
* [Jaeger](https://www.jaegertracing.io/) - Jaeger tracing
* [Bluemonday](https://github.com/microcosm-cc/bluemonday) - HTML sanitizer
* [Gomail](https://github.com/go-gomail/gomail/tree/v2) - Simple and efficient package to send emails
* [Go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) - Sql mock driver for golang to test database interactions
* [Go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware) - interceptor chaining, auth, logging, retries and more
* [Opentracing-go](https://github.com/opentracing/opentracing-go) - OpenTracing API for Go
* [Prometheus-go-client](https://github.com/prometheus/client_golang) - Prometheus instrumentation library for Go applications
* [Cadvisor for Container Metrics](https://prometheus.io/docs/guides/cadvisor/) - cAdvisor (short for container Advisor) analyzes and exposes resource usage and performance data from running containers
* [Echo](https://github.com/labstack/echo) - High performance, minimalist Go web framework
* [Evans](https://github.com/ktr0731/evans) - Evans has been created to use easier than other existing gRPC clients.
```shell
~/Devel/go/src/github.com/triumphpc/Go-tools-microservice (master*) » evans internal/email/proto/email.proto -p 5001                                                      triumphpc@MacBook-Pro-triumphpc

  ______
 |  ____|
 | |__    __   __   __ _   _ __    ___
 |  __|   \ \ / /  / _. | | '_ \  / __|
 | |____   \ V /  | (_| | | | | | \__ \
 |______|   \_/    \__,_| |_| |_| |___/

 more expressive universal gRPC client


emailService.EmailService@127.0.0.1:5001> show package
+--------------+
|   PACKAGE    |
+--------------+
| emailService |
+--------------+

emailService.EmailService@127.0.0.1:5001> pac
command pac: unknown command

emailService.EmailService@127.0.0.1:5001> package em
command package: unknown package name 'em'

emailService.EmailService@127.0.0.1:5001> package emailService

emailService@127.0.0.1:5001> show service
+--------------+----------------------+-----------------------------+------------------------------+
|   SERVICE    |         RPC          |        REQUEST TYPE         |        RESPONSE TYPE         |
+--------------+----------------------+-----------------------------+------------------------------+
| EmailService | SendEmails           | SendEmailRequest            | SendEmailResponse            |
| EmailService | FindEmailById        | FindEmailByIdRequest        | FindEmailByIdResponse        |
| EmailService | FindEmailsByReceiver | FindEmailsByReceiverRequest | FindEmailsByReceiverResponse |
+--------------+----------------------+-----------------------------+------------------------------+

emailService@127.0.0.1:5001> call SendEmails
command call: failed to get the RPC descriptor for: SendEmails: unknown service name

emailService@127.0.0.1:5001> se
command se: unknown command

emailService@127.0.0.1:5001> service  EmailService
```

* [Hey](https://github.com/rakyll/hey) - hey is a tiny program that sends some load to a web application.
```shell
hey -n 2 -c 2 -m GET https://yandex.ru                                                              triumphpc@MacBook-Pro-triumphpc

Summary:
  Total:	0.0722 secs
  Slowest:	0.0722 secs
  Fastest:	0.0649 secs
  Average:	0.0685 secs
  Requests/sec:	27.7053

  Total data:	16654 bytes
  Size/request:	8327 bytes

Response time histogram:
  0.065 [1]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.066 [0]	|
  0.066 [0]	|
  0.067 [0]	|
  0.068 [0]	|
  0.069 [0]	|
  0.069 [0]	|
  0.070 [0]	|
  0.071 [0]	|
  0.071 [0]	|
  0.072 [1]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10% in 0.0722 secs
  0% in 0.0000 secs
  0% in 0.0000 secs
  0% in 0.0000 secs
  0% in 0.0000 secs
  0% in 0.0000 secs
  0% in 0.0000 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0420 secs, 0.0649 secs, 0.0722 secs
  DNS-lookup:	0.0138 secs, 0.0138 secs, 0.0138 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0000 secs
  resp wait:	0.0120 secs, 0.0097 secs, 0.0144 secs
  resp read:	0.0016 secs, 0.0016 secs, 0.0016 secs

Status code distribution:
  [200]	2 responses

```


#### Recommendation for local development most comfortable usage:
    make local // run all containers
    make run // run the application

#### 🙌👨‍💻🚀 Docker-compose files:
    docker-compose.local.yml - run rabbitmq, postgresql, jaeger, prometheus, grafana containers
    docker-compose.yml - run all in docker

### Docker development usage:
    make docker

### Local development usage:
    make local
    make run

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### RabbitMQ UI:

http://localhost:15672


### Email example
```json
{
  "emailId" : "ebae32c84f744cfb9f94a100faca18d1",
  "to" : ["triumph.job@gmail.com"],
  "body" : "Hello <b>Bob</b> and <i>Cora</i>!",
  "subject" : "Subject Hello!", 
  "from" : "triumph.job@yandex.ru"
}
```
