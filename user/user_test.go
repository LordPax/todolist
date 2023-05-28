package user

import (
	"fmt"
	"testing"
	"time"
	taskLib "todolist/task"
	"todolist/utils"

	fakerLib "github.com/jaswdr/faker"
)

var (
	faker = fakerLib.New()
	tasks = []taskLib.Task{
		{
			Name:        faker.Lorem().Word(),
			Completed:   faker.Bool(),
			Description: faker.Lorem().Sentence(10),
			EndDate:     time.Now().AddDate(0, 0, faker.IntBetween(1, 10)),
		},
	}
	users = []User{
		{
			Firstname: faker.Person().Name(),
			Lastname:  faker.Person().Name(),
			Email:     faker.Internet().Email(),
			Birthdate: faker.Time().Time(time.Now()),
			Password:  faker.Internet().Password(),
			Tasks:     tasks,
		},
		{
			Firstname: faker.Person().Name(),
			Lastname:  faker.Person().Name(),
			Email:     faker.Internet().Email(),
			Birthdate: faker.Time().Time(time.Now()),
			Password:  faker.Internet().Password(),
			Tasks:     tasks,
		},
		{
			Firstname: faker.Person().Name(),
			Lastname:  faker.Person().Name(),
			Email:     faker.Internet().Email(),
			Birthdate: faker.Time().Time(time.Now()),
			Password:  faker.Internet().Password(),
			Tasks:     tasks,
		},
	}
)

func TestNewUser(t *testing.T) {
	user := NewUser(users[0].Firstname, users[0].Lastname, users[0].Email, users[0].Tasks)

	if user.Firstname != users[0].Firstname {
		t.Error("Firstname should be", users[0].Firstname, "but got", user.Firstname)
	}

	if user.Lastname != users[0].Lastname {
		t.Error("Lastname should be", users[0].Lastname, "but got", user.Lastname)
	}

	if user.Email != users[0].Email {
		t.Error("Email should be", users[0].Email, "but got", user.Email)
	}

	if len(user.Tasks) != len(users[0].Tasks) {
		t.Error("Tasks lenght should be", len(users[0].Tasks), "but got", len(user.Tasks))
	}

	for i, task := range user.Tasks {
		if task.Name != users[0].Tasks[i].Name {
			t.Error("Task name should be", users[0].Tasks[i].Name, "but got", task.Name)
		}
		if task.Description != users[0].Tasks[i].Description {
			t.Error("Task description should be", users[0].Tasks[i].Description, "but got", task.Description)
		}
	}
}

func TestSaveUser(t *testing.T) {
	utils.SqliteInstance, _ = utils.ConnectDB(true)
	defer utils.SqliteInstance.Close()

	// insert lists
	for i, user := range users {
		if user.Id != 0 {
			t.Error("User id should be 0 but is", user.Id)
		}

		user.Save()

		if user.Id != int64(i+1) {
			t.Error("User id should be", i+1, "but is", user.Id)
		}
	}

	usersDB, _ := GetUsers()

	for i, userDB := range usersDB {
		if userDB.Id != int64(i+1) {
			t.Error("First id should be", i+1, "but got", userDB.Id)
		}
		if userDB.Firstname != users[i].Firstname {
			t.Error("Firstname should be", users[i].Firstname, "but got", userDB.Firstname)
		}
		if userDB.Lastname != users[i].Lastname {
			t.Error("Lastname should be", users[i].Lastname, "but got", userDB.Lastname)
		}
		if userDB.Email != users[i].Email {
			t.Error("Email should be", users[i].Email, "but got", userDB.Email)
		}
		if !userDB.Birthdate.Equal(users[i].Birthdate) {
			t.Error("Birthdate should be", users[i].Birthdate, "but got", userDB.Birthdate)
		}
		if userDB.Password != users[i].Password {
			t.Error("Password should be", users[i].Password, "but got", userDB.Password)
		}
		if len(userDB.Tasks) != len(tasks) {
			t.Error("Tasks lenght should be", len(tasks), "but got", len(userDB.Tasks))
		}
		for i, task := range userDB.Tasks {
			if task.Name != tasks[i].Name {
				t.Error("Task name should be", tasks[i].Name, "but got", task.Name)
			}
		}
	}

	// update list
	for _, userDB := range usersDB {
		userDB.Firstname = fmt.Sprintf("%s_updated", userDB.Firstname)
		userDB.Save()
	}

	usersDBUpdate, _ := GetUsers()

	for i, userDBUpdate := range usersDBUpdate {
		if userDBUpdate.Firstname != fmt.Sprintf("%s_updated", users[i].Firstname) {
			t.Error("Firstname should be", fmt.Sprintf("%s_updated", users[i].Firstname), "but got", userDBUpdate.Firstname)
		}
		if userDBUpdate.Lastname != users[i].Lastname {
			t.Error("Lastname should be", users[i].Lastname, "but got", userDBUpdate.Lastname)
		}
		if userDBUpdate.Email != users[i].Email {
			t.Error("Email should be", users[i].Email, "but got", userDBUpdate.Email)
		}
		if !userDBUpdate.Birthdate.Equal(users[i].Birthdate) {
			t.Error("Birthdate should be", users[i].Birthdate, "but got", userDBUpdate.Birthdate)
		}
		if userDBUpdate.Password != users[i].Password {
			t.Error("Password should be", users[i].Password, "but got", userDBUpdate.Password)
		}
		if len(userDBUpdate.Tasks) != len(tasks) {
			t.Error("Tasks length should be", len(tasks), "but got", len(userDBUpdate.Tasks))
		}
		for i, task := range userDBUpdate.Tasks {
			if task.Name != tasks[i].Name {
				t.Error("Task name should be", tasks[i].Name, "but got", task.Name)
			}
		}
	}
}

// func mockSendEmail(email string, subject string, body string) error {
// 	return nil
// }

func TestAddTask(t *testing.T) {
	// origSendEmail := services.SendEmail
	// defer func() { services.SendEmail = origSendEmail }()
	// services.SendEmail = mockSendEmail

	user := NewUser(faker.Person().Name(), faker.Person().Name(), faker.Internet().Email(), nil)

	for i := 0; i < 5; i++ {
		err := user.AddTask(tasks[0])
		if err != nil {
			t.Error("Should not return an error")
		}
		if len(user.Tasks) != i+1 {
			t.Error("Tasks length should be", i+1, "but got", len(user.Tasks))
		}
		if user.Tasks[i].Name != tasks[0].Name {
			t.Error("Task name should be", tasks[0].Name, "but got", user.Tasks[i].Name)
		}
		if user.Tasks[i].Description != tasks[0].Description {
			t.Error("Task description should be", tasks[0].Description, "but got", user.Tasks[i].Description)
		}
	}

	user = NewUser(faker.Person().Name(), faker.Person().Name(), faker.Internet().Email(), nil)

	for i := 0; i < 10; i++ {
		user.Tasks = append(user.Tasks, tasks[0])
	}

	err := user.AddTask(tasks[0])

	if err == nil {
		t.Error("Should return an error")
	}
	if len(user.Tasks) > 10 {
		t.Error("Tasks length should be lower than 10 but got", len(user.Tasks))
	}
}
