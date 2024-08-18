package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/types"
	"net/http"
)

// http 커넥션을 webSocket 커넥션으로 업그레이드!
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Room struct {
	Forward chan *message    //수신되는 메시지를 보관, 다른 클라이언트들에게 전송
	Join    chan *Client     //소켓이 연결되는 경우 작동
	Leave   chan *Client     //소켓이 끊어지는 경우 작동
	Clients map[*Client]bool //현재 방에있는 클라이언트 정보 저장
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client
	defer func() { r.Leave <- client }() // 커넥션이 끊김으로써 함수를 벗어나는 상황이 발생했을때 leave
}
