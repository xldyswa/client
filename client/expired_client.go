package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

func main() {
    url := "wss://101.32.223.170:8080/ws" // 使用你的 IP 地址

    // 读取自签名证书
    caCert, err := ioutil.ReadFile("expired.crt")
    if err != nil {
        log.Fatal("Failed to read CA certificate:", err)
    }

    // 创建一个新的证书池
    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        log.Fatal("Failed to append CA certificate to pool")
    }

    // 自定义 TLS 配置
    tlsConfig := &tls.Config{
        RootCAs: caCertPool,
    }

    dialer := websocket.Dialer{
        TLSClientConfig: tlsConfig,
    }

    conn, _, err := dialer.Dial(url, nil)
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer conn.Close()

    for {
        // 发送消息
        msg := []byte("Hello, Server!")
        err := conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Println("Write error:", err)
            return
        }
        fmt.Println("Sent:", string(msg))

        // 接收消息
        _, msg, err = conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            return
        }
        fmt.Printf("Received: %s\n", msg)

        time.Sleep(2 * time.Second) // 每2秒发送一次消息
    }
}
