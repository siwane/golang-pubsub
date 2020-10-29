module web-app

go 1.15

require (
	cloud.google.com/go/pubsub v1.8.2
	github.com/joho/godotenv v1.3.0
)

replace (
	web-app/utils => ./utils/env.go
	web-app/pubsub/publisher => ./pubsub/publisher/publisher.go
	web-app/pubsub/subscriber => ./pubsub/subscriber/subscriber.go
)
