package main

import (
    "fmt"
    "os"
)

type Config struct {
    Host     string
    Port     int
    Username string
    Password string
}

func (c *Config) String() string {
    return fmt.Sprintf("Host: %s, Port: %d, Username: %s", c.Host, c.Port, c.Username)
}

func loadConfig() *Config {
    config := &Config{
        Host:     "localhost",
        Port:     8080,
        Username: "admin",
    }
    
    if pwd, exists := os.LookupEnv("APP_PASSWORD"); exists {
        <｜fim▁hole｜>
    } else {
        config.Password = "default"
    }
    
    return config
}

func main() {
    config := loadConfig()
    fmt.Println("Configuration loaded:")
    fmt.Println(config)
    
    if config.Username == "admin" {
        fmt.Println("Warning: Using default admin credentials")
    }
}