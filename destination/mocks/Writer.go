// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

// Writer is an autogenerated mock type for the Writer type
type Writer struct {
	mock.Mock
}

// Write provides a mock function with given fields: ctx, records
func (_m *Writer) Write(ctx context.Context, records []sdk.Record) error {
	ret := _m.Called(ctx, records)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []sdk.Record) error); ok {
		r0 = rf(ctx, records)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewWriterT interface {
	mock.TestingT
	Cleanup(func())
}

// NewWriter creates a new instance of Writer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWriter(t NewWriterT) *Writer {
	mock := &Writer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
