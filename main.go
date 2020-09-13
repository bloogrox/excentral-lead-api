package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/robfig/cron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/adapter/advertiser"
	conversionrepo "gitlab.com/cpanova/excentral/adapter/conversion"
	"gitlab.com/cpanova/excentral/adapter/fakeadvertiser"
	leadrepo "gitlab.com/cpanova/excentral/adapter/lead"
	postbackrepo "gitlab.com/cpanova/excentral/adapter/postback"

	advapi "gitlab.com/cpanova/excentral/ext/excentral"

	conversionLoadWrk "gitlab.com/cpanova/excentral/worker/conversion/load"
	conversionPersistWrk "gitlab.com/cpanova/excentral/worker/conversion/persist"
	conversionPostbackWrk "gitlab.com/cpanova/excentral/worker/conversion/postback"
	conversionProcessWrk "gitlab.com/cpanova/excentral/worker/conversion/process"

	leadHandler "gitlab.com/cpanova/excentral/delivery/rest/lead"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	leadRepo := leadrepo.New(db)
	conversionRepo := conversionrepo.New(db)
	postbackRepo := postbackrepo.New(db)
	_ = fakeadvertiser.New()

	exAffIDStr := os.Getenv("EXCENTRAL_AFF_ID")
	exPID := os.Getenv("EXCENTRAL_PID")
	exAffID, err := strconv.Atoi(exAffIDStr)
	if err != nil {
		panic(err)
	}
	advAPI := advapi.New(exAffID, exPID)

	advService := advertiser.New(advAPI)

	// PubSub router
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		middleware.CorrelationID,
		// middleware.Retry{
		// 	MaxRetries:      3,
		// 	InitialInterval: time.Second * 1,
		// 	Logger:          logger,
		// }.Middleware,
		// wMiddleware.Recoverer,
		// middleware.RandomFail(0.3),
		// middleware.RandomPanic(0.3),
	)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router.AddHandler(
		"load-conversions",
		"load-conversions",
		pubSub,
		"conversions",
		pubSub,
		conversionLoadWrk.NewHandler(advAPI).Handler,
	)

	router.AddHandler(
		"process-conversion",
		"conversions",
		pubSub,
		"lead.converted",
		pubSub,
		conversionProcessWrk.NewHandler(conversionRepo).Handler,
	)

	router.AddNoPublisherHandler(
		"persist-conversion",
		"lead.converted",
		pubSub,
		conversionPersistWrk.NewHandler(conversionRepo).Handler,
	)

	router.AddNoPublisherHandler(
		"postback-conversion",
		"lead.converted",
		pubSub,
		conversionPostbackWrk.NewHandler(leadRepo, postbackRepo).Handler,
	)

	go func() {
		ctx := context.Background()
		if err := router.Run(ctx); err != nil {
			panic(err)
		}
	}()

	// Cron
	crn := cron.New()

	crn.AddFunc("@hourly", func() {
		fmt.Println("Start syncing conversions")
		pubSub.Publish("load-conversions", message.NewMessage(watermill.NewUUID(), message.Payload([]byte(""))))
	})

	crn.Start()

	// HTTP server
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post(
		"/leads",
		leadHandler.NewHandler(
			advService,
			leadRepo,
		).Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :80")
		errs <- http.ListenAndServe(":80", r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
