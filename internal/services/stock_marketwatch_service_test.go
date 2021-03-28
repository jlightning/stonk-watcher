package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataFromMorningstar(t *testing.T) {
	s := assert.New(t)

	data, err := GetDataFromMarketWatch("SNE")
	s.NoError(err)
	s.NotNil(data)
}
