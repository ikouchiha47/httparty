package tests

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"httparty"
	"httparty/engine"
	"testing"
)

func TestHttpartySingleLevel(t *testing.T) {
	executor := &MockHttpExecutor{}
	hp := httparty.NewHttparty(
		httparty.WithExecutor(executor),
	)

	parsed, err := engine.ParseFile("./files/rules.json")
	require.NoError(t, err, "should not have failed to parse file")

	step0, step1 := parsed[0].Steps[0], parsed[0].Steps[1]

	accessTokenRespBody := map[string]interface{}{"token": "AccessToken"}
	finalRespBody := map[string]interface{}{"data": "hello"}

	step0 = &engine.Step{
		Name:     step0.Name,
		Request:  step0.Request,
	}

	step1 = &engine.Step{
		Name: step1.Name,
		Request: step1.Request,
	}

	executor.On("RunIt", step1).Return(engine.Response{Body: finalRespBody}, nil)
	executor.On("RunIt", step0).Return(engine.Response{Body: accessTokenRespBody}, nil)

	xhp, err := hp.RunIt(parsed[0])
	require.NoError(t, err, "should not have failed to execute response")

	assert.Equal(t, xhp.GetRespValue(step1.Name), engine.Response{finalRespBody})
	executor.AssertExpectations(t)
}

func TestHttpartyMultipleSteps(t *testing.T) {
	executor := &MockHttpExecutor{}
	hp := httparty.NewHttparty(
		httparty.WithExecutor(executor),
	)

	parsed, err := engine.ParseFile("./files/multiple.json")
	require.NoError(t, err, "should not have failed to parse file")

	step0, step1 := parsed[0].Steps[0], parsed[1].Steps[0]

	require.Equal(t, step0.Name, "anthem_access_token")
	require.Equal(t, step1.Name, "anthem_plans")

	executor.On("RunIt", step0).Return(engine.Response{Body: engine.Body{"data": "foo"}}, nil)
	executor.On("RunIt", step1).Return(engine.Response{Body: engine.Body{"data": "bar"}}, nil)

	xhps, err := hp.RunAll(parsed)
	require.NoError(t, err, "should not have failed to make request")

	assert.Equal(t, xhps[0].GetRespValue("anthem_access_token"), engine.Response{engine.Body{"data": "foo"}})
	executor.AssertExpectations(t)
}

func TestHttpartyCyclicSteps(t *testing.T) {
	executor := &MockHttpExecutor{}
	hp := httparty.NewHttparty(
		httparty.WithExecutor(executor),
	)

	parsed, err := engine.ParseFile("./files/cyclic.json")
	require.NoError(t, err, "should not have failed to parse file")

	step0, step1 := parsed[0].Steps[0], parsed[0].Steps[1]

	require.Equal(t, step0.Name, "anthem_access_token")
	require.Equal(t, step1.Name, "anthem_plans")

	_, err = hp.RunAll(parsed)
	require.Error(t, err, "CyclicDependency")

	executor.AssertExpectations(t)
}
