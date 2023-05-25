package user

import (
	"net/mail"
	"time"
	taskLib "todolist/task"
	"todolist/utils"
)

type User struct {
	Id        int64          `json:"id"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Email     string         `json:"email"`
	Birthdate time.Time      `json:"birthdate"`
	Password  string         `json:"password"`
	Tasks     []taskLib.Task `json:"task"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UserInterface interface {
	IsValid() bool
	GetAge() int
	AddTask(task taskLib.Task)
	GetTask(index int64) taskLib.Task
	GetTasks() []taskLib.Task
	DeleteTask(index int64) error
	CompleteTask(index int64)
	Save() error
}

func IsListExist(name string) bool {
	var count int64

	row := utils.SqliteInstance.DB.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", name)
	row.Scan(&count)

	return count > 0
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func NewUser(firstname string, lastname string, email string, tasks []taskLib.Task) User {
	if tasks == nil {
		tasks = []taskLib.Task{}
	}

	return User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Tasks:     tasks,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetUser(id int64) (User, error) {
	var user User

	row := utils.SqliteInstance.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)

	row.Scan(
		&user.Id,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Birthdate,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	tasks, err := taskLib.GetTasksByUserId(id)

	if err != nil {
		return user, err
	}

	user.Tasks = tasks

	return user, nil

}

func GetUsers() ([]User, error) {
	var users []User

	rows, err := utils.SqliteInstance.DB.Query("SELECT id FROM users")

	if err != nil {
		return users, err
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		user, err := GetUser(id)

		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *User) IsValid() bool {
	if u.Firstname == "" ||
		u.Lastname == "" ||
		u.Email == "" ||
		!isValidEmail(u.Email) ||
		!isValidDate(u.Birthdate.Format("2006-01-02")) ||
		u.GetAge() < 13 {
		return false
	}

	return true
}

func (u *User) GetAge() int {
	return time.Now().Year() - u.Birthdate.Year()
}

func (u *User) AddTask(task taskLib.Task) {
	u.Tasks = append(u.Tasks, task)
}

func (u *User) GetTask(index int64) taskLib.Task {
	return u.Tasks[index]
}

func (u *User) GetTasks() []taskLib.Task {
	return u.Tasks
}

func (u *User) DeleteTask(index int64) error {
	err := u.Tasks[index].Delete()

	if err != nil {
		return err
	}

	u.Tasks = append(u.Tasks[:index], u.Tasks[index+1:]...)

	return nil
}

func (u *User) CompleteTask(index int64) {
	u.Tasks[index].Complete()
}

func (u *User) Save() error {
	var query string

	if u.Id == 0 {
		query = "INSERT INTO users (firstname, lastname, email, birthdate, password, created_at, updated_at) VALUES (?, ?, ?, ?)"
	} else {
		query = "UPDATE users SET firstname = ?, lastname = ?, email = ?, birthdate = ?, password = ?, created_at = ?, updated_at = ? WHERE id = ?"
	}

	stmt, err := utils.SqliteInstance.DB.Prepare(query)

	if err != nil {
		return err
	}

	if u.Id == 0 {
		res, err := stmt.Exec(u.Firstname, u.Lastname, u.Email, u.Birthdate, u.Password, u.CreatedAt, u.UpdatedAt)

		if err != nil {
			return err
		}

		u.Id, _ = res.LastInsertId()
	} else {
		u.UpdatedAt = time.Now()
		_, err := stmt.Exec(u.Firstname, u.Lastname, u.Email, u.Birthdate, u.Password, u.CreatedAt, u.UpdatedAt, u.Id)

		if err != nil {
			return err
		}
	}

	for _, task := range u.Tasks {
		task.UserId = u.Id
		task.Save()
	}

	return nil
}
