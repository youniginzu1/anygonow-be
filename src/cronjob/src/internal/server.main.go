package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aqaurius6666/cronjob/src/internal/cronjob"
	"github.com/aqaurius6666/cronjob/src/internal/db"
	"github.com/aqaurius6666/cronjob/src/internal/lib/unleash"
	"github.com/aqaurius6666/cronjob/src/internal/var/c"
	"github.com/aqaurius6666/cronjob/src/services/chatservice"
	"github.com/aqaurius6666/cronjob/src/services/mailservice"
	"github.com/aqaurius6666/cronjob/src/services/payment"
	"github.com/aqaurius6666/cronjob/src/services/redis"
	commongrpc "github.com/aqaurius6666/go-utils/common_grpc"
	commonpb "github.com/aqaurius6666/go-utils/common_grpc/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	tp *tracesdk.TracerProvider
)

func runMain(appCtx *cli.Context) error {
	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	errChan := make(chan error, 1)
	if appCtx.Bool("disable-tracing") {
		logger.Info("Tracing disabled.")
	} else {
		logger.Info("Tracing enabled.")
		go initOtelCollector(ctx, appCtx.String("otel-address"))
	}
	if appCtx.Bool("disable-profiler") {
		logger.Info("Profiling disabled.")
	} else {
		logger.Info("Profiling enabled.")
		// go initProfiling(c.SERVICE_NAME, appCtx.String("runtime-version"))
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
		DBDsn:                db.DBDsn(appCtx.String("db-uri")),
		Key:                  payment.STRIPE_API_KEY(appCtx.String("stripe-key")),
		QuantityInterval:     cronjob.QUANTITY_INTERVAL(appCtx.Int("quantity-interval")),
		UnitInterval:         cronjob.UNIT_INTERVAL(appCtx.String("unit-interval")),
		PaymentDay:           cronjob.PAYMENT_DAY(appCtx.String("payment_day")),
		MailserviceAddr:      mailservice.MailserviceAddr(appCtx.String("mailservice-address")),
		ChatQuantityInterval: cronjob.CHAT_QUANTITY_INTERVAL(appCtx.Int("chat-quantity-interval")),
		RedisUri:             redis.REDIS_URI(appCtx.String("redis-uri")),
		RedisUser:            redis.REDIS_USER(appCtx.String("redis-user")),
		RedisPass:            redis.REDIS_PASS(appCtx.String("redis-pass")),
		ChatserviceAddr:      chatservice.ChatserviceAddr(appCtx.String("chatservice-address")),
	})
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer func() {
		_ = httpListner.Close()
		_ = mainServer.MainRepo.Close()
	}()
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
			srv = grpc.NewServer(
				grpc.ChainUnaryInterceptor(
				// otelgrpc.UnaryServerInterceptor(),
				// grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(MyLogger)),
				))
		} else {
			logger.Info("Stats enabled.")
			// srv = grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
		}
		healthpb.RegisterHealthServer(srv, commonServer)
		commonpb.RegisterCommonServer(srv, commonServer)
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
	// Run job
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, job := range mainServer.Jobs {
			wg.Add(1)
			go func(job cronjob.Cronjob) {
				defer wg.Done()
				job.Run(ctx)
			}(job)
		}
	}()

	// Watch kill signal
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
		select {
		case s := <-sigCh:
			cancelFn()
			// Handle graceful shutdown here
			if tp != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				if err := tp.Shutdown(ctx); err != nil {
					panic(err)
				}
			}

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

func initOtelCollector(ctx context.Context, otelAgentAddr string) {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
	)
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		panic(err)
	}

	bsp := tracesdk.NewBatchSpanProcessor(traceExp)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.SERVICE_NAME),
			attribute.Int64("ID", c.ID),
		)),
		tracesdk.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)
}

// func initStats(exporter *stackdriver.Exporter) {
// 	view.SetReportingPeriod(60 * time.Second)
// 	view.RegisterExporter(exporter)
// 	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
// 		logger.Warn("Error registering default server views")
// 	} else {
// 		logger.Info("Registered default server views")
// 	}
// }

// func initStackdriverTracing() {
// 	// since they are not sharing packages.
// 	for i := 1; i <= 3; i++ {
// 		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
// 		if err != nil {
// 			logger.Infof("failed to initialize stackdriver exporter: %+v", err)
// 		} else {
// 			trace.RegisterExporter(exporter)
// 			logger.Info("registered Stackdriver tracing")

// 			// Register the views to collect server stats.
// 			initStats(exporter)
// 			return
// 		}
// 		d := time.Second * 10 * time.Duration(i)
// 		logger.Infof("sleeping %v to retry initializing Stackdriver exporter", d)
// 		time.Sleep(d)
// 	}
// 	logger.Warn("could not initialize Stackdriver exporter after retrying, giving up")
// }

// func initProfiling(service, version string) {
// 	// since they are not sharing packages.
// 	for i := 1; i <= 3; i++ {
// 		if err := profiler.Start(profiler.Config{
// 			Service:        service,
// 			ServiceVersion: version,
// 			// ProjectID must be set if not running on GCP.
// 			// ProjectID: "my-project",
// 		}); err != nil {
// 			logger.Warnf("failed to start profiler: %+v", err)
// 		} else {
// 			logger.Info("started Stackdriver profiler")
// 			return
// 		}
// 		d := time.Second * 10 * time.Duration(i)
// 		logger.Infof("sleeping %v to retry initializing Stackdriver profiler", d)
// 		time.Sleep(d)
// 	}
// 	logger.Warn("could not initialize Stackdriver profiler after retrying, giving up")
// }

func MyLogger(p interface{}) error {
	err, _ := p.(error)
	logger.Errorf("%+v", err)
	err = utils.Unwrap(err)
	return status.Error(codes.Internal, err.Error())
}
