package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"` // 1-5, 5 is highest
	Status      string    `json:"status"`  // "todo", "in_progress", "done"
	DueDate     time.Time `json:"due_date"`
	Tags        []string  `json:"tags"`
	Assignee    string    `json:"assignee"`
}

type TaskManager struct {
	tasks      map[int]*Task
	categories map[string][]int
	mutex      sync.RWMutex
}

type TaskHeap []*Task

func (h TaskHeap) Len() int           { return len(h) }
func (h TaskHeap) Less(i, j int) bool { 
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	return h[i].DueDate.Before(h[j].DueDate)
}
func (h TaskHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TaskHeap) Push(x interface{}) {
	*h = append(*h, x.(*Task))
}

func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return item
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:      make(map[int]*Task),
		categories: make(map[string][]int),
		mutex:      sync.RWMutex{},
	}
}

func (tm *TaskManager) AddTask(task *Task) int {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	task.ID = len(tm.tasks)
	tm.tasks[task.ID] = task
	
	// Add to categories
	for _, tag := range task.Tags {
		tm.categories[tag] = append(tm.categories[tag], task.ID)
	}
	
	log.Printf("Task added: %s (ID: %d)", task.Title, task.ID)
	return task.ID
}

func (tm *TaskManager) GetTask(id int) (*Task, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	task, exists := tm.tasks[id]
	return task, exists
}

func (tm *TaskManager) UpdateTask(id int, updates map[string]interface{}) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	task, exists := tm.tasks[id]
	if !exists {
		return false
	}
	
	// Remove from old categories
	for _, tag := range task.Tags {
		for i, taskID := range tm.categories[tag] {
			if taskID == id {
				tm.categories[tag] = append(tm.categories[tag][:i], tm.categories[tag][i+1:]...)
				break
			}
		}
	}
	
	// Apply updates
	for key, value := range updates {
		switch key {
		case "title":
			if title, ok := value.(string); ok {
				task.Title = title
			}
		case "description":
			if description, ok := value.(string); ok {
				task.Description = description
			}
		case "priority":
			if priority, ok := value.(int); ok && priority >=1 && priority <= 5 {
				task.Priority = priority
			}
		case "status":
			if status, ok := value.(string); ok {
				if status == "todo" || status == "in_progress" || status == "done" {
					task.Status = status
				}
			}
		case "due_date":
			if dueDate, ok := value.(time.Time); ok {
				task.DueDate = dueDate
			}
		case "tags":
			if tags, ok := value.([]string); ok {
				task.Tags = tags
			}
		case "assignee":
			if assignee, ok := value.(string); ok {
				task.Assignee = assignee
			}
		}
	}
	
	// Add to new categories
	for _, tag := range task.Tags {
		tm.categories[tag] = append(tm.categories[tag], id)
	}
	
	log.Printf("Task updated: %s (ID: %d)", task.Title, id)
	return true
}

func (tm *TaskManager) DeleteTask(id int) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	task, exists := tm.tasks[id]
	if !exists {
		return false
	}
	
	// Remove from categories
	for _, tag := range task.Tags {
		for i, taskID := range tm.categories[tag] {
			if taskID == id {
				tm.categories[tag] = append(tm.categories[tag][:i], tm.categories[tag][i+1:]...)
				break
			}
		}
	}
	
	delete(tm.tasks, id)
	log.Printf("Task deleted: %s (ID: %d)", task.Title, id)
	return true
}

func (tm *TaskManager) GetTasksByTag(tag string) []*Task {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	var tasks []*Task
	for _, taskID := range tm.categories[tag] {
		if task, exists := tm.tasks[taskID]; exists {
			tasks = append(tasks, task)
		}
	}
	
	return tasks
}

func (tm *TaskManager) GetTasksByStatus(status string) []*Task {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	var tasks []*Task
	for _, task := range tm.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	
	return tasks
}

func (tm *TaskManager) GetAllTasks() []*Task {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	
	return tasks
}

func (tm *TaskManager) GetHighPriorityTasks() []*Task {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	var tasks []*Task
	for _, task := range tm.tasks {
		if task.Priority >= 4 && task.Status != "done" {
			tasks = append(tasks, task)
		}
	}
	
	return tasks
}

func (tm *TaskManager) SortTasksByPriority(tasks []*Task) []*Task {
	h := &TaskHeap{}
	heap.Init(h)
	
	for _, task := range tasks {
		heap.Push(h, task)
	}
	
	sorted := make([]*Task, len(tasks))
	for i := range sorted {
		sorted[i] = heap.Pop(h).(*Task)
	}
	
	return sorted
}

