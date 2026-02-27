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

## Component system
Component system gives ability to create component based application. Each component can be dependend on other components.
It's solves problem, order of component initialization. Like you have Kafka, Postgres and you need Kafka and Postgres started
before your web application start. Component system can handle this problem.

Possible fail scenarios:
* What to do if one child component fails?
* Should app shutdown if child component fails (parent or child)?
* Shutdown all components if error occur in parent
* Should app retry if any child component fails?
* How retry logic should be implemented (shoutdown all child components, then parent and restart parent and then child components)?

All of these are possible scenarios hard to address properly. There is a lot of possible scenarios and to cover all of them is not possible. 
Solution can easy became way too complex then problem is trying to solve.  

Default approach is in case any error occurs in parent or child component, shutdown all components in proper order, from child to parent. 
Except default approach, there another strategy to handle component it's restart when component get signal it can restart full hierarchy. 
It's cover possible scenario when component like Kafka need to restart when configuration of topics is changed