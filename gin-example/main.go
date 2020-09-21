package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/?type=catalog&mode=checkauth", func(c *gin.Context) {
		c.XML(200, gin.H{
			"message": "OK",
		})
	})
	r.Run(":8070") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
