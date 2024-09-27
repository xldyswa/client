package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    // 创建一个自定义的 HTTP 客户端
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true, // 忽略证书验证（只用于开发环境）
            },
        },
    }

    // 发送 GET 请求
    resp, err := client.Get("https://101.32.223.170:443")
    if err != nil {
        fmt.Println("Failed to connect to HTTPS server:", err)
        return
    }
    defer resp.Body.Close()

    // 读取响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Failed to read response body:", err)
        return
    }

    fmt.Println("Response from server:", string(body))
}
