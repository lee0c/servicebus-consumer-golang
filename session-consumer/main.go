package main

import (
	"context"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	log "github.com/sirupsen/logrus"
)


type TimeoutSessionHandler struct {
	ms *servicebus.MessageSession
	timer *time.Timer
	timeout time.Duration
}

func (sh *TimeoutSessionHandler) Start(ms *servicebus.MessageSession) error {
	sh.ms = ms
	sh.timeout = time.Second * 30
	log.Infof("Accepted session lock\n")
	return nil
}

func (sh *TimeoutSessionHandler) Handle(ctx context.Context, msg *servicebus.Message) error {
	if sh.timer != nil {
		sh.timer.Stop()
	}

	log.Infof("Received message:\n\t[Session %s; Id %s] %s\n", *msg.SessionID, msg.ID, string(msg.Data))
	sh.timer = time.NewTimer(sh.timeout)

	go func() {
		<-sh.timer.C
		sh.ms.Close()
	}()

	return msg.Complete(ctx)
}

func (sh *TimeoutSessionHandler) End() {
	log.Infof("Released session lock\n")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName := os.Getenv("QUEUE_NAME")

	if connStr == "" {
		log.Errorf("Connection string unset")
		os.Exit(1)
	}
	if queueName == "" {
		queueName = "sessionqueue"
	}

	// Set up sb entities
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
    if err != nil {
        log.Errorf("Failed to get namespace %s\n", err)
        return
    }

    q, err := ns.NewQueue(queueName)
    if err != nil {
        log.Errorf("Failed to get queue %s\n", err)
        return
    }
	defer q.Close(ctx)

	// Set up session receiver and handler
	for {
		queueSession := q.NewSession(nil)
		err = queueSession.ReceiveOne(ctx, new(TimeoutSessionHandler))
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			queueSession.Close(ctx)
		}
	}
}