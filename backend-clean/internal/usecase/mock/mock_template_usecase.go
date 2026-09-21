// Code generated manually for gomock-based tests.
package mockusecase

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"immortal-architecture-cqrs/backend/internal/domain/template"
)

// MockTemplateRepository is a mock of port.TemplateRepository.
type MockTemplateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTemplateRepositoryMockRecorder
}

// MockTemplateRepositoryMockRecorder records invocations.
type MockTemplateRepositoryMockRecorder struct {
	mock *MockTemplateRepository
}

// NewMockTemplateRepository creates a new mock.
func NewMockTemplateRepository(ctrl *gomock.Controller) *MockTemplateRepository {
	mock := &MockTemplateRepository{ctrl: ctrl}
	mock.recorder = &MockTemplateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockTemplateRepository) EXPECT() *MockTemplateRepositoryMockRecorder {
	return m.recorder
}

func (m *MockTemplateRepository) List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filters)
	res0, _ := ret[0].([]template.WithUsage)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockTemplateRepositoryMockRecorder) List(ctx, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTemplateRepository)(nil).List), ctx, filters)
}

func (m *MockTemplateRepository) Get(ctx context.Context, id string) (*template.WithUsage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	res0, _ := ret[0].(*template.WithUsage)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockTemplateRepositoryMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTemplateRepository)(nil).Get), ctx, id)
}

func (m *MockTemplateRepository) Create(ctx context.Context, tpl template.Template) (*template.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tpl)
	res0, _ := ret[0].(*template.Template)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockTemplateRepositoryMockRecorder) Create(ctx, tpl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTemplateRepository)(nil).Create), ctx, tpl)
}

func (m *MockTemplateRepository) Update(ctx context.Context, tpl template.Template) (*template.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tpl)
	res0, _ := ret[0].(*template.Template)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockTemplateRepositoryMockRecorder) Update(ctx, tpl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTemplateRepository)(nil).Update), ctx, tpl)
}

func (m *MockTemplateRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTemplateRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTemplateRepository)(nil).Delete), ctx, id)
}

func (m *MockTemplateRepository) ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceFields", ctx, templateID, fields)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTemplateRepositoryMockRecorder) ReplaceFields(ctx, templateID, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceFields", reflect.TypeOf((*MockTemplateRepository)(nil).ReplaceFields), ctx, templateID, fields)
}

// MockTxManager is a mock of port.TxManager.
type MockTxManager struct {
	ctrl     *gomock.Controller
	recorder *MockTxManagerMockRecorder
}

// MockTxManagerMockRecorder records calls.
type MockTxManagerMockRecorder struct {
	mock *MockTxManager
}

// NewMockTxManager creates a new mock.
func NewMockTxManager(ctrl *gomock.Controller) *MockTxManager {
	mock := &MockTxManager{ctrl: ctrl}
	mock.recorder = &MockTxManagerMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockTxManager) EXPECT() *MockTxManagerMockRecorder {
	return m.recorder
}

func (m *MockTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithinTransaction", ctx, fn)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTxManagerMockRecorder) WithinTransaction(ctx, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithinTransaction", reflect.TypeOf((*MockTxManager)(nil).WithinTransaction), ctx, fn)
}

// MockTemplateOutputPort is a mock of port.TemplateOutputPort.
type MockTemplateOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockTemplateOutputPortMockRecorder
}

// MockTemplateOutputPortMockRecorder records calls.
type MockTemplateOutputPortMockRecorder struct {
	mock *MockTemplateOutputPort
}

// NewMockTemplateOutputPort creates a new mock.
func NewMockTemplateOutputPort(ctrl *gomock.Controller) *MockTemplateOutputPort {
	mock := &MockTemplateOutputPort{ctrl: ctrl}
	mock.recorder = &MockTemplateOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockTemplateOutputPort) EXPECT() *MockTemplateOutputPortMockRecorder {
	return m.recorder
}

func (m *MockTemplateOutputPort) PresentTemplateList(ctx context.Context, templates []template.WithUsage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentTemplateList", ctx, templates)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTemplateOutputPortMockRecorder) PresentTemplateList(ctx, templates any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentTemplateList", reflect.TypeOf((*MockTemplateOutputPort)(nil).PresentTemplateList), ctx, templates)
}

func (m *MockTemplateOutputPort) PresentTemplate(ctx context.Context, tpl *template.WithUsage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentTemplate", ctx, tpl)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTemplateOutputPortMockRecorder) PresentTemplate(ctx, tpl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentTemplate", reflect.TypeOf((*MockTemplateOutputPort)(nil).PresentTemplate), ctx, tpl)
}

func (m *MockTemplateOutputPort) PresentTemplateDeleted(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentTemplateDeleted", ctx)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockTemplateOutputPortMockRecorder) PresentTemplateDeleted(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentTemplateDeleted", reflect.TypeOf((*MockTemplateOutputPort)(nil).PresentTemplateDeleted), ctx)
}
