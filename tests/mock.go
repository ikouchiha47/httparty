package tests

import (
	"github.com/stretchr/testify/mock"
	"httparty/engine"
)

type MockHttpExecutor struct {
	mock.Mock
}

func (me *MockHttpExecutor) RunIt(step *engine.Step) (engine.Response, error) {
	args := me.Called(step)

	return args.Get(0).(engine.Response), args.Error(1)
}