package main

import (
	"context"
	"flag"
	"io/ioutil"
	stdHttp "net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	kafkaAddr = flag.String("kafka", "localhost:9092", "The address of the kafka broker")
	httpAddr  = flag.String("http", ":8080", "The address for the http subscriber")
)

func main() {
	flag.Parse()
	logger := watermill.NewStdLogger(true, true).With(
		watermill.LogFields{},
	)

	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{*kafkaAddr},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w stdHttp.ResponseWriter, r *stdHttp.Request) {
		w.Write([]byte("kukusiki"))
	})

	httpSubscriber, err := http.NewSubscriber(
		*httpAddr,
		http.SubscriberConfig{
			Router: router,
			UnmarshalMessageFunc: func(topic string, request *stdHttp.Request) (*message.Message, error) {
				body, _ := ioutil.ReadAll(request.Body)

				logger.Info("UnmarshalMessageFunc Invoked", nil)

				return message.NewMessage(
					watermill.NewUUID(),
					body,
				), nil
			},
		},
		logger,
	)
	if err != nil {
		logger.Error("Could not create HTTP subscriber", err, nil)
	}

	r, err := message.NewRouter(
		message.RouterConfig{},
		logger,
	)
	if err != nil {
		panic(err)
	}

	r.AddHandler(
		"http_to_kafka", // name for debug
		"/webhooks",     // api url
		httpSubscriber,  // subscriber
		"webhooks",      // name of topic
		kafkaPublisher,  // publisher
		func(msg *message.Message) ([]*message.Message, error) {
			logger.Info("HandleFunc invoked", nil)
			return []*message.Message{msg}, nil
		},
	)
	// Run router asynchronously
	go r.Run(context.Background())
	// Check if router is running then start HTTP server
	<-r.Running()

	logger.Info("Starting HTTP server", nil)
	err = httpSubscriber.StartHTTPServer()
	if err != nil {
		logger.Error("Could not start HTTP server", err, nil)
	}
}
