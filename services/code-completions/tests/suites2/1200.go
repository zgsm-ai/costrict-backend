package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)
<｜fim▁hole｜>
// Task 表示一个任务
type Task struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	Result    string     `json:"result"`
	Error     string     `json:"error"`
}

// TaskHandler 处理任务接口
type TaskHandler interface {
	Handle(ctx context.Context, task *Task) error
	Name() string
}

// TaskQueue 任务队列
type TaskQueue struct {
	tasks    []*Task
	handlers map[string]TaskHandler
	mutex    sync.Mutex
	cond     *sync.Cond
	running  bool
	workers  int
}

// NewTaskQueue 创建新的任务队列
func NewTaskQueue(workers int) *TaskQueue {
	tq := &TaskQueue{
		tasks:    make([]*Task, 0),
		handlers: make(map[string]TaskHandler),
		workers:  workers,
	}
	tq.cond = sync.NewCond(&tq.mutex)
	return tq
}

// RegisterHandler 注册任务处理器
func (tq *TaskQueue) RegisterHandler(name string, handler TaskHandler) {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()
	tq.handlers[name] = handler
}

// AddTask 添加任务
func (tq *TaskQueue) AddTask(task *Task) error {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()

	task.CreatedAt = time.Now()
	task.Status = TaskStatusPending
	tq.tasks = append(tq.tasks, task)
	tq.cond.Signal()

	return nil
}

// getNextTask 获取下一个任务
func (tq *TaskQueue) getNextTask() *Task {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()

	for len(tq.tasks) == 0 {
		tq.cond.Wait()
	}

	// 找到第一个待处理的任务
	for i, task := range tq.tasks {
		if task.Status == TaskStatusPending {
			// 从队列中移除任务
			tq.tasks = append(tq.tasks[:i], tq.tasks[i+1:]...)
			return task
		}
	}

	return nil
}

// processTask 处理任务
func (tq *TaskQueue) processTask(ctx context.Context, task *Task) {
	defer func() {
		if r := recover(); r != nil {
			task.Status = TaskStatusFailed
			task.Error = fmt.Sprintf("panic: %v", r)
		}
	}()

	task.Status = TaskStatusRunning

	handler, ok := tq.handlers[task.Name]
	if !ok {
		task.Status = TaskStatusFailed
		task.Error = fmt.Sprintf("未找到处理器: %s", task.Name)
		return
	}

	err := handler.Handle(ctx, task)
	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err.Error()
		return
	}

	task.Status = TaskStatusCompleted
	task.Result = "任务处理完成"
}

// worker 工作线程
func (tq *TaskQueue) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := tq.getNextTask()
			if task != nil {
				tq.processTask(ctx, task)
			}
		}
	}
}

// Start 启动任务队列
func (tq *TaskQueue) Start(ctx context.Context) {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()

	if tq.running {
		return
	}

	tq.running = true

	// 启动工作线程
	for i := 0; i < tq.workers; i++ {
		go tq.worker(ctx)
	}
}

// SimpleTaskHandler 简单任务处理器
type SimpleTaskHandler struct{}

func (h *SimpleTaskHandler) Handle(ctx context.Context, task *Task) error {
	// 模拟任务处理
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (h *SimpleTaskHandler) Name() string {
	return "simple"
}

func main() {
	// 创建任务队列
	queue := NewTaskQueue(3)

	// 注册任务处理器
	queue.RegisterHandler("simple", &SimpleTaskHandler{})

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动任务队列
	go queue.Start(ctx)

	// 创建并添加任务
	for i := 1; i <= 5; i++ {
		task := &Task{
			ID:   fmt.Sprintf("task-%d", i),
			Name: "simple",
		}
		queue.AddTask(task)
	}

	// 等待任务完成
	time.Sleep(2 * time.Second)

	// 停止任务队列
	cancel()

	fmt.Println("任务队列示例完成")
}
