package rabbitmq

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type producer struct {
	ctx            context.Context
	exchange       string
	key            string
	contentType    string
	deliveryMode   uint8 //2
	msgType        string
	channel        *channel
	sendBodyLength int
	sendBody       chan []byte
	alarm          Alarm
}

func (p *producer) Validate() error {
	if p.exchange == "" {
		return errors.New("exchange is nil,please check it")
	}
	if p.contentType == "" {
		p.contentType = "text/plain"
	}
	if p.channel.config.ChannelNum == 0 {
		p.channel.config.ChannelNum = 1
	}
	if p.channel == nil {
		return errors.New("channel is nil,please init channel")
	}
	if p.sendBodyLength == 0 {
		// 4096 is the default message queue buffer size (number of messages buffered in memory)
		p.sendBodyLength = 4096
	}
	p.sendBody = make(chan []byte, p.sendBodyLength)
	return nil
}

func (p *producer) Run() {
	go func() {
		for {
			select {
			case msg := <-p.sendBody:
				p.Send(msg, p.channel)
			}
		}
	}()
}

func (p *producer) Send(body []byte, channel *channel) {
	retryCount := 0
	// maxReconnectCount limits reconnection attempts before triggering an alarm notification
	maxReconnectCount := 3
	for {
		err := channel.Chan[0].PublishWithContext(
			p.ctx,
			p.exchange,
			p.key,
			false,
			false,
			amqp.Publishing{
				ContentType:  p.contentType,
				DeliveryMode: p.deliveryMode,
				Body:         body,
				Type:         p.msgType,
			})
		if err == amqp.ErrClosed {
			if err := channel.Init(); err != nil {
				retryCount++
				continue
			}
		}
		if err != nil {
			if retryCount >= maxReconnectCount {
				slog.Error("RabbitMQ.Producer", slog.String("info", err.Error()))
				if p.alarm != nil {
					if err := p.alarm.SetMsg(map[string]string{
						"Title":   "RabbitMQ-Producer 连接失败超出阈值",
						"Address": p.channel.config.Address,
						"Queue":   p.channel.config.QueueName,
					}); err != nil {
						slog.Error("RabbitMQ.Producer alarm.SetMsg", slog.String("error", err.Error()))
					}
					if err := p.alarm.Do(); err != nil {
						slog.Error("RabbitMQ.Producer alarm.Do", slog.String("error", err.Error()))
					}
				}
				return
			}
			retryCount++
			continue
		}
		return
	}
}
