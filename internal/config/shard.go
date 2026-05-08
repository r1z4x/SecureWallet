package config

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ShardMode int

const (
	ShardModeSingle ShardMode = iota
	ShardModeMulti
)

type ShardManager struct {
	mu        sync.RWMutex
	mode      ShardMode
	shards    []*gorm.DB
	shardDSNs []string
}

var shardMgr *ShardManager

func GetShardManager() *ShardManager {
	return shardMgr
}

func (sm *ShardManager) Mode() ShardMode {
	return sm.mode
}

func (sm *ShardManager) ShardCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.shards)
}

func (sm *ShardManager) ShardIndexForUser(userID uuid.UUID) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.mode == ShardModeSingle || len(sm.shards) == 0 {
		return 0
	}
	return hashUserIDToShard(userID, len(sm.shards))
}

func (sm *ShardManager) GetDB(userID uuid.UUID) *gorm.DB {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.mode == ShardModeSingle || len(sm.shards) == 0 {
		return DB
	}

	idx := hashUserIDToShard(userID, len(sm.shards))
	return sm.shards[idx]
}

func (sm *ShardManager) GetDefaultDB() *gorm.DB {
	return DB
}

func (sm *ShardManager) GetShardByIndex(idx int) (*gorm.DB, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if idx < 0 || idx >= len(sm.shards) {
		return nil, fmt.Errorf("shard index %d out of range [0, %d)", idx, len(sm.shards))
	}
	return sm.shards[idx], nil
}

func (sm *ShardManager) AllShards() []*gorm.DB {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make([]*gorm.DB, len(sm.shards))
	copy(result, sm.shards)
	return result
}

func hashUserIDToShard(userID uuid.UUID, shardCount int) int {
	hash := sha256.Sum256([]byte(userID.String()))
	val := binary.BigEndian.Uint64(hash[:8])
	return int(val % uint64(shardCount))
}

func InitShards() error {
	shardCountStr := os.Getenv("DB_SHARD_COUNT")
	if shardCountStr == "" {
		shardMgr = &ShardManager{mode: ShardModeSingle}
		log.Println("ShardManager: single-DB mode (DB_SHARD_COUNT not set)")
		return nil
	}

	count, err := strconv.Atoi(shardCountStr)
	if err != nil || count < 1 {
		shardMgr = &ShardManager{mode: ShardModeSingle}
		log.Printf("ShardManager: single-DB mode (DB_SHARD_COUNT=%s invalid)", shardCountStr)
		return nil
	}

	if count == 1 {
		shardMgr = &ShardManager{mode: ShardModeSingle}
		log.Println("ShardManager: single-DB mode (DB_SHARD_COUNT=1)")
		return nil
	}

	shards := make([]*gorm.DB, 0, count)
	dsns := make([]string, 0, count)

	for i := 0; i < count; i++ {
		envKey := fmt.Sprintf("DB_SHARD_%d_DSN", i)
		dsn := os.Getenv(envKey)
		if dsn == "" {
			envKeyHost := fmt.Sprintf("DB_SHARD_%d_HOST", i)
			envKeyPort := fmt.Sprintf("DB_SHARD_%d_PORT", i)
			envKeyName := fmt.Sprintf("DB_SHARD_%d_NAME", i)
			envKeyUser := fmt.Sprintf("DB_SHARD_%d_USER", i)
			envKeyPass := fmt.Sprintf("DB_SHARD_%d_PASSWORD", i)

			host := os.Getenv(envKeyHost)
			port := os.Getenv(envKeyPort)
			name := os.Getenv(envKeyName)
			user := os.Getenv(envKeyUser)
			pass := os.Getenv(envKeyPass)

			if host == "" {
				host = os.Getenv("DB_HOST")
			}
			if port == "" {
				port = os.Getenv("DB_PORT")
			}
			if name == "" {
				name = fmt.Sprintf("%s_shard_%d", os.Getenv("DB_NAME"), i)
			}
			if user == "" {
				user = os.Getenv("DB_USER")
			}
			if pass == "" {
				pass = os.Getenv("DB_PASSWORD")
			}

			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				user, pass, host, port, name)
		}

		db, err := openShardDB(dsn, i)
		if err != nil {
			closeShards(shards)
			return fmt.Errorf("failed to connect to shard %d: %w", i, err)
		}

		shards = append(shards, db)
		dsns = append(dsns, dsn)
	}

	shardMgr = &ShardManager{
		mode:      ShardModeMulti,
		shards:    shards,
		shardDSNs: dsns,
	}

	log.Printf("ShardManager: multi-DB mode initialized with %d shards", count)
	return nil
}

func openShardDB(dsn string, index int) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("ShardManager: shard %d connected", index)
	return db, nil
}

func closeShards(shards []*gorm.DB) {
	for i, db := range shards {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		log.Printf("ShardManager: closed shard %d", i)
	}
}


