# go-h5-video-demo
gocv + websoket + h5 实现视频直播

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
### 参考
https://github.com/kataras/iris/tree/v12/_examples/websocket
