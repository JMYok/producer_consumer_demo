package model

// Task 任务
type Task struct {
	ID int64
}

// 消费任务
func (t *Task) Run() {
	OutObj.Println(t.ID)
}
