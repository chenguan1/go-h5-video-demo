package main

import (
	"encoding/base64"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"gocv.io/x/gocv"
	"log"
	"time"
)

func main() {
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
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

	ws.OnUpgradeError = func(err error) {
		log.Printf("Upgrade Error: %v", err)
	}

	go func() {
		camera, err := gocv.VideoCaptureDevice(0)
		if err != nil {
			panic(err)
		}

		img := gocv.NewMat()

		for {
			camera.Read(&img)
			data, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				fmt.Println(err)
			} else {
				n := base64.StdEncoding.EncodedLen(len(data))
				dst := make([]byte, n)
				base64.StdEncoding.Encode(dst, data)
				urldata := "data:image/jpeg;base64," + string(dst)
				mg := websocket.Message{
					Body:     []byte(urldata),
					IsNative: true,
				}
				ws.Broadcast(nil, mg)
			}
			time.Sleep(time.Millisecond * time.Duration(50))
		}
	}()

	app := iris.New()
	app.HandleDir("/", "./html")
	app.Get("/video", websocket.Handler(ws))

	app.Run(iris.Addr(":8080"))
}
