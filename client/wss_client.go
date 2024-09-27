package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

func main() {
    // 连接 WSS 服务器
    wssConfig := &tls.Config{
        InsecureSkipVerify: true, // 注意：在生产中请不要使用此设置
    }
    wssDialer := &websocket.Dialer{
        TLSClientConfig: wssConfig,
    }

    wss, _, err := wssDialer.Dial("wss://101.32.223.170:443/wss", nil)
    if err != nil {
        log.Fatalf("Failed to connect to WSS server: %v", err)
    }
    defer wss.Close()

    // 持续发送消息
    for {
        msg := []byte("Hello from WSS!")
        err = wss.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Fatal("WriteMessage: ", err)
        }
        fmt.Printf("Sent: %s\n", msg)

        // 等待1秒钟后再发送下一条消息
        time.Sleep(1 * time.Second)

        // 可选：读取服务器的响应
        _, wssMsg, err := wss.ReadMessage()
        if err != nil {
            log.Fatal("ReadMessage: ", err)
        }
        fmt.Printf("Received from WSS server: %s\n", wssMsg)
    }
}
