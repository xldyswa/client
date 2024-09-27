package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "time"

    "github.com/gorilla/websocket"
)

// 获取本机的IP地址
func getLocalIP() (string, error) {
    // 获取本机所有网络接口
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    // 遍历网络接口
    for _, iface := range interfaces {
        // 排除回环地址和非活动接口
        if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
            addresses, err := iface.Addrs()
            if err != nil {
                return "", err
            }
            for _, addr := range addresses {
                // 确保是IPv4地址
                if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
                    return ipNet.IP.String(), nil
                }
            }
        }
    }
    return "", fmt.Errorf("no valid IP address found")
}

func main() {
    // 获取本机的IP地址
    localIP, err := getLocalIP()
    if err != nil {
        log.Fatal("Failed to get local IP address:", err)
    }

    // 设置端口号
    port := 99 // 在此处直接设置端口号

    // 定义远程服务器的URL
    url := "wss://101.32.223.170:8080/ws" // 远程服务器的URL

    // 读取自签名证书
    caCert, err := ioutil.ReadFile("server.crt")
    if err != nil {
        log.Fatal("Failed to read CA certificate:", err)
    }

    // 创建新的证书池
    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        log.Fatal("Failed to append CA certificate to pool")
    }

    // 自定义TLS配置
    tlsConfig := &tls.Config{
        RootCAs: caCertPool,
    }

    // 创建一个TCP连接的拨号器
    dialer := &websocket.Dialer{
        TLSClientConfig: tlsConfig,
        NetDial: func(network, address string) (net.Conn, error) {
            // 解析远程地址
            remoteAddr, err := net.ResolveTCPAddr(network, address)
            if err != nil {
                return nil, err
            }

            // 建立TCP连接并绑定本地地址
            return net.DialTCP(network, &net.TCPAddr{IP: net.ParseIP(localIP), Port: port}, remoteAddr)
        },
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
