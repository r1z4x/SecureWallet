package services

import (
	"fmt"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShardResolver struct{}

func NewShardResolver() *ShardResolver {
	return &ShardResolver{}
}

func (r *ShardResolver) DBForUser(userID uuid.UUID) *gorm.DB {
	return config.GetShardManager().GetDB(userID)
}

func (r *ShardResolver) DefaultDB() *gorm.DB {
	return config.GetShardManager().GetDefaultDB()
}

func (r *ShardResolver) IsSharded() bool {
	return config.GetShardManager().Mode() == config.ShardModeMulti
}

func (r *ShardResolver) ResolveUserDB(email string) (*gorm.DB, *models.User, error) {
	defaultDB := r.DefaultDB()
	user, err := findUserByEmail(defaultDB, email)
	if err != nil {
		return nil, nil, err
	}
	return r.DBForUser(user.ID), user, nil
}

func (r *ShardResolver) ResolveUserDBByUsername(username string) (*gorm.DB, *models.User, error) {
	defaultDB := r.DefaultDB()
	user, err := findUserByUsername(defaultDB, username)
	if err != nil {
		return nil, nil, err
	}
	return r.DBForUser(user.ID), user, nil
}

func (r *ShardResolver) FindUserAcrossShards(username string) (*models.User, *gorm.DB, error) {
	sm := config.GetShardManager()

	if !r.IsSharded() {
		db := r.DefaultDB()
		user, err := findUserByUsername(db, username)
		if err != nil {
			return nil, nil, err
		}
		return user, db, nil
	}

	// In single-DB mode or as a last resort, try default DB first
	defaultDB := r.DefaultDB()
	user, err := findUserByUsername(defaultDB, username)
	if err == nil {
		return user, defaultDB, nil
	}

	// Search shards sequentially; this is acceptable for login (infrequent)
	shardCount := sm.ShardCount()
	for i := 0; i < shardCount; i++ {
		shardDB, err := sm.GetShardByIndex(i)
		if err != nil {
			continue
		}
		user, err := findUserByUsername(shardDB, username)
		if err == nil {
			expectedIdx := sm.ShardIndexForUser(user.ID)
			if expectedIdx == i {
				return user, shardDB, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("user not found: %s", username)
}

func findUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found by email: %s", email)
	}
	return &user, nil
}

func findUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found by username: %s", username)
	}
	return &user, nil
}


