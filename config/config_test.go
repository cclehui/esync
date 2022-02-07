package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DecodeFromFile(t *testing.T) {
	configData := &Config{}

	_, err := configData.DecodeFromFile("./config.sample.yaml")
	assert.Equal(t, err, nil)
}
