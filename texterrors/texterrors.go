package texterrors

import "log"

var (
	FailedConnectionWithRabbitMQ = "Failed to connect to RabbitMQ"
	FailedToOpenChannel          = "Failed to open a channel"
	QueueNotFound                = "Failed to send. Queue not found"
	NoMessageToSend              = "No message to send"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("% s:% s", msg, err)
	}
}
