package broker

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"h12.me/kpax/model"
)

var (
	ErrCorrelationIDMismatch = errors.New("correlationID mismatch")
	ErrBrokerClosed          = errors.New("broker is closed")
)

type B struct {
	Addr     string
	Timeout  time.Duration
	QueueLen int
	cid      int32
	conn     *connection
	mu       sync.Mutex
}

type brokerJob struct {
	req     model.Request
	resp    model.Response
	errChan chan error
}

func New(addr string) model.Broker {
	return &B{
		Addr:     addr,
		Timeout:  30 * time.Second,
		QueueLen: 1000,
	}
}

func (b *B) Do(req model.Request, resp model.Response) error {
	req.SetID(atomic.AddInt32(&b.cid, 1))
	errChan := make(chan error)
	if err := b.sendJob(&brokerJob{
		req:     req,
		resp:    resp,
		errChan: errChan,
	}); err != nil {
		return err
	}
	return <-errChan
}

func (b *B) sendJob(job *brokerJob) error {
	b.mu.Lock()
	if b.conn == nil || b.conn.Closed() {
		conn, err := b.newConn()
		if err != nil {
			b.mu.Unlock()
			return err
		}
		b.conn = conn
	}
	b.mu.Unlock()
	b.conn.sendChan <- job
	return nil
}

func (b *B) newConn() (*connection, error) {
	conn, err := net.Dial("tcp", b.Addr)
	if err != nil {
		return nil, err
	}
	c := &connection{
		conn:     conn,
		sendChan: make(chan *brokerJob),
		recvChan: make(chan *brokerJob, b.QueueLen),
		timeout:  b.Timeout,
	}
	go c.sendLoop()
	go c.receiveLoop()
	return c, nil
}

func (b *B) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.conn != nil && !b.conn.Closed() {
		b.conn.Close()
		b.conn = nil
	}
}
