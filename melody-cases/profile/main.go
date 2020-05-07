package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	profileEngine := gin.Default()

	profileEngine.GET("/profile/:id", func(context *gin.Context) {
		id := context.Param("id")
		context.JSON(http.StatusOK, map[string]interface{}{
			"id":        id,
			"uuid":      "15316368801",
			"name":      "Grant",
			"age":       22,
			"role_id":   1,
			"vip":       false,
			"password":  "xxxxxxx",
			"nick_name": "Grant",
			"avatar":    "grant.jpg",
			"birth":     "xxxx-xx-xx",
		})
	})
	fmt.Println(profileEngine.Run(":9001"))
}
