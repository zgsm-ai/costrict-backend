package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// HTTP客户端配置
type HTTPClientConfig struct {
    Timeout       time.Duration `json:"timeout"`
    MaxIdleConns int          `json:"max_idle_conns"`
    IdleConnTimeout time.Duration `json:"idle_conn_timeout"`
    TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout"`
    ResponseHeaderTimeout time.Duration `json:"response_header_timeout"`
    ExpectContinueTimeout time.Duration `json:"expect_continue_timeout"`
}

// HTTP客户端
type HTTPClient struct {
    client  *http.Client
    config  *HTTPClientConfig
    headers map[string]string
}

// 创建新的HTTP客户端
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
    if config == nil {
        config = &HTTPClientConfig{
            Timeout:           30 * time.Second,
            MaxIdleConns:      10,
            IdleConnTimeout:    90 * time.Second,
            TLSHandshakeTimeout: 10 * time.Second,
        }
    }
    
    client := &http.Client{
        Timeout: config.Timeout,
        Transport: &http.Transport{
            MaxIdleConns:        config.MaxIdleConns,
            IdleConnTimeout:     config.IdleConnTimeout,
            TLSHandshakeTimeout:  config.TLSHandshakeTimeout,
            ResponseHeaderTimeout: config.ResponseHeaderTimeout,
            ExpectContinueTimeout: config.ExpectContinueTimeout,
        },
    }
    
    return &HTTPClient{
        client:  client,
        config:  config,
        headers: make(map[string]string),
    }
}

// 设置请求头
func (c *HTTPClient) SetHeader(key, value string) {
    c.headers[key] = value
}

// 设置多个请求头
func (c *HTTPClient) SetHeaders(headers map[string]string) {
    for key, value := range headers {
        c.headers[key] = value
    }
}

// GET请求
func (c *HTTPClient) Get(url string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("创建请求失败: %v", err)
    }
    
    // 设置请求头
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    return c.client.Do(req)
}

// POST请求
func (c *HTTPClient) Post(url string, body interface{}) (*http.Response, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("序列化请求体失败: %v", err)
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("创建请求失败: %v", err)
    }
    
    // 设置请求头
    req.Header.Set("Content-Type", "application/json")
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    return c.client.Do(req)
}

// PUT请求
func (c *HTTPClient) Put(url string, body interface{}) (*http.Response, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("序列化请求体失败: %v", err)
    }
    
    req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("创建请求失败: %v", err)
    }
    
    // 设置请求头
    req.Header.Set("Content-Type", "application/json")
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    return c.client.Do(req)
}

// DELETE请求
func (c *HTTPClient) Delete(url string) (*http.Response, error) {
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return nil, fmt.Errorf("创建请求失败: %v", err)
    }
    
    // 设置请求头
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    return c.client.Do(req)
}

// 解析JSON响应
func (c *HTTPClient) ParseJSONResponse(resp *http.Response, v interface{}) error {
    defer resp.Body.Close()
    
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(v); err != nil {
        return fmt.Errorf("解析JSON响应失败: %v", err)
    }
    
    return nil
}

// 获取客户端配置
func (c *HTTPClient) GetConfig() *HTTPClientConfig {
    return c.config
}

// 获取请求头
func (c *HTTPClient) GetHeaders() map[string]string {
    return c.headers
}

// 关闭客户端
func (c *HTTPClient) Close() error {
    if c.client.Transport != nil {
        if transport, ok := c.client.Transport.(*http.Transport); ok {
            transport.CloseIdleConnections()
        }
    }
    return nil
}

func main() {
    // 创建HTTP客户端配置
    config := &HTTPClientConfig{
        Timeout:           10 * time.Second,
        MaxIdleConns:      5,
        IdleConnTimeout:    30 * time.Second,
        TLSHandshakeTimeout: 5 * time.Second,
    }
    
    // 创建HTTP客户端
    client := NewHTTPClient(config)
    
    // 设置请求头
    client.SetHeader("User-Agent", "Go HTTP Client")
    client.SetHeader("Accept", "application/json")
    
    // 发送GET请求
    resp, err := client.Get("https://httpbin.org/get")
    if err != nil {
        log.Fatalf("GET请求失败: %v", err)
    }
    defer resp.Body.Close()
    
    fmt.Printf("GET响应状态码: %d\n", resp.StatusCode)
    
    // 发送POST请求
    data := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }
    
    resp, err = client.Post("https://httpbin.org/post", data)
    if err != nil {
        log.Fatalf("POST请求失败: %v", err)
    }
    defer resp.Body.Close()
    
    fmt.Printf("POST响应状态码: %d\n", resp.StatusCode)
    
    // 解析POST响应
    var result map[string]interface{}
    if err := client.ParseJSONResponse(resp, &result); err != nil {
        log.Fatalf("解析POST响应失败: %v", err)
    }
    
    <｜fim▁hole｜>fmt.Printf("POST响应数据: %+v\n", result)
    
    // 关闭客户端
    if err := client.Close(); err != nil {
        log.Fatalf("关闭客户端失败: %v", err)
    }
    
    fmt.Println("HTTP客户端示例完成")
}