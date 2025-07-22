package handlers

import (
	"database/sql"

	"zplus-saas/apps/backend/shared/config"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	Auth   *AuthHandler
	Tenant *TenantHandler
	Module *ModuleHandler
	Proxy  *ProxyHandler
}

func New(db *sql.DB, redis *redis.Client, mongo *mongo.Database, cfg *config.Config) *Handlers {
	return &Handlers{
		Auth:   NewAuthHandler(db, redis, cfg),
		Tenant: NewTenantHandler(db, cfg),
		Module: NewModuleHandler(db, cfg),
		Proxy:  NewProxyHandler(cfg),
	}
}
