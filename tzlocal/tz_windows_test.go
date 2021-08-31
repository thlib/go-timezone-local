package tzlocal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalTZ(t *testing.T) {
	s, err := localTZfromTzutil()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	s, err = localTZfromReg()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	s, err = LocalTZ()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	tmp := WinTZtoIANA
	WinTZtoIANA = map[string]string{}
	s, err = LocalTZ()
	assert.Error(t, err)
	assert.Empty(t, s)
	WinTZtoIANA = tmp
}
