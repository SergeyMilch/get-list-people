package router

import (
	"context"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db/redisdb"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, db db.Database, rdb redisdb.RedisClient) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)   // "db" в контекст (для GraphQL)
		c.Set("rdb", rdb) // "rdb" в контекст (для GraphQL)
		c.Next()
	})
	handler.SetupRoutes(router, db, rdb)
	return router
}
