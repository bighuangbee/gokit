package mqtt

import (
	"fmt"
	"sync"
	"time"

	"github.com/bighuangbee/gokit/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	QosLevel         byte = 1
	publishTimeout        = 3
	subscribeTimeout      = 3
	connectTimeout        = 3
)

type MessageHandler func(client mqtt.Client, msg mqtt.Message)

type subscribeEvent struct {
	Qos     byte
	Topic   string
	Handler MessageHandler
}

type Client struct {
	sdkClient        mqtt.Client
	logger           log.Logger
	rwMux            sync.RWMutex
	subscribeHistory map[string]subscribeEvent
}

type Option struct {
	ClientId string
	Addr     string
	Logger   log.Logger
	User     string
	Password string
}

func NewClient(opt Option) (*Client, error) {
	c := &Client{
		logger:           opt.Logger,
		subscribeHistory: make(map[string]subscribeEvent),
	}
	cli, err := c.newSdkClient(opt)
	if err != nil {
		return nil, err
	}
	c.sdkClient = cli
	return c, nil
}

func (c *Client) newSdkClient(opt Option) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", opt.Addr))
	opts.SetClientID(opt.ClientId)
	opts.SetOrderMatters(false)
	opts.SetKeepAlive(time.Second * 10)
	opts.ConnectTimeout = time.Second
	opts.WriteTimeout = time.Second
	opts.PingTimeout = time.Second
	opts.SetConnectRetryInterval(time.Second * 3)
	opts.AutoReconnect = true
	opts.ConnectRetry = true
	opts.Username = opt.User
	opts.Password = opt.Password

	opts.OnConnect = func(client mqtt.Client) {
		r := client.OptionsReader()
		c.logger.Infow("MQTT connected", "clientID", r.ClientID())
		c.reSubscribe(client)
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		r := client.OptionsReader()
		c.logger.Infof("MQTT connect lost", r.ClientID(), err)
	}
	opts.OnReconnecting = func(client mqtt.Client, opts *mqtt.ClientOptions) {
		c.logger.Infow("MQTT connected", "clientID", opts.ClientID)
	}
	opts.DefaultPublishHandler = func(client mqtt.Client, msg mqtt.Message) {
		c.logger.Infow("MQTT UNEXPECTED MESSAGE", "clientID", opts.ClientID, "msg", msg)
	}

	cli := mqtt.NewClient(opts)
	token := cli.Connect()
	if ok := token.WaitTimeout(connectTimeout * time.Second); !ok {
		return nil, errors.New("MQTT connect timeout")
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *Client) Publish(topic string, payload []byte) error {
	if !c.sdkClient.IsConnected() {
		return errors.New("MQTT not connected")
	}
	t := c.sdkClient.Publish(topic, QosLevel, false, payload)
	if ok := t.WaitTimeout(publishTimeout * time.Second); !ok {
		return errors.New("MQTT publish timeout")
	}
	return t.Error()
}

func (c *Client) PublishAsync(topic string, payload []byte, callback func(error)) {
	if !c.sdkClient.IsConnected() {
		callback(errors.New("MQTT not connected"))
		return
	}
	t := c.sdkClient.Publish(topic, QosLevel, false, payload)
	go func() {
		if ok := t.WaitTimeout(publishTimeout * time.Second); !ok {
			callback(errors.New("MQTT publish timeout"))
		} else {
			callback(t.Error())
		}
	}()
}

func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	c.logger.Infow("subscribe topic", "topic", topic)
	c.addSubscribeHistory(topic, QosLevel, handler)
	return c.subscribe(c.sdkClient, topic, QosLevel, wrapMessageHandler(handler))
}

func (c *Client) subscribe(cli mqtt.Client, topic string, qos byte, handler mqtt.MessageHandler) error {
	t := cli.Subscribe(topic, qos, handler)
	if ok := t.WaitTimeout(subscribeTimeout * time.Second); !ok {
		return errors.New("MQTT subscribe timeout")
	}
	return t.Error()
}

func (c *Client) addSubscribeHistory(topic string, qos byte, handler MessageHandler) {
	c.rwMux.Lock()
	c.subscribeHistory[topic] = subscribeEvent{
		Qos:     qos,
		Topic:   topic,
		Handler: handler,
	}
	c.rwMux.Unlock()
}

func (c *Client) reSubscribe(cli mqtt.Client) {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()
	for _, sub := range c.subscribeHistory {
		if err := c.subscribe(cli, sub.Topic, sub.Qos, wrapMessageHandler(sub.Handler)); err != nil {
			c.logger.Errorf("MQTT reSubscribe error: %v", err)
		}
	}
}

func wrapMessageHandler(fn MessageHandler) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		fn(cli, msg)
	}
}
