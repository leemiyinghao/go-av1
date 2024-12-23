// Code generated by mockery v2.46.3. DO NOT EDIT.

package task_flow

import (
	task_flow "github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	mock "github.com/stretchr/testify/mock"
)

// MockExecutableFactory is an autogenerated mock type for the ExecutableFactory type
type MockExecutableFactory struct {
	mock.Mock
}

type MockExecutableFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExecutableFactory) EXPECT() *MockExecutableFactory_Expecter {
	return &MockExecutableFactory_Expecter{mock: &_m.Mock}
}

// GenerateExecutable provides a mock function with given fields: file
func (_m *MockExecutableFactory) GenerateExecutable(file string) task_flow.Executable {
	ret := _m.Called(file)

	if len(ret) == 0 {
		panic("no return value specified for GenerateExecutable")
	}

	var r0 task_flow.Executable
	if rf, ok := ret.Get(0).(func(string) task_flow.Executable); ok {
		r0 = rf(file)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(task_flow.Executable)
		}
	}

	return r0
}

// MockExecutableFactory_GenerateExecutable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateExecutable'
type MockExecutableFactory_GenerateExecutable_Call struct {
	*mock.Call
}

// GenerateExecutable is a helper method to define mock.On call
//   - file string
func (_e *MockExecutableFactory_Expecter) GenerateExecutable(file interface{}) *MockExecutableFactory_GenerateExecutable_Call {
	return &MockExecutableFactory_GenerateExecutable_Call{Call: _e.mock.On("GenerateExecutable", file)}
}

func (_c *MockExecutableFactory_GenerateExecutable_Call) Run(run func(file string)) *MockExecutableFactory_GenerateExecutable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockExecutableFactory_GenerateExecutable_Call) Return(_a0 task_flow.Executable) *MockExecutableFactory_GenerateExecutable_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExecutableFactory_GenerateExecutable_Call) RunAndReturn(run func(string) task_flow.Executable) *MockExecutableFactory_GenerateExecutable_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExecutableFactory creates a new instance of MockExecutableFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecutableFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecutableFactory {
	mock := &MockExecutableFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
