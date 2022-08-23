package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nepomuceno/cloud-ipam/ipconfig"
)

type ApiClient struct {
	ipClient *ipconfig.IpConfigClient
}

func NewApiClient(ipClient *ipconfig.IpConfigClient) *ApiClient {
	return &ApiClient{ipClient: ipClient}
}

func (client *ApiClient) RegisterRoutes(r *gin.Engine) error {
	api := r.Group("/api")
	client.registerEnvironmentRoutes(api)
	client.registerSubscriptionRoutes(api)
	return nil
}

func (client *ApiClient) registerSubscriptionRoutes(r *gin.RouterGroup) {
	subscription := r.Group("/environment/:environmentID/subscription")
	subscription.GET("/:subscriptionID", func(ctx *gin.Context) {
		environmentID := ctx.Param("environmentID")
		subscriptionID := ctx.Param("subscriptionID")
		subscription, err := client.ipClient.GetSubscription(environmentID, subscriptionID)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		ctx.JSON(200, subscription)
	})
	subscription.GET("/", func(ctx *gin.Context) {
		environmentID := ctx.Param("environmentID")
		subscriptions, err := client.ipClient.ListSubscriptions(environmentID)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		ctx.JSON(200, subscriptions)
	})
}

func (client *ApiClient) registerEnvironmentRoutes(r *gin.RouterGroup) {
	envRoute := r.Group("/environment")
	envRoute.GET("/:environmentID", func(c *gin.Context) {
		environmentID := c.Param("environment")
		environment, err := client.ipClient.GetEnvironment(environmentID)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, environment)
	})
	envRoute.GET("/", func(c *gin.Context) {
		environment, err := client.ipClient.ListEnvironments()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, environment)
	})

}
