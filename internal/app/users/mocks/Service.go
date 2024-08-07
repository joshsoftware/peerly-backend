// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import dto "github.com/joshsoftware/peerly-backend/internal/pkg/dto"
import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetIntranetUserData provides a mock function with given fields: ctx, req
func (_m *Service) GetIntranetUserData(ctx context.Context, req dto.GetIntranetUserDataReq) (dto.IntranetUserData, error) {
	ret := _m.Called(ctx, req)

	var r0 dto.IntranetUserData
	if rf, ok := ret.Get(0).(func(context.Context, dto.GetIntranetUserDataReq) dto.IntranetUserData); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(dto.IntranetUserData)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.GetIntranetUserDataReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListIntranetUsers provides a mock function with given fields: ctx, reqData
func (_m *Service) ListIntranetUsers(ctx context.Context, reqData dto.GetUserListReq) ([]dto.IntranetUserData, error) {
	ret := _m.Called(ctx, reqData)

	var r0 []dto.IntranetUserData
	if rf, ok := ret.Get(0).(func(context.Context, dto.GetUserListReq) []dto.IntranetUserData); ok {
		r0 = rf(ctx, reqData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.IntranetUserData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.GetUserListReq) error); ok {
		r1 = rf(ctx, reqData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginUser provides a mock function with given fields: ctx, u
func (_m *Service) LoginUser(ctx context.Context, u dto.IntranetUserData) (dto.LoginUserResp, error) {
	ret := _m.Called(ctx, u)

	var r0 dto.LoginUserResp
	if rf, ok := ret.Get(0).(func(context.Context, dto.IntranetUserData) dto.LoginUserResp); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Get(0).(dto.LoginUserResp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.IntranetUserData) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: ctx, u
func (_m *Service) RegisterUser(ctx context.Context, u dto.IntranetUserData) (dto.User, error) {
	ret := _m.Called(ctx, u)

	var r0 dto.User
	if rf, ok := ret.Get(0).(func(context.Context, dto.IntranetUserData) dto.User); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Get(0).(dto.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.IntranetUserData) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatePeerly provides a mock function with given fields: ctx, authToken
func (_m *Service) ValidatePeerly(ctx context.Context, authToken string) (dto.ValidateResp, error) {
	ret := _m.Called(ctx, authToken)

	var r0 dto.ValidateResp
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.ValidateResp); ok {
		r0 = rf(ctx, authToken)
	} else {
		r0 = ret.Get(0).(dto.ValidateResp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, authToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
