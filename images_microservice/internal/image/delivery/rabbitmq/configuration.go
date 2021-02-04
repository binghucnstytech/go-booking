package rabbitmq

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false

	ImagesExchange = "images"

	ResizeQueueName   = "resize_queue"
	ResizeConsumerTag = "resize_consumer"
	ResizeWorkers     = 10
	ResizeBindingKey  = "resize_image_key"

	CreateQueueName   = "create_queue"
	CreateConsumerTag = "create_consumer"
	CreateWorkers     = 5
	CreateBindingKey  = "create_image_key"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_incoming_messages_total",
		Help: "The total number of incoming RabbitMQ messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_success_messages_total",
		Help: "The total number of success incoming success RabbitMQ messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_error_messages_total",
		Help: "The total number of error incoming success RabbitMQ messages",
	})
)
