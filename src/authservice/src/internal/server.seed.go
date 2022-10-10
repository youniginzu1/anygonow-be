package main

import (
	"context"
	"sync"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/db"
	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/authservice/src/internal/db/seed"
	"github.com/aqaurius6666/authservice/src/services/mailservice"
	"github.com/urfave/cli/v2"
)

func seedData(appCtx *cli.Context) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancelFn()
	mainServer, err := InitMainServer(ctx, logger, ServerOptions{
		DBDsn:    db.DBDsn(appCtx.String("db-uri")),
		MailAddr: mailservice.MailServiceAddr(appCtx.String("mailservice-address")),
	})
	defer func(db db.ServerRepo) {
		e := db.Close()
		if e != nil {
			panic("cannot close DB")
		}
	}(mainServer.MainRepo)
	if appCtx.Bool("clean") {
		logger.Debugf("start cleaning DB")
		err := clean(appCtx)
		if err != nil {
			logger.Errorf("failed cleaning DB")
			return err
		}
		logger.Debugf("sucessed cleaning DB")
	}
	err = mainServer.MainRepo.Migrate()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, sql := range seed.SQL {
		wg.Add(1)
		go func(sql string) {
			defer wg.Done()
			err = mainServer.MainRepo.RawSQL(sql)
		}(sql)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func clean(appCtx *cli.Context) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancelFn()
	mainServer, err := InitMainServer(ctx, logger, ServerOptions{
		DBDsn:    db.DBDsn(appCtx.String("db-uri")),
		MailAddr: mailservice.MailServiceAddr(appCtx.String("mailservice-address")),
	})
	defer func(db db.ServerRepo) {
		e := db.Close()
		if e != nil {
			panic("cannot close DB")
		}
	}(mainServer.MainRepo)
	if err != nil {
		return err
	}
	return mainServer.MainRepo.Drop()
}

func seedRole(db db.ServerRepo) error {
	return db.RawSQL(role.SQL)
}
