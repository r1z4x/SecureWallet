package config

import (
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	defaultAlertCooldown       = 300 * time.Second
	defaultAlertWindow         = 60 * time.Second
	defaultUnevenThreshold     = 3.0
	defaultLatencyThresholdMs  = 100.0
	defaultLatencyWindowSize   = 256
	minShardCountForAlert      = 2
)

type shardMetrics struct {
	queryCount   atomic.Uint64
	latencySum   atomic.Uint64
	latencyCount atomic.Uint64
	latencyRing  []time.Duration
	ringPos      uint32
	lastAlerted  time.Time
}

type ShardMonitor struct {
	mu             sync.RWMutex
	metrics        []*shardMetrics
	window         time.Duration
	latencyWindow  int
	unevenThresh   float64
	latencyThresh  time.Duration
	alertCooldown  time.Duration
	started        atomic.Bool
	stopCh         chan struct{}
	lastSnapshot   []SnapshotEntry
}

type SnapshotEntry struct {
	ShardIndex int           `json:"shard_index"`
	QueryCount uint64        `json:"query_count"`
	AvgLatency time.Duration `json:"avg_latency_ns"`
	P50Latency time.Duration `json:"p50_latency_ns"`
	P95Latency time.Duration `json:"p95_latency_ns"`
	P99Latency time.Duration `json:"p99_latency_ns"`
}

type ShardHealth struct {
	Mode       string          `json:"mode"`
	ShardCount int             `json:"shard_count"`
	Snapshot   []SnapshotEntry `json:"snapshot"`
	Alerts     []ShardAlert    `json:"alerts,omitempty"`
}

type ShardAlert struct {
	Type      string  `json:"type"`
	Severity  string  `json:"severity"`
	Message   string  `json:"message"`
	ShardIdx  int     `json:"shard_idx,omitempty"`
	MaxQPS    float64 `json:"max_qps,omitempty"`
	MeanQPS   float64 `json:"mean_qps,omitempty"`
	Ratio     float64 `json:"ratio,omitempty"`
	P95Ms     float64 `json:"p95_ms,omitempty"`
	Threshold float64 `json:"threshold_ms,omitempty"`
}

var shardMonitor *ShardMonitor

func getEnvFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	s, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return time.Duration(s) * time.Second
}

func InitShardMonitor() *ShardMonitor {
	m := &ShardMonitor{
		window:        getEnvDuration("SHARD_ALERT_WINDOW_SECONDS", defaultAlertWindow),
		latencyWindow: defaultLatencyWindowSize,
		unevenThresh:  getEnvFloat("SHARD_ALERT_UNEVEN_THRESHOLD", defaultUnevenThreshold),
		latencyThresh: time.Duration(getEnvFloat("SHARD_ALERT_LATENCY_THRESHOLD_MS", defaultLatencyThresholdMs)) * time.Millisecond,
		alertCooldown: getEnvDuration("SHARD_ALERT_COOLDOWN_SECONDS", defaultAlertCooldown),
		stopCh:        make(chan struct{}),
	}

	sm := GetShardManager()
	count := 1
	if sm != nil {
		count = sm.ShardCount()
		if count == 0 {
			count = 1
		}
	}

	m.metrics = make([]*shardMetrics, count)
	for i := range m.metrics {
		m.metrics[i] = &shardMetrics{
			latencyRing: make([]time.Duration, m.latencyWindow),
		}
	}

	shardMonitor = m
	return m
}

func GetShardMonitor() *ShardMonitor {
	return shardMonitor
}

func (m *ShardMonitor) Start() {
	if !m.started.CompareAndSwap(false, true) {
		return
	}
	go m.loop()
	log.Println("ShardMonitor: started alert loop")
}

func (m *ShardMonitor) Stop() {
	if m.started.CompareAndSwap(true, false) {
		close(m.stopCh)
		log.Println("ShardMonitor: stopped")
	}
}

func (m *ShardMonitor) loop() {
	ticker := time.NewTicker(m.window)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.evaluate()
		}
	}
}

func (m *ShardMonitor) RecordQuery(userID uuid.UUID) {
	sm := GetShardManager()
	if sm == nil {
		return
	}
	idx := sm.ShardIndexForUser(userID)
	m.mu.RLock()
	if idx >= 0 && idx < len(m.metrics) {
		m.metrics[idx].queryCount.Add(1)
	}
	m.mu.RUnlock()
}

func (m *ShardMonitor) RecordLatency(userID uuid.UUID, d time.Duration) {
	sm := GetShardManager()
	if sm == nil {
		return
	}
	idx := sm.ShardIndexForUser(userID)
	if idx < 0 || idx >= len(m.metrics) {
		return
	}
	met := m.metrics[idx]
	met.latencySum.Add(uint64(d.Nanoseconds()))
	met.latencyCount.Add(1)
	pos := atomic.AddUint32(&met.ringPos, 1) % uint32(len(met.latencyRing))
	met.latencyRing[pos] = d
}

