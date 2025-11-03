module github.com/streamverse/search-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/elastic/go-elasticsearch/v8 v8.11.0
	github.com/streamverse/common-go v0.0.0
	go.uber.org/zap v1.26.0
)

replace github.com/streamverse/common-go => ../../packages/common-go

