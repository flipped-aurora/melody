package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	profileEngine := gin.Default()

	profileEngine.GET("/role/:id", func(context *gin.Context) {
		roleId := context.Param("id")
		context.JSON(http.StatusOK, map[string]interface{}{
			"role_id":   roleId,
			"role_name": "Administrator",
			"role_tag":  "Admin",
			"authorities": []map[string]interface{}{
				{
					"tag":    "user add",
					"path":   "/user/add",
					"method": "post",
					"active": true,
				},
				{
					"tag":    "user delete",
					"path":   "/user/delete",
					"method": "get",
					"active": true,
				},
			},
		})
	})
	fmt.Println(profileEngine.Run(":9002"))
}
