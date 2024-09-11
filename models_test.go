package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {

	acc, err := NewAccount("test", "name", "ibrox")

	assert.Nil(t, err)

	log.Printf("%+v\n", acc)

	assert.Equal(t, 1, 1)
}
