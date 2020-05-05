# go-h5-video-demo
gocv + websoket + h5 实现视频直播

[toc]
### 数据流向图
![数据流](https://img-blog.csdnimg.cn/20200505215829619.png)
### 实现方式
1. 通过opencv抓取摄像头的视频数据按帧处理
2. 将每一帧压缩成jpg格式并编码成base64格式
3. 通过websocket协议将base64图像传输给前端页面
4. 前端解析每一帧并更新显示

本文使用了iris框架的websocket封装，因此opencv也使用了go语言的版本GoCV。

### 主要功能点
#### 从摄像头获取视频数据
```go
img:=gocv.NewMat()
camera,err:=gocv.VideoCaptureDevice(0)
camera.Read(&img)
```
#### 图像Base64编码
```go
data,err:=gocv.IMEncode(".jpg",img)
n:=base64.StdEncoding.EncodedLen(len(data))
dst:=make([]byte,n)
base64.StdEncoding.Encode(dst,data)
urldata:="data:image/jpeg;base64,"+string(dst)
```

#### Websocket服务
```go
ws:=websocket.New(websocket.DefaultGorillaUpgrader,websocket.Events{
 websocket.OnNativeMessage:func(nsConn*websocket.NSConn,msgwebsocket.Message)error{
  log.Printf("Servergot:%sfrom[%s]",msg.Body,nsConn.Conn.ID())
  return nil
 },
})
app:=iris.New()
app.Get("/video",websocket.Handler(ws))
app.Run(iris.Addr(":8080"))
```
#### 通过websocket广播
```go
mg:=websocket.Message{
 Body:[]byte(urldata),
 IsNative:true,
}
ws.Broadcast(nil,mg)
```
### 完整代码
#### go
```go
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
     Body: []byte(urldata),
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
```
#### html5
```html
<html>
<head>
    <title>video</title>
</head>
<body style="padding:10px;">
<img id="target" style="display:inline;"/>
<canvas id="canvas" style="display:inline;"/>
<script type="text/javascript">
    var HOST = "localhost:8080"

    w = new WebSocket("ws://" + HOST + "/video");
    w.onopen = function () {
        console.log("Websocket connection enstablished");
    };

    w.onclose = function () {
        console.log("Websocket disconnected");
    };

    var canvas = document.getElementById("canvas");
    var img = new Image()
    img.onload = function(){
        canvas.width = img.width
        canvas.height = img.height
        canvas.getContext("2d").drawImage(img,0,0,img.width,img.height);
    };
    w.onmessage = function (message) {
        img.src = message.data;
    };
</script>
</body>
</html>
```

### 效果图
![效果图](https://img-blog.csdnimg.cn/2020050522104329.gif)
### 参考
https://github.com/kataras/iris/tree/v12/_examples/websocket
