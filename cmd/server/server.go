package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/wwwthomson/monitoring/pkg/agent"
	"fmt"
)

func HttpServer(data *Monitor) *gin.Engine {
	route := gin.Default()
	route.POST("/network", Network(data))
	route.POST("/memory", Memory(data))
	route.POST("/swap", Swap(data))
	route.POST("/cpu", CPU(data))
	return route
}

func Network(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Network
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		data.AddNetwork(request)
	}
}

func Memory(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Memory
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		data.AddMemory(request)
	}
}

func Swap(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Swap
		err := c.ShouldBindJSON(&request)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		data.AddSwap(request)
	}
}

func CPU(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.CPU
		err := c.ShouldBindJSON(&request)
		if err == nil {
			fmt.Printf("%#v\n", request)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		data.AddCPU(request)
	}
}
