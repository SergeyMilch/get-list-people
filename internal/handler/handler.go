package handler

import (
	"github.com/SergeyMilch/get-list-people-effective-mobile/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRoutes(router *gin.Engine, db *sqlx.DB) {
	router.GET("/people", api.GetPeople(db))
	router.POST("/people", api.AddPerson(db))
	router.DELETE("/people/:id", api.DeletePerson(db))
	router.PUT("/people/:id", api.UpdatePerson(db))
}
