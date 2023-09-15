package router

import (
	"context"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

func NewRouter(ctx context.Context, db *sqlx.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db) // "db" в контекст (для GraphQL)
		c.Next()
	})
	handler.SetupRoutes(router, db, rdb)
	return router
}
