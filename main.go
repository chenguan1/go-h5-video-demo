package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"gocv.io/x/gocv"
	"log"
	"strings"
)

var vpath = "E:/go/src/go-h5-video-demo/data/vtest.avi"

func playVideo()  {
	w := gocv.NewWindow("video")

	img := gocv.IMRead("./data/500w01.jpg",1)
	for{
		
		w.IMShow(img)
		w.WaitKey(1)
	}

}


func main() {

	go playVideo();

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			ping := string(msg.Body)
			pong := strings.Replace(ping,"？", "！", len(ping))
			pong = strings.Replace(pong, "么","",len(pong))

			mg := websocket.Message{
				Body:[]byte(pong),
				IsNative:true,
			}

			nsConn.Conn.Write(mg)

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


	app := iris.New()
	app.HandleDir("/","./html")
	app.Get("/msg", websocket.Handler(ws))

	app.Run(iris.Addr(":8080"))
}
