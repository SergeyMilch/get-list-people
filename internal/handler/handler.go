package handler

import (
	"github.com/SergeyMilch/get-list-people-effective-mobile/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

func SetupRoutes(router *gin.Engine, db *sqlx.DB, rdb *redis.Client) {
	router.GET("/people", api.GetPeople(db, rdb))
	router.GET("/people/:id", api.GetPersonByID(db, rdb))
	router.POST("/people", api.AddPerson(db, rdb))
	router.DELETE("/people/:id", api.DeletePerson(db, rdb))
	router.PUT("/people/:id", api.UpdatePerson(db, rdb))

	router.POST("/graphql", api.HandleGraphQL(db, rdb))
}
