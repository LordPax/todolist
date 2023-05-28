package task

import (
	"fmt"
	"testing"
	"time"
	"todolist/utils"

	_ "github.com/glebarez/go-sqlite"
	fakerLib "github.com/jaswdr/faker"
)

const (
	TASK_NAME = "TestTask"
)

var (
	faker = fakerLib.New()
	tasks = []Task{
		{
			Name:        fmt.Sprintf("%s1", TASK_NAME),
			Completed:   faker.Bool(),
			Description: faker.Lorem().Sentence(10),
			EndDate:     time.Now().AddDate(0, 0, faker.IntBetween(1, 10)),
		},
		{
			Name:        fmt.Sprintf("%s2", TASK_NAME),
			Completed:   faker.Bool(),
			Description: faker.Lorem().Sentence(10),
			EndDate:     time.Now().AddDate(0, 0, faker.IntBetween(1, 10)),
		},
		{
			Name:        fmt.Sprintf("%s3", TASK_NAME),
			Completed:   faker.Bool(),
			Description: faker.Lorem().Sentence(10),
			EndDate:     time.Now().AddDate(0, 0, faker.IntBetween(1, 10)),
		},
	}
)

func TestNewTask(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()
	task := NewTask(TASK_NAME)

	if task.Id != 0 {
		t.Error("Task id should be 0 but is", task.Id)
	}

	if task.Name != TASK_NAME {
		t.Error("Task name should be", TASK_NAME)
	}

	task.Save()

	row := utils.SqliteInstance.DB.QueryRow("SELECT id, name FROM tasks")

	var taskDB Task
	row.Scan(&taskDB.Id, &taskDB.Name)

	if taskDB.Id != 1 {
		t.Error("Task id should be 1 but is", taskDB.Id)
	}

	if taskDB.Name != TASK_NAME {
		t.Error("Task name should be", TASK_NAME, "but is", taskDB.Name)
	}
}

func TestGetTask(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()

	// insert tasks
	for _, task := range tasks {
		stmt, _ := utils.SqliteInstance.DB.Prepare("INSERT INTO tasks (name, completed, description, end_date) VALUES (?, ?, ?, ?)")
		stmt.Exec(task.Name, task.Completed, task.Description, task.EndDate)
	}

	for i, task := range tasks {
		taskDB := GetTask(int64(i + 1))

		if taskDB.Id != int64(i+1) {
			t.Error("Task id should be", i+1)
		}
		if taskDB.Name != task.Name {
			t.Error("Task name should be", task.Name)
		}
		if taskDB.Completed != task.Completed {
			t.Error("Task completed should be", task.Completed)
		}
		if taskDB.Description != task.Description {
			t.Error("Task description should be", task.Description)
		}
		if !taskDB.EndDate.Equal(task.EndDate) {
			t.Error("Task end date should be", task.EndDate, "but is", taskDB.EndDate)
		}
	}
}

func TestGetTasksByUserId(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()
	var userId int64 = 1

	// insert tasks
	for _, task := range tasks {
		task.UserId = userId
		task.Save()
	}

	tasksDB, _ := GetTasksByUserId(userId)

	if len(tasksDB) != len(tasks) {
		t.Error("Tasks should have", len(tasks), "tasks but has", len(tasksDB))
	}

	for i, taskDB := range tasksDB {
		if taskDB.Id != int64(i+1) {
			t.Error("Task id should be", i+1)
		}
		if taskDB.Name != tasks[i].Name {
			t.Error("Task name should be", tasks[i].Name)
		}
		if taskDB.Completed != tasks[i].Completed {
			t.Error("Task completed should be", tasks[i].Completed)
		}
		if taskDB.Description != tasks[i].Description {
			t.Error("Task description should be", tasks[i].Description)
		}
		if !taskDB.EndDate.Equal(tasks[i].EndDate) {
			t.Error("Task end date should be", tasks[i].EndDate, "but is", taskDB.EndDate)
		}
	}
}

