package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VladisP/media-savior/internal/common/validator"
)

func TestNewConfig(t *testing.T) {
	err := os.Setenv(configPathEnv, "./config.example.yml")
	if err != nil {
		t.Error(err)
	}

	v := validator.NewValidator()
	_, err = NewConfig(v)

	assert.Nil(t, err)
}
