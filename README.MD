# Go-Microservice-Template
A template for microservice development in Go. Most of the boilerplate code is already done. Setup and configuration is done with environment variables.  
Defaults are set for Redis, Postgres, Datadog, Sentry, Gin, Franz Kafka, and Redis.
## Dependencies

```shell
go get github.com/joho/godotenv
go get github.com/caarlos0/env/v11
go get github.com/jackc/pgx/v5
go get github.com/redis/go-redis/v9
go get github.com/getsentry/sentry-go
go get gopkg.in/DataDog/dd-trace-go.v1/ddtrace
go get gopkg.in/DataDog/dd-trace-go.v1/profiler
go get github.com/gin-gonic/gin
```