package main

import (
	"os"
	"time"

	cli2 "github.com/aqaurius6666/go-utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	serviceName = "main-service"
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
		Name:                 serviceName,
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
		Flags:  makeFlags(cli2.GormFlag, cli2.PrometheusFlag, cli2.CommonServerFlag, cli2.FeatureToggleFlag, cli2.LoggerFlag, CustomFlag),
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run server",
				Action:  runMain,
				Flags:   makeFlags(cli2.GormFlag, cli2.PrometheusFlag, cli2.CommonServerFlag, cli2.FeatureToggleFlag, cli2.LoggerFlag, CustomFlag),
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
			Name:     "authservice-address",
			EnvVars:  []string{"CONFIG_AUTHSERVICE_ADDRESS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-bucket-name",
			EnvVars:  []string{"CONFIG_AWS_BUCKET_NAME", "AWS_BUCKET_NAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-access-key-id",
			EnvVars:  []string{"CONFIG_AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-access-secret-key",
			EnvVars:  []string{"CONFIG_AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-region",
			EnvVars:  []string{"CONFIG_AWS_REGION", "AWS_REGION"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "otel-address",
			EnvVars: []string{"CONFIG_OTEL_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "chatservice-address",
			EnvVars: []string{"CONFIG_CHATSERVICE_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "mailservice-address",
			EnvVars: []string{"CONFIG_MAILSERVICE_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "stripe-key",
			EnvVars: []string{"CONFIG_STRIPE_KEY"},
		},
		&cli.StringFlag{
			Name:    "stripe-signature-verification",
			EnvVars: []string{"CONFIG_STRIPE_SIGNATURE_VERIFICATION"},
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
