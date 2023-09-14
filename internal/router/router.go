package router

import (
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	handler.SetupRoutes(router, db)
	return router
}
