package routes

import (
	"net/http"

	"securewallet/internal/config"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupHealthRoutes(r *gin.Engine) {
	r.GET("/health", handleLegacyHealth)
	r.GET("/health/live", handleLiveness)
	r.GET("/health/ready", handleReadiness)
	r.GET("/health/heartbeat", handleHeartbeat)
	r.GET("/health/system", handleSystemInfo)
	r.GET("/health/shards", handleShardHealth)
}

// @Summary Liveness probe
// @Description Returns whether the application process is alive. Used by Kubernetes liveness probes.
// @Tags health
// @Produce json
// @Success 200 {object} services.HealthStatus
// @Router /health/live [get]
func handleLiveness(c *gin.Context) {
	status := services.Liveness()
	c.JSON(http.StatusOK, status)
}

// @Summary Readiness probe
// @Description Returns whether the application is ready to serve traffic, including dependency health checks.
// @Tags health
// @Produce json
// @Success 200 {object} services.HealthStatus
// @Success 503 {object} services.HealthStatus
// @Router /health/ready [get]
func handleReadiness(c *gin.Context) {
	status := services.Readiness()
	if status.Status != "ok" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	c.JSON(http.StatusOK, status)
}

// @Summary Workflow heartbeat
// @Description Returns workflow health including database, redis, and agent workflow status with uptime.
// @Tags health
// @Produce json
// @Success 200 {object} services.HealthStatus
// @Router /health/heartbeat [get]
func handleHeartbeat(c *gin.Context) {
	status := services.Heartbeat()
	code := http.StatusOK
	if status.Status == "degraded" {
		code = http.StatusPartialContent
	}
	c.JSON(code, status)
}

// @Summary Legacy health endpoint
// @Description Simple health check for backward compatibility and Docker healthcheck.
// @Tags health
// @Produce json
// @Success 200 {object} services.HealthStatus
// @Router /health [get]
func handleLegacyHealth(c *gin.Context) {
	status := services.LegacyHealth()
	if status.Status == "degraded" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	c.JSON(http.StatusOK, status)
}

// @Summary System information
// @Description Returns runtime system information including memory, goroutines, and uptime.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/system [get]
func handleSystemInfo(c *gin.Context) {
	info := services.SystemInfo()
	c.JSON(http.StatusOK, info)
}

// @Summary Shard health
// @Description Returns per-shard metrics including query counts, latency percentiles, and active alerts.
// @Tags health
// @Produce json
// @Success 200 {object} config.ShardHealth
// @Router /health/shards [get]
func handleShardHealth(c *gin.Context) {
	monitor := config.GetShardMonitor()
	if monitor == nil {
		c.JSON(http.StatusOK, config.ShardHealth{
			Mode:       "single",
			ShardCount: 1,
		})
		return
	}
	c.JSON(http.StatusOK, monitor.Health())
}
