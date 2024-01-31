package main

import (
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	"github.com/hson98/ecommerce-microservice/src/auth-service/internal/server"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/constants"
	jaegerTracer "github.com/hson98/ecommerce-microservice/src/auth-service/pkg/jaeger"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/logger"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/myjwt"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/postgres"
	"github.com/opentracing/opentracing-go"
	"log"
	"os"
)

func main() {
	config, err := config.LoadConfig(os.Getenv(constants.ConfigPath))
	if err != nil {
		log.Fatalf("Loading config: %v", err)

	}

	psqlDB := postgres.NewPostgresDB(config)
	appLogger := logger.NewAPILogger(config)

	//khởi tạo JWT
	jwtMaker, err := myjwt.NewJwtMaker(config.Server.SecretKeyJWT)
	if err != nil {
		panic(err)
	}

	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		config.Server.AppVersion,
		config.Logger.Level,
		config.Server.Mode,
		config.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", config.Server.AppVersion)

	tracer, closer, err := jaegerTracer.InitJaeger(config)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	authServer := server.NewAuthServer(psqlDB, jwtMaker, config, appLogger)
	appLogger.Fatal(authServer.Run())
}
