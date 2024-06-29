package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-container-api/config"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("api")

func ServerAPI(cfg *config.Config) {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Use(Cors())
	r.Use(RequestLoggerMiddleware())

	apiV1 := r.Group("/api/v1")
	authMiddleware, err := jwtGinMiddleware(cfg.SecretKey)
	if err != nil {
		log.Fatalf("jwt auth middleware: %v", err)
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		log.Fatalf("authMiddleware.MiddlewareInit: %v", err)
	}

	user := apiV1.Group("/user")
	user.POST("/login", authMiddleware.LoginHandler)
	user.POST("/logout", authMiddleware.LogoutHandler)
	user.GET("/refresh_token", authMiddleware.RefreshHandler)

	user.Use(authMiddleware.MiddlewareFunc())
	user.POST("/info", GetUserInfoHandler)

	container := apiV1.Group("/container")
	container.Use(authMiddleware.MiddlewareFunc())
	container.GET("/providers", GetProvidersHandler)
	container.GET("/deployments", GetDeploymentsHandler)
	container.GET("/deployment/manifest", GetDeploymentManifestHandler)
	container.POST("/deployment/create", CreateDeploymentHandler)
	container.POST("/deployment/delete", DeleteDeploymentHandler)
	container.POST("/deployment/update", UpdateDeploymentHandler)
	container.GET("/deployment/logs", GetDeploymentLogsHandler)
	container.GET("/deployment/event", GetDeploymentEventsHandler)
	container.GET("/deployment/domains", GetDeploymentDomainHandler)
	container.POST("/deployment/domain/add", AddDeploymentDomainHandler)
	container.POST("/deployment/domain/del", DeleteDeploymentDomainHandler)
	container.GET("/deployment/shell", GetDeploymentShellHandler)

	if err := r.Run(cfg.ApiListen); err != nil {
		log.Fatalf("starting server: %v\n", err)
	}
}
