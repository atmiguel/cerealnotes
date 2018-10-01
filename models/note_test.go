package models_test

import (
	"testing"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/test_util"
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
			test_util.Ok(t, err)
			test_util.Equals(t, val, cat)
		})
	}

}
