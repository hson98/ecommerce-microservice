package server

import (
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	"github.com/hson98/ecommerce-microservice/src/auth-service/internal/interceptors"
	"github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/delivery/grpc/service"
	user_repository "github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/repository"
	user_usecase "github.com/hson98/ecommerce-microservice/src/auth-service/internal/user/usecase"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/logger"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/metric"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/myjwt"
	userService "github.com/hson98/ecommerce-microservice/src/auth-service/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	db       *gorm.DB
	jwtMaker myjwt.Maker
	config   *config.Config
	logger   logger.Logger
}

func NewAuthServer(db *gorm.DB, jwtMaker myjwt.Maker, config *config.Config, logger logger.Logger) *Server {
	return &Server{db: db, jwtMaker: jwtMaker, config: config, logger: logger}
}
func (s *Server) Run() error {
	//setup metrics
	metrics, err := metric.CreateMetrics(s.config.Metrics.URL, s.config.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.config.Metrics.URL,
		s.config.Metrics.ServiceName,
	)
	im := interceptors.NewInterceptorManager(s.logger, s.config, metrics)
	userRepo := user_repository.NewUserPgRepo(s.db)
	userUC := user_usecase.NewUserUC(userRepo, s.config, s.jwtMaker)

	l, err := net.Listen("tcp", s.config.Server.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.config.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.config.Server.Timeout * time.Second,
		MaxConnectionAge:  s.config.Server.MaxConnectionAge * time.Minute,
		Time:              s.config.Server.Timeout * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	if s.config.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGRPCServer := service.NewAuthServerGRPC(s.logger, s.config, userUC, s.jwtMaker)
	userService.RegisterUserServiceServer(server, authGRPCServer)

	grpc_prometheus.Register(server)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		s.logger.Infof("Server is listening on port: %v", s.config.Server.Port)
		if err := server.Serve(l); err != nil {
			s.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
