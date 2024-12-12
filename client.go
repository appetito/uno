package uno

import (
	"encoding/json"
	"strconv"
	"time"

	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var UnoClientDefaultTimeout = time.Second * 2


type UnoClientConfig struct {
	TimeOut time.Duration
}

type UnoClient struct {
	nc *nats.Conn
	cfg *UnoClientConfig
}

type RequestOptions struct{
	Timeout time.Duration
}

func NewUnoClient(nc *nats.Conn, cfg *UnoClientConfig) (*UnoClient) {
	if cfg == nil {
		cfg = &UnoClientConfig{
			TimeOut: UnoClientDefaultTimeout,
		}
	}else{
		if cfg.TimeOut == 0{
			cfg.TimeOut = UnoClientDefaultTimeout
		}
	}
	return &UnoClient{
		nc: nc,
		cfg: cfg,
	}
}


func (c *UnoClient) RequestJSON(ctx context.Context, subject string, obj interface{}) (*nats.Msg, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	msg := nats.NewMsg(subject)
	msg.Data = data
	d, ok := ctx.Deadline()
	if !ok{
		d = time.Now().Add(c.cfg.TimeOut)
	}
	rid := ctx.Value(RequestIDHeader)
	if rid == nil {
		rid = uuid.New().String()
	}
	rctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	
	msg.Header.Set(DeadlineHeader, strconv.FormatInt(d.UnixMicro(), 10))
	msg.Header.Set(RequestIDHeader, rid.(string))
	reply, err := c.nc.RequestMsgWithContext(rctx, msg)
	return reply, err
	
}


