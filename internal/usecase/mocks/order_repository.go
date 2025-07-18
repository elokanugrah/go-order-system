// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/elokanugrah/go-order-system/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: ctx, order
func (_m *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrderRepository creates a new instance of OrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderRepository {
	mock := &OrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
