package main
/*
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

*/

/*
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"gocv.io/x/gocv"
	"log"
	"math/rand"
	"sync"
	"time"
)

var vpath = "E:/go/src/go-h5-video-demo/data/vtest.avi"

var lk sync.Mutex
var urldata string

func playVideo()  {
	rand.Seed(time.Now().Unix())

	//w := gocv.NewWindow("video")

	img := gocv.IMRead("./data/500w01.jpg",1)
	if img.Empty(){
		panic("image is empty")
	}


	for{
		x := int(rand.Float64() * float64(img.Cols()-1))
		y := int(rand.Float64() * float64(img.Rows()-1))
		//z := int(rand.Float64() * float64(2))
		v1 := int(rand.Float64() * float64(255))
		v2 := int(rand.Float64() * float64(255))
		v3 := int(rand.Float64() * float64(255))
		v := v1 << 16 + v2 << 8 + v3
		img.SetIntAt(y, x, int32(v))
		//fmt.Println(x,y,z,v)
		//img.SetUCharAt(y,x,uint8(v))
		//w.IMShow(img)
		//w.WaitKey(1)

		data, err := gocv.IMEncode(".jpg",img)
		if err != nil{
			fmt.Println(err)
		}else {
			n := base64.StdEncoding.EncodedLen(len(data)) // 计算加密后数据的长度
			dst := make([]byte,n) // 创建容器
			base64.StdEncoding.Encode(dst,data) // 加密数据
			lk.Lock()
			urldata = "data:image/jpeg;base64,"+string(dst)
			lk.Unlock()
		}

	}

}


func main() {

	go playVideo();

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			go func() {
				for{
					mg := websocket.Message{
						Body:[]byte(urldata),
						IsNative:true,
					}

					nsConn.Conn.Write(mg)
					time.Sleep(time.Millisecond * time.Duration(100))

					if nsConn.Conn.IsClosed() {
						break
					}
				}

			}()



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

*/