package main

import (
	"os"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	cli2 "github.com/aqaurius6666/go-utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var logger *logrus.Logger

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
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
			{
				Name:   "seed",
				Usage:  "seed data",
				Action: seedData,
				Flags: makeFlags(cli2.GormFlag, cli2.LoggerFlag, []cli.Flag{
					&cli.BoolFlag{
						Name:  "clean",
						Usage: "clean before seed",
						Value: false,
					},
				}),
			},
			{
				Name:   "clean",
				Usage:  "clean database",
				Action: clean,
				Flags: makeFlags(cli2.GormFlag, cli2.LoggerFlag, []cli.Flag{
					&cli.BoolFlag{
						Name:  "clean",
						Usage: "clean before seed",
						Value: false,
					},
				}),
			},
		},
	}
	return app

}

var (
	CustomFlag = []cli.Flag{
		&cli.StringFlag{
			Name:     "twillo-callback-url",
			EnvVars:  []string{"CONFIG_TWILLO_CALLBACK_URL"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "mailservice-address",
			EnvVars:  []string{"CONFIG_MAILSERVICE_ADDRESS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "otel-address",
			EnvVars: []string{"CONFIG_OTEL_ADDRESS"},
		},
		&cli.StringFlag{
			Name:     "authservice-address",
			EnvVars:  []string{"CONFIG_AUTHSERVICE_ADDRESS"},
			Required: true,
		}, &cli.StringFlag{
			Name:     "twillo-account-sid",
			EnvVars:  []string{"TWILIO_ACCOUNT_SID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "twillo-auth-token",
			EnvVars:  []string{"TWILIO_AUTH_TOKEN"},
			Required: true,
		},
	}
)

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
)

func makeFlags(lists ...interface{}) []cli.Flag {
	flags := make([]cli.Flag, 0)
	for _, f := range lists {
		tmp, _ := f.([]cli.Flag)
		flags = append(flags, tmp...)
	}
	return flags
}
