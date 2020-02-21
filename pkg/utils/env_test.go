package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	// Assert correct value
	err := os.Setenv("TEST_STRING", "my_value")
	assert.Nil(t, err)
	assert.Equal(t, "my_value", GetEnv("TEST_STRING", "my_value"))

	// Assert default value
	assert.Equal(t, "unset_value", GetEnv("UNSET_TEST_STRING", "unset_value"))

	// Assert incorrect value
	err = os.Setenv("TEST_STRING", "updated_value")
	assert.Nil(t, err)
	assert.False(t, "old_value" == GetEnv("TEST_STRING", "old_value"))
}

func TestGetEnvBool(t *testing.T) {

	trues := []string{"1", "t", "T", "true", "TRUE", "True"}
	falses := []string{"0", "f", "F", "false", "FALSE", "False"}
	invalid := []string{"set", "foo", "bar", "2", "@"}

	// Assert true values are parsed correctly
	for _, s := range trues {
		err := os.Setenv("TEST_BOOL", s)
		assert.Nil(t, err)
		assert.True(t, GetEnvBool("TEST_BOOL"))
	}

	// Assert false values are parsed correctly
	for _, s := range falses {
		err := os.Setenv("TEST_BOOL", s)
		assert.Nil(t, err)
		assert.False(t, GetEnvBool("TEST_BOOL"))
	}

	// Assert non-boolean values return false
	for _, s := range invalid {
		err := os.Setenv("TEST_BOOL", s)
		assert.Nil(t, err)
		assert.False(t, GetEnvBool("TEST_BOOL"))
	}
}
