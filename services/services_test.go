package services

import (
	"testing"

	fakerLib "github.com/jaswdr/faker"
)

var faker = fakerLib.New()

func TestSendEmail(t *testing.T) {
	err := SendEmail(faker.Internet().Email(), faker.Lorem().Word(), faker.Lorem().Sentence(10))
	if err != nil {
		t.Error("Error should be nil but got", err)
	}
}
