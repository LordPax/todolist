package services

import (
	"fmt"
	"testing"
	"time"
	"todolist/services"

	_ "github.com/glebarez/go-sqlite"
	fakerLib "github.com/jaswdr/faker"
)


func TestSendEmail(t *testing.T) {
	err := services.SendEmail(fakerLib.Internet().Email(), fakerLib.Lorem().Word(), fakerLib.Lorem().Sentence(10))*
	if err != nil {
		t.Error("Error should be nil but got", err)
	}
}
