package main

import (
	"github.com/gin-gonic/gin"
	//"log"
	"net/http"
	"github.com/wwwthomson/monitoring/pkg/agent"
	"fmt"
)



func HttpServer(data *Store) *gin.Engine {
	route := gin.Default()
	route.POST("/network", Network(data))
	route.POST("/memory", Memory(data))
	route.POST("/swap", Swap)
	route.POST("/cpu", CPU)
	return route
}

func Network(data *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json agent.Network
		err := c.ShouldBindJSON(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		fmt.Println(data)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
	}
}

func Memory(data *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json agent.Memory
		err := c.ShouldBindJSON(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		fmt.Println(data)
	}
}

func Swap(c *gin.Context) {
	var json agent.Swap
	if err := c.ShouldBindJSON(&json); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ooops. I'm sorry =(",
		})
	}
}
func CPU(c *gin.Context) {
	var request agent.CPU
	if err := c.ShouldBindJSON(&request); err == nil {
		fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ooops. I'm sorry =(",
		})
	}
}