func (tm *TaskManager) GetTaskStatistics() map[string]interface{} {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	stats := make(map[string]interface{})
	stats["total_tasks"] = len(tm.tasks)
	
	// Count by status
	statusCount := make(map[string]int)
	for _, task := range tm.tasks {
		statusCount[task.Status]++
	}
	stats["by_status"] = statusCount
	
	// Count by priority
	priorityCount := make(map[int]int)
	for _, task := range tm.tasks {
		priorityCount[task.Priority]++
	}
	stats["by_priority"] = priorityCount
	
	// Count overdue tasks
	now := time.Now()
	overdueCount := 0
	for _, task := range tm.tasks {
		if task.Status != "done" && task.DueDate.Before(now) {
			overdueCount++
		}
	}
	stats["overdue_tasks"] = overdueCount
	
	// Get most used tags
	tagCount := make(map[string]int)
	for tag, taskIDs := range tm.categories {
		tagCount[tag] = len(taskIDs)
	}
	
	type tagPair struct {
		tag   string
		count int
	}
	
	var tagPairs []tagPair
	for tag, count := range tagCount {
		tagPairs = append(tagPairs, tagPair{tag, count})
	}
	
	sort.Slice(tagPairs, func(i, j int) bool {
		return tagPairs[i].count > tagPairs[j].count
	})
	
	var topTags []string
	for i := 0; i < 5 && i < len(tagPairs); i++ {
		topTags = append(topTags, tagPairs[i].tag)
	}
	stats["top_tags"] = topTags
	
	return stats
}

func (tm *TaskManager) SaveToFile(filename string) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return err
	}
	
	<｜fim▁hole｜>
	
	return nil
}

func (tm *TaskManager) SearchTasks(query string) []*Task {
	query = strings.ToLower(query)
	
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	var results []*Task
	for _, task := range tm.tasks {
		if strings.Contains(strings.ToLower(task.Title), query) ||
		   strings.Contains(strings.ToLower(task.Description), query) {
			results = append(results, task)
		}
	}
	
	return results
}

func generateSampleTasks() []*Task {
	now := time.Now()
	
	return []*Task{
		{
			Title:       "完成项目提案",
			Description: "为新产品准备项目提案文档",
			Priority:    5,
			Status:      "in_progress",
			DueDate:     now.AddDate(0, 0, 3),
			Tags:        []string{"工作", "重要", "文档"},
			Assignee:    "张三",
		},
		{
			Title:       "代码审查",
			Description: "审查团队成员提交的代码",
			Priority:    3,
			Status:      "todo",
			DueDate:     now.AddDate(0, 0, 1),
			Tags:        []string{"开发", "代码"},
			Assignee:    "李四",
		},
		{
			Title:       "客户会议",
			Description: "与客户讨论项目需求",
			Priority:    4,
			Status:      "todo",
			DueDate:     now.AddDate(0, 0, 2),
			Tags:        []string{"会议", "客户"},
			Assignee:    "张三",
		},
		{
			Title:       "修复Bug",
			Description: "修复用户报告的登录Bug",
			Priority:    4,
			Status:      "in_progress",
			DueDate:     now.AddDate(0, 0, 1),
			Tags:        []string{"开发", "Bug", "紧急"},
			Assignee:    "王五",
		},
	}
}

func main() {
	fmt.Println("Task Manager Demo")
	fmt.Println("================")
	
	// Create a new task manager
	manager := NewTaskManager()
	
	// Add sample tasks
	tasks := generateSampleTasks()
	taskIDs := make([]int, 0)
	for _, task := range tasks {
		id := manager.AddTask(task)
		taskIDs = append(taskIDs, id)
	}
	
	fmt.Printf("Added %d tasks\n", len(tasks))
	
	// Display statistics
	stats := manager.GetTaskStatistics()
	fmt.Printf("\nTask Statistics:\n")
	fmt.Printf("Total tasks: %d\n", stats["total_tasks"])
	fmt.Printf("By status: %v\n", stats["by_status"])
	fmt.Printf("By priority: %v\n", stats["by_priority"])
	fmt.Printf("Overdue tasks: %d\n", stats["overdue_tasks"])
	fmt.Printf("Top tags: %v\n", stats["top_tags"])
	
	// Display all tasks
	fmt.Println("\nAll Tasks:")
	allTasks := manager.GetAllTasks()
	for _, task := range allTasks {
		fmt.Printf("- %s (Priority: %d, Status: %s)\n", task.Title, task.Priority, task.Status)
	}
	
	// Display high priority tasks
	fmt.Println("\nHigh Priority Tasks:")
	highPriorityTasks := manager.GetHighPriorityTasks()
	for _, task := range highPriorityTasks {
		fmt.Printf("- %s (Priority: %d, Due: %s)\n", task.Title, task.Priority, task.DueDate.Format("2006-01-02"))
	}
	
	// Sort tasks by priority
	fmt.Println("\nTasks Sorted by Priority:")
	sortedTasks := manager.SortTasksByPriority(allTasks)
	for _, task := range sortedTasks {
		fmt.Printf("- %s (Priority: %d)\n", task.Title, task.Priority)
	}
	
	// Search tasks
	fmt.Println("\nSearch Results for 'Bug':")
	searchResults := manager.SearchTasks("Bug")
	for _, task := range searchResults {
		fmt.Printf("- %s\n", task.Title)
	}
	
	// Update a task
	if len(taskIDs) > 0 {
		updates := map[string]interface{}{
			"status": "done",
		}
		if manager.UpdateTask(taskIDs[0], updates) {
			fmt.Printf("\nTask %d status updated to 'done'\n", taskIDs[0])
		}
	}
	
	// Save tasks to file
	if err := manager.SaveToFile("tasks.json"); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
	} else {
		fmt.Println("\nTasks saved to tasks.json")
	}
	
	fmt.Println("\nTask manager demo completed")
}