package sessions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Signing(t *testing.T) {
	token := Create(SessionData{Name: "yundream", Email: "yundream@gmail.com", Admin: false})
	t.Log(token)
	session, ok := Validation(token)
	assert.True(t, ok, "")
	t.Log(session)
}
