// Code generated by mockery v1.0.0
package automock

import context "context"

import io "io"
import mock "github.com/stretchr/testify/mock"

// Mutator is an autogenerated mock type for the Mutator type
type Mutator struct {
	mock.Mock
}

// Mutate provides a mock function with given fields: ctx, contentType, reader, metadata
func (_m *Mutator) Mutate(ctx context.Context, contentType string, reader io.Reader, metadata string) ([]byte, error) {
	ret := _m.Called(ctx, contentType, reader, metadata)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string, io.Reader, string) []byte); ok {
		r0 = rf(ctx, contentType, reader, metadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, io.Reader, string) error); ok {
		r1 = rf(ctx, contentType, reader, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
