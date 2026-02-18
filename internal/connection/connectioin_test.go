package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConnection(t *testing.T) {
	conn, err := InitConnection()
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}
