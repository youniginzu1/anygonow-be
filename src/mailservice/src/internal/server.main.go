package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/stackdriver"
	commongrpc "github.com/aqaurius6666/go-utils/common_grpc"
	commonpb "github.com/aqaurius6666/go-utils/common_grpc/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/aqaurius6666/mailservice/src/internal/db"
	"github.com/aqaurius6666/mailservice/src/internal/lib"
	"github.com/aqaurius6666/mailservice/src/internal/lib/unleash"
	"github.com/aqaurius6666/mailservice/src/internal/mail"
	"github.com/aqaurius6666/mailservice/src/pb/mailpb"
	"github.com/aqaurius6666/mailservice/src/services/fcm"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/reflection"
)

func runMain(appCtx *cli.Context) error {
	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	errChan := make(chan error, 1)
	if appCtx.Bool("disable-tracing") && false {
		logger.Info("Tracing disabled.")
	} else {
		logger.Info("Tracing enabled.")
		go initTracing()
	}
	if appCtx.Bool("disable-profiler") {
		logger.Info("Profiling disabled.")
	} else {
		logger.Info("Profiling enabled.")
		go initProfiling(serviceName, appCtx.String("runtime-version"))
	}
	if appCtx.Bool("feature-toggle") {
		err := unleash.RegisterFeatureToggle(appCtx.String("unleash-app-name"), appCtx.String("unleash-token"), appCtx.String("unleash-api-url"))
		if err != nil {
			return err
		}
	}
	// Start HTTP Server
	httpListner, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("http-port")))
	if err != nil {
		logger.Fatal(err)
		return err
	}
	mainServer, err := InitMainServer(ctx, logger, ServerOptions{
		Username:    mail.MailUsername(appCtx.String("mail-username")),
		Password:    mail.MailPassword(appCtx.String("mail-password")),
		Host:        mail.SMTPHost(appCtx.String("smtp-host")),
		Port:        mail.SMTPPort(appCtx.String("smtp-port")),
		Sender:      mail.SMTPSender(appCtx.String("smtp-sender")),
		DBDsn:       db.DBDsn(appCtx.String("db-uri")),
		FirebaseKey: fcm.FB_PRIVATE_KEY(appCtx.String("firebase-private-key")),
	})
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer func() {
		_ = httpListner.Close()
	}()
	if err != nil {
		logger.Fatal(err)
		return err
	}
	if err := mainServer.Repo.Migrate(); err != nil {
		logger.Fatal(err)
		return err
	}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	mainServer.ApiServer.RegisterEndpoint()
	// 	logger.WithField("port", appCtx.Int("http-port")).Info("listening for HTTP connections")
	// 	if err := http.Serve(httpListner, mainServer.ApiServer.G); err != nil {
	// 		logger.Fatalf("failed to serve: %v", err)
	// 	}
	// }()

	// Start GRPC Server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("grpc-port")))
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer func() { _ = grpcListener.Close() }()
	var srv *grpc.Server
	commonServer := commongrpc.NewCommonServer(logger, appCtx.Bool("allow-kill"))
	wg.Add(1)
	go func() {
		defer wg.Done()
		if appCtx.Bool("disable-stats") {
			logger.Info("Stats disabled.")
			srv = grpc.NewServer(grpc.ChainUnaryInterceptor(
				otelgrpc.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(MyLogger)),
				lib.UnaryServerLogRequestInterceptor,
			))
		} else {
			logger.Info("Stats enabled.")
			srv = grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
		}
		healthpb.RegisterHealthServer(srv, commonServer)
		commonpb.RegisterCommonServer(srv, commonServer)
		mailpb.RegisterMailServiceServer(srv, mainServer.ApiServer)
		reflection.Register(srv)
		logger.WithField("port", appCtx.Int("grpc-port")).Info("listening for gRPC connections")
		if err := srv.Serve(grpcListener); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start pprof server
	pprofListener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("pprof-port")))
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer func() {
		_ = pprofListener.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.WithField("port", appCtx.Int("pprof-port")).Info("listening for pprof requests")
		sSrv := new(http.Server)
		_ = sSrv.Serve(pprofListener)
	}()

	if !appCtx.Bool("disable-prometheus") {
		// Start Prometheus Server
		promListener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("prometheus-port")))
		if err != nil {
			return err
		}
		defer func() {
			_ = promListener.Close()
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			promServer := http.NewServeMux()
			promServer.Handle("/metrics", promhttp.Handler())
			logger.WithField("port", appCtx.Int("prometheus-port")).Info("listening for metrics requests")
			if err := http.Serve(promListener, promServer); err != nil {
				errChan <- err
			}
		}()
	} else {
		logger.Info("Prometheus disabled.")
	}

	// Watch kill signal
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
		select {
		case s := <-sigCh:
			cancelFn()
			// Handle graceful shutdown here

			logger.WithField("signal", s.String()).Infof("shutting down due to signal")
		case <-ctx.Done():

		case err := <-errChan:
			cancelFn()
			logger.WithField("error", err.Error()).Errorf("shutting down due to error")
		}
	}()
	wg.Wait()
	return nil
}

func initTracing() {
	initJaegerTracing()
	// initStackdriverTracing()
}

func initJaegerTracing() {
	// addr := os.Getenv("CONFIG_JAGGER_ADDRESS")
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: fmt.Sprintf("http://%s", "jeager:14268"),
		Process: jaeger.Process{
			ServiceName: "checkoutservice",
		},
		ServiceName: "test",
	})
	if err != nil {
		logger.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	logger.Info("jaeger initialization completed.")
}
func initStats(exporter *stackdriver.Exporter) {
	view.SetReportingPeriod(60 * time.Second)
	view.RegisterExporter(exporter)
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		logger.Warn("Error registering default server views")
	} else {
		logger.Info("Registered default server views")
	}
}

func initStackdriverTracing() {
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			logger.Infof("failed to initialize stackdriver exporter: %+v", err)
		} else {
			trace.RegisterExporter(exporter)
			logger.Info("registered Stackdriver tracing")

			// Register the views to collect server stats.
			initStats(exporter)
			return
		}
		d := time.Second * 10 * time.Duration(i)
		logger.Infof("sleeping %v to retry initializing Stackdriver exporter", d)
		time.Sleep(d)
	}
	logger.Warn("could not initialize Stackdriver exporter after retrying, giving up")
}

func initProfiling(service, version string) {
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: version,
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",
		}); err != nil {
			logger.Warnf("failed to start profiler: %+v", err)
		} else {
			logger.Info("started Stackdriver profiler")
			return
		}
		d := time.Second * 10 * time.Duration(i)
		logger.Infof("sleeping %v to retry initializing Stackdriver profiler", d)
		time.Sleep(d)
	}
	logger.Warn("could not initialize Stackdriver profiler after retrying, giving up")
}

func MyLogger(p interface{}) error {
	var err error
	switch t := p.(type) {
	case error:
		err = t
		err = utils.Unwrap(err)
	case string:
		err = errors.New(t)
	default:
		err = errors.New(fmt.Sprintf("%+v", t))
	}
	logger.Errorf("%+v", err)
	return status.Error(codes.Internal, err.Error())
}
