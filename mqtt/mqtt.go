package mqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Option struct {
	Addr     string
	User     string
	Password string
}

type Client struct {
	mqtt mqtt.Client
}

const timeout = 3
const QosLevel byte = 0
const retained = false

func NewClient(opt *Option) (*Client, error) {
	cliOps := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", opt.Addr))
	cliOps.SetUsername(opt.User)
	cliOps.SetPassword(opt.Password)
	cliOps.SetCleanSession(true)
	cliOps.SetAutoReconnect(true)
	cliOps.SetConnectRetry(true)
	cliOps.SetKeepAlive(30 * time.Second)

	mqttCli := mqtt.NewClient(cliOps)
	token := mqttCli.Connect()

	if ok := token.WaitTimeout(timeout * time.Second); !ok {
		return nil, errors.New("connect timeout")
	}
	if token.Error() != nil {
		return nil, token.Error()
	}
	return &Client{mqtt: mqttCli}, nil
}
func (c *Client) Subscribe(topic string, callback mqtt.MessageHandler) error {
	t := c.mqtt.Subscribe(topic, QosLevel, callback)
	if ok := t.WaitTimeout(timeout * time.Second); !ok {
		return errors.New("subscribe timeout")
	}
	return t.Error()
}

func (c *Client) Publish(topic string, payload []byte) error {
	t := c.mqtt.Publish(topic, QosLevel, retained, payload)
	if ok := t.WaitTimeout(timeout * time.Second); !ok {
		return errors.New("Publish timeout")
	}
	return t.Error()
}
