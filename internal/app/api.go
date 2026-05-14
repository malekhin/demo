package app

import (
	"context"
	_ "demo/docs"
	"demo/internal/common/auth"
	"demo/internal/common/config"
	"demo/internal/common/metrics"
	"demo/internal/common/wrapper/database"
	skAdmin "demo/internal/domain/sk/http/admin"
	skPublic "demo/internal/domain/sk/http/public"
	"demo/internal/validator"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RouterParams struct {
	fx.In

	ErrorMiddleware             ErrorMiddleware
	PublicAgentAuthMiddleware   auth.PublicAgentAuthMiddleware
	AdminOperatorAuthMiddleware auth.AdminOperatorAuthMiddleware
	Config                      *config.Config
	Logger                      *zap.Logger
	Db                          database.Database
	HTTPMetricsCollector        metrics.HTTPMiddleware
	MetricsHandler              MetricsHandler

	/*
		Правило формирования алиаса:
		transdekra/http/api -> transdekraApi

		Правило формирования полей структуры:
		Правило для алиса + названия хендлера
		transdekraApi + TransdekraHandlers -> TransdekraApiTransdekra
	*/

	SkAdminSk  *skAdmin.SkHandlers
	SkPublicSk *skPublic.SkHandlers
}

func NewRouter(p RouterParams) *gin.Engine {
	switch config.GetEnvironment() {
	case config.EnvironmentTest:
		gin.SetMode(gin.TestMode)
	case config.EnvironmentLocal, config.EnvironmentDev:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	RegisterValidatorTagNameFunc()
	validator.RegisterCustomValidators(p.Db, p.Logger)

	router := gin.New()

	router.Use(ginzap.Ginzap(p.Logger.Named("GIN"), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(p.Logger.Named("GIN"), true))

	router.Use(gin.HandlerFunc(p.HTTPMetricsCollector))
	router.Use(gin.HandlerFunc(p.ErrorMiddleware))

	router.GET("/metrics", gin.HandlerFunc(p.MetricsHandler))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRoutes(router, p)
	apiRoutes(router, p)
	adminRoutes(router, p)

	return router
}

func publicRoutes(router gin.IRouter, p RouterParams) {
	v1 := router.Group("/public/v1")
	v1.Use(gin.HandlerFunc(p.PublicAgentAuthMiddleware))

	sk := v1.Group("/sk")
	{
		sk.GET("", p.SkPublicSk.SkList)
	}
}

func apiRoutes(router gin.IRouter, p RouterParams) {
	_ = router.Group("/api/v1")
	{
	}
}

func adminRoutes(router gin.IRouter, p RouterParams) {
	v1 := router.Group("/admin/v1")
	v1.
		Use(gin.HandlerFunc(p.AdminOperatorAuthMiddleware))
	{
		sk := v1.Group("/sk")
		{
			sk.GET("", p.SkAdminSk.SkList)
			sk.POST("", p.SkAdminSk.SkAdd)
			sk.POST("/:id", p.SkAdminSk.SkEdit)
		}
	}
}

func NewHttpServer(lc fx.Lifecycle, router *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:              ":8000",
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				_ = srv.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

var HttpServerModule = fx.Module("httpserver",
	fx.Provide(NewRouter),
	fx.Provide(NewHttpServer),
)
