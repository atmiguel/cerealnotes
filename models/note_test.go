package models_test

import (
	"testing"

	"github.com/atmiguel/cerealnotes/models"
)

var deserializationTests = []models.Category{
	models.MARGINALIA,
	models.META,
	models.QUESTIONS,
	models.PREDICTIONS,
}

func TestDeserialization(t *testing.T) {
	for _, val := range deserializationTests {
		t.Run(val.String(), func(t *testing.T) {
			cat, err := models.DeserializeCategory(val.String())
			ok(t, err)
			equals(t, val, cat)
		})
	}

}
