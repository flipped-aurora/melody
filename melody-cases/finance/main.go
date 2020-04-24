package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	profileEngine := gin.Default()

	profileEngine.GET("/finance/:uuid", func(context *gin.Context) {
		uuid := context.Param("uuid")
		context.JSON(http.StatusOK, map[string]interface{}{
			"id": 1,
			"uuid": uuid,
			"name": "Grant",
			"level": "Gold member",
			"active_balance": 5000.22,
			"frozen_balance": 1000.00,
			"borrow_balance": 50000.00,
		})
	})
	fmt.Println(profileEngine.Run(":9004"))
}
