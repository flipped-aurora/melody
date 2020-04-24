package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	profileEngine := gin.Default()

	profileEngine.GET("/coupon/:uuid", func(context *gin.Context) {
		uuid := context.Param("uuid")
		context.JSON(http.StatusOK, map[string]interface{}{
			"id":     1,
			"uuid":   uuid,
			"coupons": []map[string]interface{}{
				{
					"title": "5元优惠券",
					"condition": "age > 20",
					"start": "xxxx-xx-xx",
					"end": "xxxx-xx-xx",
				},
				{
					"title": "7元优惠券",
					"condition": "none",
					"start": "xxxx-xx-xx",
					"end": "xxxx-xx-xx",
				},
				{
					"title": "7元优惠券",
					"condition": "none",
					"start": "xxxx-xx-xx",
					"end": "xxxx-xx-xx",
				},
				{
					"title": "5元优惠券",
					"condition": "vip true",
					"start": "xxxx-xx-xx",
					"end": "xxxx-xx-xx",
				},
			},
		})
	})
	fmt.Println(profileEngine.Run(":9003"))
}
