package services

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Encoding []float64

//TODO custom errors?

const similarityMaxThreshold = 0.64
const encodingMinlength = 42

func NewEncoding(encoding string) (Encoding, error) {
	encodingStringLen := utf8.RuneCountInString(encoding)

	if encodingStringLen < encodingMinlength { //usually encoding is 128 float numbers, so it's pretty long in text form
		return nil, errors.New(fmt.Sprint("too short encoding, probably incorrect data:", encoding))
	}

	numbers := strings.Split(encoding, " ")

	var result Encoding

	for _, numberStr := range numbers {
		if utf8.RuneCountInString(numberStr) == 0 {
			continue
		}

		numberStr = strings.Trim(numberStr, "\n")
		number, err := strconv.ParseFloat(numberStr, 16) //technically we get 16 bit length from ML model
		if err != nil {
			return nil, err
		}

		result = append(result, number)
	}

	return Encoding(result), nil
}

//GetDist returns L2 (Euclidean) distance between to encodings/vectors
//TODO should it be method for Encoding, or a helper method in package?
func (e Encoding) GetDist(otherEncoding Encoding) (float64, error) {
	result := 0.0

	if len(e) != len(otherEncoding) {
		return -1, fmt.Errorf("different encodings length: %d and %d", len(e), len(otherEncoding))
	}

	for i := 0; i < len(e); i++ {
		v1 := e[i]
		v2 := otherEncoding[i]

		result += math.Pow(v1-v2, 2)
	}

	result = math.Sqrt(result)

	return result, nil
}

//IsSame checks if encodings are for the same person
//returns flag, distance, and error
func (e Encoding) IsSame(otherEncoding Encoding) (bool, float64, error) {
	dist, err := e.GetDist(otherEncoding)
	if err != nil {
		return false, -1, err
	}

	return dist < similarityMaxThreshold, dist, nil
}

func (e Encoding) String() string {
	var buffer bytes.Buffer

	for _, num := range e {
		buffer.WriteString(fmt.Sprintf("%f", num))
	}

	return buffer.String()
}
