package rabbitmq

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMQ struct {
	ctx      context.Context
	alarm    Alarm
	consumer *consumer
	producer *producer
	handle   func(context.Context, *amqp.Delivery) Action
	config   *Config
}

// WithContext sets the context for the RabbitMQ instance and returns it for method chaining.
func (r *RabbitMQ) WithContext(ctx context.Context) *RabbitMQ {
	r.ctx = ctx
	return r
}

// SetAlarm sets the alarm handler for the RabbitMQ instance and returns it for method chaining.
func (r *RabbitMQ) SetAlarm(alarm Alarm) *RabbitMQ {
	r.alarm = alarm
	return r
}

// SetConfig sets the RabbitMQ configuration and returns it for method chaining.
func (r *RabbitMQ) SetConfig(config *Config) *RabbitMQ {
	r.config = config
	return r
}

// SetHandle sets the message handler function for the RabbitMQ instance and returns it for method chaining.
func (r *RabbitMQ) SetHandle(handle func(context.Context, *amqp.Delivery) Action) *RabbitMQ {
	r.handle = handle
	return r
}

// Validate checks if the RabbitMQ instance has been properly configured.
// Returns an error if config is nil. Sets context to TODO() if it's nil.
func (r *RabbitMQ) Validate() error {
	if r.ctx == nil {
		r.ctx = context.TODO()
	}
	if r.config == nil {
		return errors.New("config is nil")
	}
	return nil
}

// StartConsumer starts the RabbitMQ consumer in a goroutine.
// Panics if validation fails, if there's no message handler, or if consumer initialization fails.
func (r *RabbitMQ) StartConsumer() {
	if err := r.Validate(); err != nil {
		panic(err)
	}
	// 初始化队列
	r.consumer = &consumer{
		ctx:   r.ctx,
		alarm: r.alarm,
		channel: &channel{
			config: r.config,
		},
	}
	// 检查是否有 Channel
	if err := r.consumer.Validate(); err != nil {
		panic(err)
		return
	}
	if r.handle == nil {
		panic(errors.New("handle is nil"))
		return
	}
	go r.consumer.Run(r.handle)
	go r.consumer.gracefulShutdown()
}

// StartProducer starts the RabbitMQ producer and returns a channel for sending messages.
// Returns an error if validation or channel initialization fails.
func (r *RabbitMQ) StartProducer() (chan<- []byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	var err error
	r.producer = &producer{
		ctx:            r.ctx,
		key:            r.config.BindKey,
		exchange:       r.config.ExchangeName,
		deliveryMode:   r.config.DeliveryMode,
		sendBodyLength: 4096,
		alarm:          r.alarm,
	}
	r.producer.channel = &channel{
		config: r.config,
	}
	r.config.ChannelNum = 1
	err = r.producer.channel.Init()
	if err != nil {
		return nil, err
	}
	err = r.producer.Validate()
	if err != nil {
		return nil, err
	}

	go r.producer.Run()
	time.Sleep(time.Second)
	return r.producer.sendBody, nil
}
