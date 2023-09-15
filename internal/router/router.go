package router

import (
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db) // "db" в контекст (для GraphQL)
		c.Next()
	})
	handler.SetupRoutes(router, db)
	return router
}
