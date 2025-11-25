package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

type Server struct {
    clients    map[net.Addr]*Client
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
    ctx        context.Context
    cancel     context.CancelFunc
}

type Client struct {
    conn     net.Conn
    send     chan []byte
    server   *Server
    username string
}

type Message struct {
    Username string    `json:"username"`
    Content string    `json:"content"`
    Time    time.Time `json:"time"`
}

func NewServer() *Server {
    ctx, cancel := context.WithCancel(context.Background())
    return &Server{
        clients:    make(map[net.Addr]*Client),
        broadcast:  make(chan *Message, 100),
        register:   make(chan *Client, 100),
        unregister: make(chan *Client, 100),
        mu:         sync.Mutex{},
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (s *Server) Run() {
    log.Println("Server started")
    
    for {
        select {
        case <-s.ctx.Done():
            log.Println("Server shutting down...")
            return
            
        case client := <-s.register:
            s.mu.Lock()
            s.clients[client.conn.RemoteAddr()] = client
            s.mu.Unlock()
            
            log.Printf("Client %s connected", client.conn.RemoteAddr())
            
        case client := <-s.unregister:
            s.mu.Lock()
            if _, ok := s.clients[client.conn.RemoteAddr()]; ok {
                delete(s.clients, client.conn.RemoteAddr())
                close(client.send)
            }
            s.mu.Unlock()
            
            log.Printf("Client %s disconnected", client.conn.RemoteAddr())
            
        case message := <-s.broadcast:
            s.mu.Lock()
            for _, client := range s.clients {
                select {
                case client.send <- []byte(fmt.Sprintf("[%s] %s: %s", 
                    message.Time.Format("15:04:05"), message.Username, message.Content)):
                default:
                    close(client.send)
                    delete(s.clients, client.conn.RemoteAddr())
                }
            }
            s.mu.Unlock()
        }
    }
}

func (s *Server) Stop() {
    s.cancel()
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for addr, client := range s.clients {
        client.conn.Close()
        close(client.send)
        delete(s.clients, addr)
    }
    
    log.Println("Server stopped")
}

func (s *Server) Listen(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    defer listener.Close()
    
    log.Printf("Listening on %s", addr)
    
    go s.Run()
    
    for {
        select {
        case <-s.ctx.Done():
            return nil
            
        default:
            conn, err := listener.Accept()
            if err != nil {
                log.Printf("Accept error: %v", err)
                continue
            }
            
            client := &Client{
                conn:   conn,
                send:   make(chan []byte, 256),
                server: s,
            }
            
            go client.readPump()
            go client.writePump()
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.server.unregister <- c
        c.conn.Close()
    }()
    
    buf := make([]byte, 1024)
    for {
        select {
        case <-c.server.ctx.Done():
            return
            
        default:
            n, err := c.conn.Read(buf)
            if err != nil {
                log.Printf("Read error: %v", err)
                return
            }
            
            if n == 0 {
                return
            }
            
            message := string(buf[:n])
            
            // Handle special commands
            if message == "/quit" {
                return
            }
            
            // Handle username setting
            if c.username == "" && len(message) >= 5 && message[:5] == "/nick" {
                <｜fim▁hole｜>
                continue
            }
            
            // Broadcast regular messages
            if c.username != "" {
                c.server.broadcast <- &Message{
                    Username: c.username,
                    Content:  message,
                    Time:     time.Now(),
                }
            }
        }
    }
}

func (c *Client) writePump() {
    defer func() {
        c.conn.Close()
    }()
    
    for {
        select {
        case <-c.server.ctx.Done():
            return
            
        case message, ok := <-c.send:
            if !ok {
                c.conn.Write([]byte("Server closed connection"))
                return
            }
            
            _, err := c.conn.Write(message)
            if err != nil {
                log.Printf("Write error: %v", err)
                return
            }
        }
    }
}

func main() {
    server := NewServer()
    defer server.Stop()
    
    go func() {
        time.Sleep(30 * time.Second)
        log.Println("Server stopping after 30 seconds for demo")
        server.Stop()
    }()
    
    err := server.Listen(":8080")
    if err != nil {
        log.Fatalf("Server error: %v", err)
    }
}