func TestSave(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()

	// insert tasks
	for i, task := range tasks {
		if task.Id != 0 {
			t.Error("Task id should be 0 but is", task.Id)
		}

		task.Save()

		if task.Id != int64(i+1) {
			t.Error("Task id should be", i+1)
		}
	}

	rows, _ := utils.SqliteInstance.DB.Query("SELECT id, name, completed, description, end_date FROM tasks")

	var tasksDB []Task
	for rows.Next() {
		var taskDB Task
		rows.Scan(
			&taskDB.Id,
			&taskDB.Name,
			&taskDB.Completed,
			&taskDB.Description,
			&taskDB.EndDate,
		)
		tasksDB = append(tasksDB, taskDB)
	}

	for i, taskDB := range tasksDB {
		if taskDB.Id != int64(i+1) {
			t.Error("Task id should be", i+1)
		}
		if taskDB.Name != tasks[i].Name {
			t.Error("Task name should be", tasks[i].Name)
		}
		if taskDB.Completed != tasks[i].Completed {
			t.Error("Task completed should be", tasks[i].Completed)
		}
		if taskDB.Description != tasks[i].Description {
			t.Error("Task description should be", tasks[i].Description)
		}
		if !taskDB.EndDate.Equal(tasks[i].EndDate) {
			t.Error("Task end date should be", tasks[i].EndDate, "but is", taskDB.EndDate)
		}
	}

	// update tasks
	for _, taskDB := range tasksDB {
		// taskDB.Description = taskDB.Description + "_updated"
		taskDB.Description = fmt.Sprintf("%s_updated", taskDB.Description)
		taskDB.Save()
	}

	rows, _ = utils.SqliteInstance.DB.Query("SELECT id, name, completed, description, end_date FROM tasks")

	var tasksDBUpdate []Task
	for rows.Next() {
		var taskDB Task
		rows.Scan(
			&taskDB.Id,
			&taskDB.Name,
			&taskDB.Completed,
			&taskDB.Description,
			&taskDB.EndDate,
		)
		tasksDBUpdate = append(tasksDBUpdate, taskDB)
	}

	for i, taskDBUpdate := range tasksDBUpdate {
		if taskDBUpdate.Id != tasksDB[i].Id {
			t.Error("Task id should be", tasksDB[i].Id)
		}
		if taskDBUpdate.Name != tasksDB[i].Name {
			t.Error("Task name should be", tasksDB[i].Name)
		}
		if taskDBUpdate.Completed != tasksDB[i].Completed {
			t.Error("Task completed should be", tasksDB[i].Completed)
		}
		if taskDBUpdate.Description != fmt.Sprintf("%s_updated", tasksDB[i].Description) {
			t.Error("Task description should be", tasksDB[i].Description)
		}
		if !taskDBUpdate.EndDate.Equal(tasksDB[i].EndDate) {
			t.Error("Task end date should be", tasksDB[i].EndDate)
		}
	}
}

func TestDelete(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()

	for _, task := range tasks {
		stmt, _ := utils.SqliteInstance.DB.Prepare("INSERT INTO tasks (name, completed, description, end_date) VALUES (?, ?, ?, ?)")
		stmt.Exec(task.Name, task.Completed, task.Description, task.EndDate)
	}

	var tasksDB []Task
	rows, _ := utils.SqliteInstance.DB.Query("SELECT id, name, completed, description, end_date FROM tasks")
	for rows.Next() {
		var taskDB Task
		rows.Scan(
			&taskDB.Id,
			&taskDB.Name,
			&taskDB.Completed,
			&taskDB.Description,
			&taskDB.EndDate,
		)
		tasksDB = append(tasksDB, taskDB)
	}

	if len(tasksDB) != len(tasks) {
		t.Error("Tasks should have", len(tasks), "tasks but has", len(tasksDB))
	}

	for _, taskDB := range tasksDB {
		taskDB.Delete()
	}

	row := utils.SqliteInstance.DB.QueryRow("SELECT COUNT(*) FROM tasks")

	var count int
	row.Scan(&count)

	if count != 0 {
		t.Error("Tasks should have 0 tasks but has", count)
	}
}

func TestComplete(t *testing.T) {
	task := NewTask(TASK_NAME)
	task.Completed = false
	task.Complete()

	if !task.Completed {
		t.Error("Task should not be complete")
	}
}

func TestIsTaskExist(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()

	for _, task := range tasks {
		stmt, _ := utils.SqliteInstance.DB.Prepare("INSERT INTO tasks (name, completed, description, end_date) VALUES (?, ?, ?, ?)")
		stmt.Exec(task.Name, task.Completed, task.Description, task.EndDate)
	}

	var count int
	row := utils.SqliteInstance.DB.QueryRow("SELECT COUNT(*) FROM tasks")
	row.Scan(&count)

	if count != len(tasks) {
		t.Error("It should have", len(tasks), "tasks but has", count)
	}

	for i := 0; i < len(tasks); i++ {
		if !IsTaskExist(int64(i + 1)) {
			t.Error("Task should exist")
		}
	}
}
