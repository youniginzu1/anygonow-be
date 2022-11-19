package main

import (
	"os"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	cli2 "github.com/aqaurius6666/go-utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var logger *logrus.Logger

func main() {
	logger = logrus.New()
	if err := makeApp().Run(os.Args); err != nil {
		logger.WithField("err", err).Error("shutting down due to error")
		_ = os.Stderr.Sync()
		os.Exit(1)
	}
}

func makeApp() *cli.App {
	app := &cli.App{
		Name:                 c.SERVICE_NAME,
		Version:              "v1.0.1",
		EnableBashCompletion: true,
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Vu Nguyen",
				Email: "aqaurius6666@gmail.com",
			},
		},
		Action: runMain,
		Flags:  makeFlags(cli2.GormFlag, cli2.CommonServerFlag, cli2.LoggerFlag, CustomFlag, cli2.FeatureToggleFlag, cli2.PrometheusFlag, RedisFlag),
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run server",
				Action:  runMain,
				Flags:   makeFlags(cli2.GormFlag, cli2.CommonServerFlag, cli2.LoggerFlag, CustomFlag, cli2.FeatureToggleFlag, cli2.PrometheusFlag, RedisFlag),
			},
		},
	}
	return app

}

var (
	RedisFlag = []cli.Flag{
		&cli.StringFlag{
			Name:     "redis-uri",
			EnvVars:  []string{"CONFIG_REDIS_URI"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "redis-pass",
			EnvVars: []string{"CONFIG_REDIS_PASS"},
		},
		&cli.StringFlag{
			Name:    "redis-user",
			EnvVars: []string{"CONFIG_REDIS_USER"},
		},
	}
	CustomFlag = []cli.Flag{
		&cli.StringFlag{
			Name:     "mailservice-address",
			EnvVars:  []string{"CONFIG_MAILSERVICE_ADDRESS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "chatservice-address",
			EnvVars:  []string{"CONFIG_CHATSERVICE_ADDRESS"},
			Required: true,
		},
		// &cli.StringFlag{
		// 	Name:     "frontend-url",
		// 	EnvVars:  []string{"CONFIG_FRONTEND_URL"},
		// 	Required: true,
		// },
		&cli.StringFlag{
			Name:    "otel-address",
			EnvVars: []string{"CONFIG_OTEL_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "stripe-key",
			EnvVars: []string{"CONFIG_STRIPE_KEY"},
		},
		&cli.IntFlag{
			Name:    "quantity-interval",
			EnvVars: []string{"CONFIG_QUANTITY_INTERVAL"},
			Value:   1,
		},
		&cli.StringFlag{
			Name:    "unit-interval",
			EnvVars: []string{"CONFIG_UNIT_INTERVAL"},
			Value:   "WEEK",
		},
		&cli.StringFlag{
			Name:    "payment_day",
			EnvVars: []string{"CONFIG_PAYMENT_DAY"},
			Value:   "Sunday",
		},
		&cli.IntFlag{
			Name:    "chat-quantity-interval",
			EnvVars: []string{"CONFIG_CHAT_QUANTITY_INTERVAL"},
			Value:   1,
		},
	}
)

func makeFlags(lists ...interface{}) []cli.Flag {
	flags := make([]cli.Flag, 0)
	for _, f := range lists {
		tmp, _ := f.([]cli.Flag)
		flags = append(flags, tmp...)
	}
	return flags
}
