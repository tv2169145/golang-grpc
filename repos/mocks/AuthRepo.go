// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tv2169145/golang-grpc/repos (interfaces: AuthRepo)

// Package mock_repos is a generated GoMock package.
package mock_repos

import (
	gomock "github.com/golang/mock/gomock"
	jwt "github.com/pascaldekloe/jwt"
	types "github.com/tv2169145/golang-grpc/types"
	reflect "reflect"
)

// MockAuthRepo is a mock of AuthRepo interface
type MockAuthRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepoMockRecorder
}

// MockAuthRepoMockRecorder is the mock recorder for MockAuthRepo
type MockAuthRepoMockRecorder struct {
	mock *MockAuthRepo
}

// NewMockAuthRepo creates a new mock instance
func NewMockAuthRepo(ctrl *gomock.Controller) *MockAuthRepo {
	mock := &MockAuthRepo{ctrl: ctrl}
	mock.recorder = &MockAuthRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthRepo) EXPECT() *MockAuthRepoMockRecorder {
	return m.recorder
}

// GetDataFromToken mocks base method
func (m *MockAuthRepo) GetDataFromToken(arg0 string) (*types.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataFromToken", arg0)
	ret0, _ := ret[0].(*types.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataFromToken indicates an expected call of GetDataFromToken
func (mr *MockAuthRepoMockRecorder) GetDataFromToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataFromToken", reflect.TypeOf((*MockAuthRepo)(nil).GetDataFromToken), arg0)
}

// GetNewClaims mocks base method
func (m *MockAuthRepo) GetNewClaims(arg0 string, arg1 map[string]interface{}) *jwt.Claims {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewClaims", arg0, arg1)
	ret0, _ := ret[0].(*jwt.Claims)
	return ret0
}

// GetNewClaims indicates an expected call of GetNewClaims
func (mr *MockAuthRepoMockRecorder) GetNewClaims(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewClaims", reflect.TypeOf((*MockAuthRepo)(nil).GetNewClaims), arg0, arg1)
}

// GetSignedToken mocks base method
func (m *MockAuthRepo) GetSignedToken(arg0 *jwt.Claims) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSignedToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSignedToken indicates an expected call of GetSignedToken
func (mr *MockAuthRepoMockRecorder) GetSignedToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSignedToken", reflect.TypeOf((*MockAuthRepo)(nil).GetSignedToken), arg0)
}
