package task

import (
	"fmt"
	"time"
	"todolist/utils"

	_ "github.com/glebarez/go-sqlite"
)

type Task struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Completed   bool      `json:"completed"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
	BeginDate   time.Time `json:"begin_date"`
	Priority    int       `json:"priority"`
	Location    string    `json:"location"`
	Label       string    `json:"label"`
	UserId      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskInterface interface {
	Complete()
	Save() error
	Delete() error
	Print()
	PrintDetails()
}

func IsTaskExist(id int64) bool {
	var count int

	row := utils.SqliteInstance.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ?", id)
	row.Scan(&count)

	return count > 0
}

func NewTask(name string) Task {
	task := Task{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return task
}

func GetTask(id int64) Task {
	var task Task

	row := utils.SqliteInstance.DB.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	row.Scan(
		&task.Id,
		&task.Name,
		&task.Description,
		&task.Completed,
		&task.EndDate,
		&task.BeginDate,
		&task.Priority,
		&task.Location,
		&task.Label,
		&task.UserId,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return task
}

func GetTasksByUserId(id int64) ([]Task, error) {
	var tasks []Task

	row, err := utils.SqliteInstance.DB.Query("SELECT * FROM tasks WHERE user_id = ?", id)

	if err != nil {
		return tasks, err
	}

	for row.Next() {
		var task Task

		row.Scan(
			&task.Id,
			&task.Name,
			&task.Description,
			&task.Completed,
			&task.EndDate,
			&task.BeginDate,
			&task.Priority,
			&task.Location,
			&task.Label,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *Task) Complete() {
	t.Completed = !t.Completed
}

func (t *Task) Save() error {
	var query string

	if t.Id == 0 {
		query = "INSERT INTO tasks (name, description, completed, end_date, begin_date, priority, location, label, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	} else {
		query = "UPDATE tasks SET name = ?, description = ?, completed = ?, end_date = ?, begin_date = ?, priority = ?, location = ?, label = ?, user_id = ?, created_at = ?, updated_at = ? WHERE id = ?"
	}

	stmt, err := utils.SqliteInstance.DB.Prepare(query)

	if err != nil {
		return err
	}

	if t.Id == 0 {
		res, err := stmt.Exec(
			t.Name,
			t.Description,
			t.Completed,
			t.EndDate,
			t.BeginDate,
			t.Priority,
			t.Location,
			t.Label,
			t.UserId,
			t.CreatedAt,
			t.UpdatedAt,
		)

		if err != nil {
			return err
		}

		t.Id, _ = res.LastInsertId()
	} else {
		t.UpdatedAt = time.Now()
		_, err := stmt.Exec(
			t.Name,
			t.Description,
			t.Completed,
			t.EndDate,
			t.BeginDate,
			t.Priority,
			t.Location,
			t.Label,
			t.UserId,
			t.CreatedAt,
			t.UpdatedAt,
			t.Id,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) Delete() error {
	stmt, err := utils.SqliteInstance.DB.Prepare("DELETE FROM tasks WHERE id = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.Id)

	if err != nil {
		return err
	}

	return nil
}

func (t *Task) PrintDetails() {
	fmt.Println("Id:          ", t.Id)
	fmt.Println("Name:        ", t.Name)
	fmt.Println("Completed:   ", t.Completed)

	if t.Description != "" {
		fmt.Println("Description: ", t.Description)
	}
	if t.Priority != 0 {
		fmt.Println("Priority:    ", t.Priority)
	}
	if t.Location != "" {
		fmt.Println("Location:    ", t.Location)
	}
	if t.Label != "" {
		fmt.Println("Label:       ", t.Label)
	}
}

func (t *Task) Print() {
	completed := " "

	if t.Completed {
		completed = "x"
	}

	fmt.Printf("[%s] [%d] %s\n", completed, t.Id, t.Name)
}
