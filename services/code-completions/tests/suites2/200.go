package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

type User struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type UserCache struct {
    users map[int]*User
    mu    sync.RWMutex
}

func NewUserCache() *UserCache {
    return &UserCache{
        users: make(map[int]*User),
    }
}

func (uc *UserCache) AddUser(user *User) {
    uc.mu.Lock()
    defer uc.mu.Unlock()
    
    uc.users[user.ID] = user
}

func (uc *UserCache) GetUser(id int) (*User, bool) {
    uc.mu.RLock()
    defer uc.mu.RUnlock()
    
    user, exists := uc.users[id]
    <｜fim▁hole｜>
}

func (uc *UserCache) RemoveUser(id int) {
    uc.mu.Lock()
    defer uc.mu.Unlock()
    
    delete(uc.users, id)
}

func (uc *UserCache) GetAllUsers() []*User {
    uc.mu.RLock()
    defer uc.mu.RUnlock()
    
    var users []*User
    for _, user := range uc.users {
        users = append(users, user)
    }
    return users
}

func userHandler(w http.ResponseWriter, r *http.Request, cache *UserCache) {
    switch r.Method {
    case http.MethodGet:
        users := cache.GetAllUsers()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    case http.MethodPost:
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        user.CreatedAt = time.Now()
        cache.AddUser(&user)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func main() {
    cache := NewUserCache()
    
    // Add some initial users
    cache.AddUser(&User{
        ID:   1,
        Name: "Alice",
        Email: "alice@example.com",
        CreatedAt: time.Now(),
    })
    
    cache.AddUser(&User{
        ID:   2,
        Name: "Bob",
        Email: "bob@example.com",
        CreatedAt: time.Now(),
    })
    
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        userHandler(w, r, cache)
    })
    
    port := "8080"
    fmt.Printf("Server starting on port %s...\n", port)
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}