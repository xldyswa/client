package main

import (
    "fmt"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

func main() {
    // 连接到 WebSocket 服务器
    url := "ws://101.32.223.170:8080/ws"
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatal("Error connecting to WebSocket server:", err)
    }
    defer conn.Close()

    // 启动一个 goroutine 来处理服务器发送的消息
    go func() {
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                log.Println("Error while reading message:", err)
                return
            }
            fmt.Printf("Received from server: %s\n", msg)
        }
    }()

    // 循环发送消息到服务器
    for {
        msg := []byte("Hello, WebSocket server!")
        err := conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Println("Error while writing message:", err)
            return
        }

        // 每隔 2 秒发送一次消息
        time.Sleep(2 * time.Second)
    }
}
