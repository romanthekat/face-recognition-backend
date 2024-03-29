package services

import (
	"github.com/EvilKhaosKat/FaceRecognitionBackend/pkg/models"
	"math"
)

type EncodingComparator struct {
	persons models.PersonModel
}

func NewEncodingComparator(persons models.PersonModel) *EncodingComparator {
	return &EncodingComparator{persons: persons}
}

func (e *EncodingComparator) FindSamePerson(encoding Encoding) (*models.Person, error) {
	persons, err := e.persons.GetAll()
	if err != nil {
		return nil, err
	}

	var closestPerson *models.Person
	var closestDist = math.MaxFloat64

	for _, person := range persons {
		for _, encodingStr := range person.Encodings {
			personEncoding, err := NewEncoding(encodingStr)
			if err != nil {
				return nil, err
			}

			samePerson, dist, err := encoding.IsSame(personEncoding)
			if err != nil {
				return nil, err
			}
			if samePerson && dist < closestDist {
				closestPerson = person
			}
		}
	}

	return closestPerson, nil
}
