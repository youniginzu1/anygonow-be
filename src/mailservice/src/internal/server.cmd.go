package main

import (
	"os"
	"time"

	cli2 "github.com/aqaurius6666/go-utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	serviceName = "mail-service"
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
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run server",
				Action:  runMain,
			},
		},
		Flags: makeFlags(cli2.GormFlag, cli2.CommonServerFlag, cli2.LoggerFlag, CustomFlag),
	}
	return app

}

var (
	CustomFlag = []cli.Flag{
		&cli.StringFlag{
			Name:    "mail-username",
			EnvVars: []string{"CONFIG_MAIL_USERNAME"},
			Usage:   "Mail username",
		},
		&cli.StringFlag{
			Name:    "mail-password",
			EnvVars: []string{"CONFIG_MAIL_PASSWORD"},
			Usage:   "Mail password",
		},
		&cli.StringFlag{
			Name:    "smtp-host",
			EnvVars: []string{"CONFIG_MAIL_SMTP_HOST"},
			Usage:   "SMTP Host",
		},
		&cli.StringFlag{
			Name:    "smtp-port",
			EnvVars: []string{"CONFIG_MAIL_SMTP_PORT"},
			Usage:   "SMTP Port",
		},
		&cli.StringFlag{
			Name:    "smtp-sender",
			EnvVars: []string{"CONFIG_MAIL_SMTP_SENDER"},
			Usage:   "Mail sender",
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
			Name:     "firebase-private-key",
			EnvVars:  []string{"CONFIG_FIREBASE_PRIVATE_KEY"},
			Required: true,
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