func (m *ShardMonitor) snapshot() []SnapshotEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.metrics) == 0 {
		return nil
	}

	entries := make([]SnapshotEntry, len(m.metrics))
	for i, met := range m.metrics {
		e := SnapshotEntry{ShardIndex: i, QueryCount: met.queryCount.Load()}
		count := met.latencyCount.Load()
		if count > 0 {
			e.AvgLatency = time.Duration(met.latencySum.Load() / count)
		}
		e.P50Latency, e.P95Latency, e.P99Latency = m.percentiles(met)
		entries[i] = e
	}
	return entries
}

func (m *ShardMonitor) percentiles(met *shardMetrics) (p50, p95, p99 time.Duration) {
	durations := make([]time.Duration, 0, len(met.latencyRing))
	for _, d := range met.latencyRing {
		if d > 0 {
			durations = append(durations, d)
		}
	}
	if len(durations) == 0 {
		return 0, 0, 0
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p50 = durations[len(durations)*50/100]
	p95 = durations[len(durations)*95/100]
	p99 = durations[len(durations)*99/100]
	return
}

func (m *ShardMonitor) evaluate() {
	m.mu.RLock()
	n := len(m.metrics)
	m.mu.RUnlock()

	if n < minShardCountForAlert {
		return
	}

	snap := m.snapshot()
	if len(snap) == 0 {
		return
	}

	m.lastSnapshot = snap

	var maxQPS float64
	var totalQPS float64
	for _, s := range snap {
		qps := float64(s.QueryCount) / m.window.Seconds()
		totalQPS += qps
		if qps > maxQPS {
			maxQPS = qps
		}
	}
	meanQPS := totalQPS / float64(len(snap))

	if meanQPS > 0 {
		ratio := maxQPS / meanQPS
		if ratio >= m.unevenThresh {
			for i := range snap {
				qps := float64(snap[i].QueryCount) / m.window.Seconds()
				if qps == maxQPS {
					if m.canAlert(i) {
						m.markAlerted(i)
						log.Printf("[SHARD_ALERT] Uneven usage: shard %d at %.1f QPS (%.1fx mean of %.1f QPS)",
							i, maxQPS, ratio, meanQPS)
					}
				}
			}
		}
	}

	for i, s := range snap {
		p95Ms := float64(s.P95Latency) / float64(time.Millisecond)
		if p95Ms > 0 && s.P95Latency >= m.latencyThresh {
			if m.canAlert(i) {
				m.markAlerted(i)
				log.Printf("[SHARD_ALERT] Latency spike: shard %d p95=%.1fms (threshold %.1fms)",
					i, p95Ms, float64(m.latencyThresh)/float64(time.Millisecond))
			}
		}
	}
}

func (m *ShardMonitor) canAlert(idx int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if idx < 0 || idx >= len(m.metrics) {
		return false
	}
	return time.Since(m.metrics[idx].lastAlerted) >= m.alertCooldown
}

func (m *ShardMonitor) markAlerted(idx int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if idx >= 0 && idx < len(m.metrics) {
		m.metrics[idx].lastAlerted = time.Now()
	}
}

func (m *ShardMonitor) ResetCounters() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, met := range m.metrics {
		met.queryCount.Store(0)
		met.latencySum.Store(0)
		met.latencyCount.Store(0)
	}
}

func (m *ShardMonitor) Health() ShardHealth {
	sm := GetShardManager()
	mode := "single"
	shardCount := 1
	if sm != nil {
		if sm.Mode() == ShardModeMulti {
			mode = "multi"
		}
		shardCount = sm.ShardCount()
		if shardCount == 0 {
			shardCount = 1
		}
	}

	h := ShardHealth{
		Mode:       mode,
		ShardCount: shardCount,
		Snapshot:   m.snapshot(),
	}

	if len(h.Snapshot) >= minShardCountForAlert {
		var maxQPS float64
		var totalQPS float64
		for _, s := range h.Snapshot {
			qps := float64(s.QueryCount) / m.window.Seconds()
			totalQPS += qps
			if qps > maxQPS {
				maxQPS = qps
			}
		}
		meanQPS := totalQPS / float64(len(h.Snapshot))
		if meanQPS > 0 {
			ratio := maxQPS / meanQPS
			if ratio >= m.unevenThresh {
				for i := range h.Snapshot {
					qps := float64(h.Snapshot[i].QueryCount) / m.window.Seconds()
					if qps == maxQPS {
						h.Alerts = append(h.Alerts, ShardAlert{
							Type:     "uneven_usage",
							Severity: "warning",
							Message:  "shard receiving disproportionate traffic",
							ShardIdx: i,
							MaxQPS:   roundFloat(maxQPS),
							MeanQPS:  roundFloat(meanQPS),
							Ratio:    roundFloat(ratio),
						})
					}
				}
			}
		}

		for i, s := range h.Snapshot {
			p95Ms := float64(s.P95Latency) / float64(time.Millisecond)
			if p95Ms > 0 && s.P95Latency >= m.latencyThresh {
				h.Alerts = append(h.Alerts, ShardAlert{
					Type:      "latency_spike",
					Severity:  "warning",
					Message:   "p95 latency exceeds threshold",
					ShardIdx:  i,
					P95Ms:     roundFloat(p95Ms),
					Threshold: float64(m.latencyThresh) / float64(time.Millisecond),
				})
			}
		}
	}

	return h
}

func roundFloat(f float64) float64 {
	return math.Round(f*100) / 100
}
