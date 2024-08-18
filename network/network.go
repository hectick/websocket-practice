package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	n := &Network{
		engin: gin.New(),
	}

	n.engin.Use(gin.Logger())   //user로 모든 api나 모든 라우터 등에 범용적인 처리 가능
	n.engin.Use(gin.Recovery()) // panic이나 특정 로직으로 인해 서버가 죽었을때 다시 서버 올려주는거
	n.engin.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"}, //모두 허용
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	return n
}

func (n *Network) StartServer() error {
	log.Println("Starting server")
	return n.engin.Run(":8080")
}
