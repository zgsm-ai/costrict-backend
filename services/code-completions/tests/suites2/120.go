package main

import (
    "fmt"
    "time"
)

type Task struct {
    ID        int
    Name      string
    Completed bool
    CreatedAt time.Time
}

func NewTask(id int, name string) *Task {
    return &Task{
        ID:        id,
        Name:      name,
        Completed: false,
        CreatedAt: time.Now(),
    }
}

func (t *Task) Complete() {
    <｜fim▁hole｜>
}

func (t *Task) String() string {
    status := "Pending"
    if t.Completed {
        status = "Completed"
    }
    return fmt.Sprintf("Task #%d: %s (%s)", t.ID, t.Name, status)
}

func main() {
    task := NewTask(1, "Learn Go")
    fmt.Println(task)
    
    task.Complete()
    fmt.Println(task)
}