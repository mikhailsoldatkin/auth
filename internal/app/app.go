package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"github.com/natefinch/lumberjack"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/interceptor"
	pbAccess "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	pbAuth "github.com/mikhailsoldatkin/auth/pkg/auth_v1"
	pbUser "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Register statik to serve Swagger UI and static files
	_ "github.com/mikhailsoldatkin/auth/statik"
)

// App represents the application with its dependencies and gRPC, HTTP and Swagger servers.
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

// NewApp initializes a new App instance with the given context and sets up the necessary dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the GRPC server, HTTP server, Swagger server, and Kafka consumer.
// It handles graceful shutdown by waiting for context cancellation or termination signals.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Fatalf("failed to run Kafka consumer: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			logger.Fatalf("failed to run Prometheus server: %v", err)
		}
	}()

	gracefulShutdown(ctx, cancel, wg)

	return nil
}

// gracefulShutdown handles the termination process by waiting for either a context cancellation
// or termination signals. It cancels the context and waits for all goroutines to finish.
func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		logger.Info("terminating: context cancelled")
	case <-waitSignal():
		logger.Info("terminating: via signal")
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

// waitSignal creates a channel to receive termination signals (SIGINT, SIGTERM).
// It returns the channel to allow waiting for these signals.
func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	return sigterm
}

// initDeps initializes the dependencies required by the App.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initLogger,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads the application configuration.
func (a *App) initConfig(_ context.Context) error {
	_, err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

// initServiceProvider initializes the service provider used by the application.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

// initLogger initializes the app logger.
func (a *App) initLogger(_ context.Context) error {
	var level zapcore.Level
	if err := level.Set(a.serviceProvider.config.Logger.Level); err != nil {
		return err
	}

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   a.serviceProvider.config.Logger.Filename,
		MaxSize:    a.serviceProvider.config.Logger.MaxSizeMB,
		MaxBackups: a.serviceProvider.config.Logger.MaxBackups,
		MaxAge:     a.serviceProvider.config.Logger.MaxAgeDays,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	logger.Init(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, zap.NewAtomicLevelAt(level)),
		zapcore.NewCore(fileEncoder, file, zap.NewAtomicLevelAt(level)),
	))

	return nil
}

// initGRPCServer initializes the gRPC server.
func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	if err != nil {
		logger.Fatal("failed to load TLS credentials from files:", zap.Error(err))
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.MetricsInterceptor,
				interceptor.ValidateInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)

	pbUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImplementation(ctx))
	pbAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImplementation(ctx))
	pbAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImplementation(ctx))

	return nil
}

// initHTTPServer initializes the HTTP server and sets up the GRPC gateway for HTTP-to-GRPC translation.
func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		logger.Fatal("failed to load TLS credentials for gRPC gateway:", zap.Error(err))
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = pbUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.config.GRPC.Address, opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.config.HTTP.Address,
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

// initSwaggerServer initializes the Swagger server to serve Swagger UI and API documentation.
func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.config.Swagger.Address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

// runHTTPServer starts the HTTP server and listens for incoming requests.
func (a *App) runHTTPServer() error {
	logger.Infof("HTTP server is running on %v", a.serviceProvider.config.HTTP.Port)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// runSwaggerServer starts the Swagger server to serve API documentation.
func (a *App) runSwaggerServer() error {
	logger.Infof("Swagger server is running on %d", a.serviceProvider.config.Swagger.Port)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// serveSwaggerFile returns an HTTP handler function to serve Swagger files.
func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Infof("serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			_ = file.Close()
		}(file)

		logger.Infof("read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("served swagger file: %s", path)
	}
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    a.serviceProvider.config.Prometheus.Address,
		Handler: mux,
	}

	return nil
}

// runPrometheusServer starts the Prometheus server.
func (a *App) runPrometheusServer() error {
	logger.Infof("Prometheus server is running on %d", a.serviceProvider.config.Prometheus.Port)

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// runGRPCServer starts the GRPC server and listens for incoming GRPC requests.
func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.config.GRPC.Address)
	if err != nil {
		return err
	}

	logger.Infof("gRPC server is running on %d", a.serviceProvider.config.GRPC.Port)

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
