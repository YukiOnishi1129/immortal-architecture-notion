package mockusecase

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"immortal-architecture-cqrs/backend/internal/domain/account"
)

// MockAccountRepository is a mock of port.AccountRepository.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder records invocations.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

func (m *MockAccountRepository) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertOAuthAccount", ctx, input)
	res0, _ := ret[0].(*account.Account)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockAccountRepositoryMockRecorder) UpsertOAuthAccount(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertOAuthAccount", reflect.TypeOf((*MockAccountRepository)(nil).UpsertOAuthAccount), ctx, input)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	res0, _ := ret[0].(*account.Account)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockAccountRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAccountRepository)(nil).GetByID), ctx, id)
}

func (m *MockAccountRepository) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	res0, _ := ret[0].(*account.Account)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockAccountRepositoryMockRecorder) GetByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockAccountRepository)(nil).GetByEmail), ctx, email)
}

// MockAccountOutputPort is a mock of port.AccountOutputPort.
type MockAccountOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockAccountOutputPortMockRecorder
}

// MockAccountOutputPortMockRecorder records invocations.
type MockAccountOutputPortMockRecorder struct {
	mock *MockAccountOutputPort
}

// NewMockAccountOutputPort creates a new mock.
func NewMockAccountOutputPort(ctrl *gomock.Controller) *MockAccountOutputPort {
	mock := &MockAccountOutputPort{ctrl: ctrl}
	mock.recorder = &MockAccountOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockAccountOutputPort) EXPECT() *MockAccountOutputPortMockRecorder {
	return m.recorder
}

func (m *MockAccountOutputPort) PresentAccount(ctx context.Context, acc *account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentAccount", ctx, acc)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockAccountOutputPortMockRecorder) PresentAccount(ctx, acc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentAccount", reflect.TypeOf((*MockAccountOutputPort)(nil).PresentAccount), ctx, acc)
}
