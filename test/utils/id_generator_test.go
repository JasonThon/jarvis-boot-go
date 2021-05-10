package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"thingworks.net/thingworks/jarvis-boot/utils"
)

func TestUniqueId(t *testing.T) {
	id, err := utils.UniqueId(1)
	assert.Nil(t, err)
	assert.True(t, id > 0)
}

