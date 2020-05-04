package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

type clientPage struct {
	Title string
	Host  string
}

func main() {
	app := iris.New()

	//app.RegisterView(iris.HTML("./templates", ".html")) // select the html engine to serve templates
	app.HandleDir("/","./html")

	// Almost all features of neffos are disabled because no custom message can pass
	// when app expects to accept and send only raw websocket native messages.
	// When only allow native messages is a fact?
	// When the registered namespace is just one and it's empty
	// and contains only one registered event which is the `OnNativeMessage`.
	// When `Events{...}` is used instead of `Namespaces{ "namespaceName": Events{...}}`
	// then the namespace is empty "".
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			nsConn.Conn.Server().Broadcast(nsConn, msg)


			data,err := ioutil.ReadFile("./img.png")
			if err == nil{
				n := base64.StdEncoding.EncodedLen(len(data)) // 计算加密后数据的长度
				dst := make([]byte,n) // 创建容器
				base64.StdEncoding.Encode(dst,data) // 加密数据
				nsConn.Emit("img",dst)

			}else{
				fmt.Println("err")
			}

			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server", c.ID())
	}

	// register the server on an endpoint.
	// see the inline javascript code i the websockets.html, this endpoint is used to connect to the server.
	app.Get("/my_endpoint", websocket.Handler(ws))


	// Target some browser windows/tabs to http://localhost:8080 and send some messages,
	// see the static/js/chat.js,
	// note that the client is using only the browser's native WebSocket API instead of the neffos one.
	app.Run(iris.Addr(":8080"))
}
