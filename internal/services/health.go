package services

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"securewallet/internal/config"
)

var (
	startTime     time.Time
	ready         atomic.Bool
	workflowOK    atomic.Bool
	lastHeartbeat atomic.Value
)

func init() {
	startTime = time.Now()
	lastHeartbeat.Store(time.Now())
}

type HealthStatus struct {
	Status    string           `json:"status"`
	Checks    map[string]Check `json:"checks,omitempty"`
	Uptime    string           `json:"uptime,omitempty"`
	Version   string           `json:"version,omitempty"`
	Timestamp string           `json:"timestamp"`
}

type Check struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

func InitHealthService() {
	ready.Store(true)
	workflowOK.Store(true)
	lastHeartbeat.Store(time.Now())
}

func MarkNotReady() {
	ready.Store(false)
}

func MarkReady() {
	ready.Store(true)
}

func MarkWorkflowOK() {
	workflowOK.Store(true)
	lastHeartbeat.Store(time.Now())
}

func MarkWorkflowDegraded() {
	workflowOK.Store(false)
}

func Liveness() HealthStatus {
	return HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func Readiness() HealthStatus {
	checks := make(map[string]Check)
	overall := "ok"

	dbCheck := checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "ok" {
		overall = "fail"
	}

	redisCheck := checkRedis()
	checks["redis"] = redisCheck
	if redisCheck.Status != "ok" {
		overall = "fail"
	}

	if !ready.Load() {
		overall = "fail"
		checks["application"] = Check{
			Status:  "fail",
			Details: "application not ready",
		}
	} else {
		checks["application"] = Check{Status: "ok"}
	}

	return HealthStatus{
		Status:    overall,
		Checks:    checks,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func Heartbeat() HealthStatus {
	checks := make(map[string]Check)
	overall := "ok"

	dbCheck := checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "ok" {
		overall = "degraded"
	}

	redisCheck := checkRedis()
	checks["redis"] = redisCheck
	if redisCheck.Status != "ok" {
		overall = "degraded"
	}

	shardCheck := checkShards()
	checks["shards"] = shardCheck
	if shardCheck.Status == "degraded" {
		if overall == "ok" {
			overall = "degraded"
		}
	}

	if !workflowOK.Load() {
		overall = "degraded"
	}
	checks["workflow"] = Check{
		Status:  map[bool]string{true: "ok", false: "degraded"}[workflowOK.Load()],
		Details: fmt.Sprintf("last_heartbeat=%s", lastHeartbeat.Load().(time.Time).UTC().Format(time.RFC3339)),
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production"
	}

	return HealthStatus{
		Status:    overall,
		Checks:    checks,
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		Version:   os.Getenv("APP_VERSION"),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func LegacyHealth() HealthStatus {
	status := "healthy"
	r := Readiness()
	if r.Status != "ok" {
		status = "degraded"
	}

	return HealthStatus{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func checkDatabase() Check {
	db := config.GetDB()
	if db == nil {
		return Check{
			Status:  "fail",
			Details: "database connection not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return Check{
			Status:  "fail",
			Details: fmt.Sprintf("failed to get sql.DB: %v", err),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return Check{
			Status:  "fail",
			Details: fmt.Sprintf("database ping failed: %v", err),
		}
	}

	stats := sqlDB.Stats()
	return Check{
		Status:  "ok",
		Details: fmt.Sprintf("open_conns=%d,in_use=%d,idle=%d", stats.OpenConnections, stats.InUse, stats.Idle),
	}
}

func checkRedis() Check {
	client := config.GetRedis()
	if client == nil {
		return Check{
			Status:  "fail",
			Details: "redis client not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return Check{
			Status:  "fail",
			Details: fmt.Sprintf("redis ping failed: %v", err),
		}
	}

	return Check{Status: "ok"}
}

func SystemInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"go_version":  runtime.Version(),
		"go_os":       runtime.GOOS,
		"go_arch":     runtime.GOARCH,
		"goroutines":  runtime.NumGoroutine(),
		"num_cpu":     runtime.NumCPU(),
		"alloc_bytes": m.Alloc,
		"total_alloc": m.TotalAlloc,
		"sys_bytes":   m.Sys,
		"gc_cycles":   m.NumGC,
		"uptime":      time.Since(startTime).Round(time.Second).String(),
		"start_time":  startTime.UTC().Format(time.RFC3339),
		"environment": getEnv("ENVIRONMENT", "production"),
		"app_version": getEnv("APP_VERSION", "unknown"),
	}
}

func checkShards() Check {
	monitor := config.GetShardMonitor()
	if monitor == nil {
		return Check{Status: "ok", Details: "shard monitoring not active"}
	}

	h := monitor.Health()
	if len(h.Alerts) > 0 {
		detail := fmt.Sprintf("alerts=%d", len(h.Alerts))
		for _, a := range h.Alerts {
			detail += fmt.Sprintf(" [shard_%d:%s]", a.ShardIdx, a.Type)
		}
		return Check{Status: "degraded", Details: detail}
	}

	detail := fmt.Sprintf("mode=%s shards=%d queries_ok", h.Mode, h.ShardCount)
	return Check{Status: "ok", Details: detail}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
