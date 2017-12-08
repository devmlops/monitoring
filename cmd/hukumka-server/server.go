package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wwwthomson/monitoring/pkg/agent"
	//"fmt"
)

func HttpServer(data *Monitor) *gin.Engine {
	route := gin.Default()
	route.GET("/health", Health())
	route.POST("/network", Network(data))
	route.POST("/memory", Memory(data))
	route.POST("/swap", Swap(data))
	route.POST("/cpu", CPU(data))
	route.POST("/disk", Disk(data))
	return route
}

func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		message := OpenMessage("message.txt")
		c.String(http.StatusOK, message)
	}
}

func Network(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Network
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		//fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		//fmt.Printf(">>>> HERE 1")
		go data.AddNetwork(request)
		//fmt.Printf(">>>> HERE 2")
	}
}

func Memory(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Memory
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		//fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		go data.AddMemory(request)
	}
}

func Swap(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Swap
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		//fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		go data.AddSwap(request)
	}
}

func CPU(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.CPU
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		//fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		go data.AddCPU(request)
	}
}

func Disk(data *Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request agent.Disk
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ooops. I'm sorry =(",
			})
			return
		}
		//fmt.Printf("%#v\n", request)
		c.JSON(http.StatusOK, gin.H{
			"message": "yeees. It's OK =)",
		})
		go data.AddDisk(request)
	}
}
