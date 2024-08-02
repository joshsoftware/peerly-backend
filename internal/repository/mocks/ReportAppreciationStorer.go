// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import dto "github.com/joshsoftware/peerly-backend/internal/pkg/dto"
import mock "github.com/stretchr/testify/mock"
import repository "github.com/joshsoftware/peerly-backend/internal/repository"

// ReportAppreciationStorer is an autogenerated mock type for the ReportAppreciationStorer type
type ReportAppreciationStorer struct {
	mock.Mock
}

// CheckAppreciation provides a mock function with given fields: ctx, reqData
func (_m *ReportAppreciationStorer) CheckAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (bool, error) {
	ret := _m.Called(ctx, reqData)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, dto.ReportAppreciationReq) bool); ok {
		r0 = rf(ctx, reqData)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ReportAppreciationReq) error); ok {
		r1 = rf(ctx, reqData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckDuplicateReport provides a mock function with given fields: ctx, reqData
func (_m *ReportAppreciationStorer) CheckDuplicateReport(ctx context.Context, reqData dto.ReportAppreciationReq) (bool, error) {
	ret := _m.Called(ctx, reqData)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, dto.ReportAppreciationReq) bool); ok {
		r0 = rf(ctx, reqData)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ReportAppreciationReq) error); ok {
		r1 = rf(ctx, reqData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckResolution provides a mock function with given fields: ctx, id
func (_m *ReportAppreciationStorer) CheckResolution(ctx context.Context, id int64) (bool, int64, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, int64) int64); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteAppreciation provides a mock function with given fields: ctx, moderationReq
func (_m *ReportAppreciationStorer) DeleteAppreciation(ctx context.Context, moderationReq dto.ModerationReq) error {
	ret := _m.Called(ctx, moderationReq)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.ModerationReq) error); ok {
		r0 = rf(ctx, moderationReq)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSenderAndReceiver provides a mock function with given fields: ctx, reqData
func (_m *ReportAppreciationStorer) GetSenderAndReceiver(ctx context.Context, reqData dto.ReportAppreciationReq) (dto.GetSenderAndReceiverResp, error) {
	ret := _m.Called(ctx, reqData)

	var r0 dto.GetSenderAndReceiverResp
	if rf, ok := ret.Get(0).(func(context.Context, dto.ReportAppreciationReq) dto.GetSenderAndReceiverResp); ok {
		r0 = rf(ctx, reqData)
	} else {
		r0 = ret.Get(0).(dto.GetSenderAndReceiverResp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ReportAppreciationReq) error); ok {
		r1 = rf(ctx, reqData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListReportedAppreciations provides a mock function with given fields: ctx
func (_m *ReportAppreciationStorer) ListReportedAppreciations(ctx context.Context) ([]repository.ListReportedAppreciations, error) {
	ret := _m.Called(ctx)

	var r0 []repository.ListReportedAppreciations
	if rf, ok := ret.Get(0).(func(context.Context) []repository.ListReportedAppreciations); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.ListReportedAppreciations)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReportAppreciation provides a mock function with given fields: ctx, reportReq
func (_m *ReportAppreciationStorer) ReportAppreciation(ctx context.Context, reportReq dto.ReportAppreciationReq) (dto.ReportAppricaitionResp, error) {
	ret := _m.Called(ctx, reportReq)

	var r0 dto.ReportAppricaitionResp
	if rf, ok := ret.Get(0).(func(context.Context, dto.ReportAppreciationReq) dto.ReportAppricaitionResp); ok {
		r0 = rf(ctx, reportReq)
	} else {
		r0 = ret.Get(0).(dto.ReportAppricaitionResp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.ReportAppreciationReq) error); ok {
		r1 = rf(ctx, reportReq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewReportAppreciationStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReportAppreciationStorer {
	mock := &ReportAppreciationStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}