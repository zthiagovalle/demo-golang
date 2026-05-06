package main

import (
	"context"
	"log"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/zthiagovalle/demo-golang/src/application/consumers"
	"github.com/zthiagovalle/demo-golang/src/application/controllers"
	"github.com/zthiagovalle/demo-golang/src/core/config"
	"github.com/zthiagovalle/demo-golang/src/core/shared"
	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/usecases"
	"github.com/zthiagovalle/demo-golang/src/infra/gateways"
	"github.com/zthiagovalle/demo-golang/src/infra/producers"
	"github.com/zthiagovalle/demo-golang/src/infra/repositories"
)

const gracefulTimeout = 10 * time.Second

type app struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	snsClient *sns.Client
	sqsClient *sqs.Client
	echo      *echo.Echo
}

func newApp(ctx context.Context, cfg *config.Config) *app {
	pool := mustOpenDatabase(ctx, cfg)
	snsClient, sqsClient := mustBuildAWSClients(ctx, cfg)
	return &app{
		cfg:       cfg,
		pool:      pool,
		snsClient: snsClient,
		sqsClient: sqsClient,
		echo:      buildEchoServer(),
	}
}

func (a *app) close() {
	a.pool.Close()
}

func (a *app) setupConsumers(ctx context.Context) {
	catalogGateway := gateways.NewCatalogGateway(a.cfg.CatalogGatewayBaseURL)
	productConsumer := consumers.NewProductConsumer(a.sqsClient, a.cfg.SQSProductCreatedQueueURL, catalogGateway)
	go productConsumer.Run(ctx)
}

func (a *app) setupControllers() {
	productsRepository := repositories.NewProductsPostgresRepository(a.pool)
	productProducer := producers.NewProductSNSProducer(a.snsClient, a.cfg.SNSProductTopicARN)

	productsController := controllers.NewProductsV1Controller(
		usecases.NewGetAllPaginatedProductsUsecase(productsRepository),
		usecases.NewCreateProductUsecase(productsRepository, productProducer),
		usecases.NewUpdateProductUsecase(productsRepository),
		usecases.NewDeleteProductUsecase(productsRepository),
		usecases.NewToggleProductStatusUsecase(productsRepository),
	)

	shared.RegisterRoutes(a.echo, productsController.Routes())
}

func (a *app) start(ctx context.Context) error {
	log.Printf("listening on :%s", a.cfg.AppPort)
	sc := echo.StartConfig{
		Address:         ":" + a.cfg.AppPort,
		GracefulTimeout: gracefulTimeout,
	}
	return sc.Start(ctx, a.echo)
}

func buildEchoServer() *echo.Echo {
	cv := shared.NewCustomValidator()
	registerCustomValidators(cv)

	e := echo.New()
	e.Validator = cv
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	return e
}

func registerCustomValidators(cv *shared.CustomValidator) {
	_ = cv.Validator.RegisterValidation("oneOfProductStatus", enums.ProductStatusValidator)
}

func mustOpenDatabase(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	return pool
}

func mustBuildAWSClients(ctx context.Context, cfg *config.Config) (*sns.Client, *sqs.Client) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKey, cfg.AWSSecretKey, "")),
	)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}

	snsClient := sns.NewFromConfig(awsCfg, func(o *sns.Options) {
		if cfg.AWSEndpointURL != "" {
			o.BaseEndpoint = awsv2.String(cfg.AWSEndpointURL)
		}
	})
	sqsClient := sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		if cfg.AWSEndpointURL != "" {
			o.BaseEndpoint = awsv2.String(cfg.AWSEndpointURL)
		}
	})
	return snsClient, sqsClient
}
