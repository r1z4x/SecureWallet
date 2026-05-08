package middleware

import (
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const dbContextKey = "_shard_db"

func ShardMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok {
			c.Next()
			return
		}

		sm := config.GetShardManager()
		if sm == nil {
			c.Set(dbContextKey, config.GetDB())
			c.Next()
			return
		}

		monitor := config.GetShardMonitor()
		if monitor != nil {
			monitor.RecordQuery(currentUser.ID)
		}

		db := sm.GetDB(currentUser.ID)
		c.Set(dbContextKey, db)

		start := time.Now()
		c.Next()

		if monitor != nil {
			monitor.RecordLatency(currentUser.ID, time.Since(start))
		}
	}
}

func DB(c *gin.Context) *gorm.DB {
	v, exists := c.Get(dbContextKey)
	if !exists {
		return config.GetDB()
	}
	db, ok := v.(*gorm.DB)
	if !ok {
		return config.GetDB()
	}
	return db
}

func DefaultDB(c *gin.Context) *gorm.DB {
	sm := config.GetShardManager()
	if sm == nil {
		return config.GetDB()
	}
	return sm.GetDefaultDB()
}
