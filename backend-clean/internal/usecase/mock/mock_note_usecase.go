// Code generated manually for gomock-based tests.
package mockusecase

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"immortal-architecture-notion/backend/internal/domain/note"
)

// MockNoteRepository is a mock of port.NoteRepository.
type MockNoteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNoteRepositoryMockRecorder
}

// MockNoteRepositoryMockRecorder records invocations.
type MockNoteRepositoryMockRecorder struct {
	mock *MockNoteRepository
}

// NewMockNoteRepository creates a new mock.
func NewMockNoteRepository(ctrl *gomock.Controller) *MockNoteRepository {
	mock := &MockNoteRepository{ctrl: ctrl}
	mock.recorder = &MockNoteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockNoteRepository) EXPECT() *MockNoteRepositoryMockRecorder {
	return m.recorder
}

func (m *MockNoteRepository) List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filters)
	res0, _ := ret[0].([]note.WithMeta)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) List(ctx, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNoteRepository)(nil).List), ctx, filters)
}

func (m *MockNoteRepository) Get(ctx context.Context, id string) (*note.WithMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	res0, _ := ret[0].(*note.WithMeta)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNoteRepository)(nil).Get), ctx, id)
}

func (m *MockNoteRepository) Create(ctx context.Context, n note.Note) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, n)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Create(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNoteRepository)(nil).Create), ctx, n)
}

func (m *MockNoteRepository) Update(ctx context.Context, n note.Note) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, n)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Update(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNoteRepository)(nil).Update), ctx, n)
}

func (m *MockNoteRepository) UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) UpdateStatus(ctx, id, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockNoteRepository)(nil).UpdateStatus), ctx, id, status)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNoteRepository)(nil).Delete), ctx, id)
}

func (m *MockNoteRepository) ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceSections", ctx, noteID, sections)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteRepositoryMockRecorder) ReplaceSections(ctx, noteID, sections any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceSections", reflect.TypeOf((*MockNoteRepository)(nil).ReplaceSections), ctx, noteID, sections)
}

// MockNoteReadModelRepository is a mock of port.NoteReadModelRepository.
type MockNoteReadModelRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNoteReadModelRepositoryMockRecorder
}

// MockNoteReadModelRepositoryMockRecorder records invocations.
type MockNoteReadModelRepositoryMockRecorder struct {
	mock *MockNoteReadModelRepository
}

// NewMockNoteReadModelRepository creates a new mock.
func NewMockNoteReadModelRepository(ctrl *gomock.Controller) *MockNoteReadModelRepository {
	mock := &MockNoteReadModelRepository{ctrl: ctrl}
	mock.recorder = &MockNoteReadModelRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockNoteReadModelRepository) EXPECT() *MockNoteReadModelRepositoryMockRecorder {
	return m.recorder
}

func (m *MockNoteReadModelRepository) List(ctx context.Context, filters note.Filters) ([]note.ReadModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filters)
	res0, _ := ret[0].([]note.ReadModel)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteReadModelRepositoryMockRecorder) List(ctx, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNoteReadModelRepository)(nil).List), ctx, filters)
}

func (m *MockNoteReadModelRepository) Get(ctx context.Context, id string) (*note.ReadModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	res0, _ := ret[0].(*note.ReadModel)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteReadModelRepositoryMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNoteReadModelRepository)(nil).Get), ctx, id)
}

func (m *MockNoteReadModelRepository) Upsert(ctx context.Context, model note.ReadModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", ctx, model)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteReadModelRepositoryMockRecorder) Upsert(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockNoteReadModelRepository)(nil).Upsert), ctx, model)
}

func (m *MockNoteReadModelRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteReadModelRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNoteReadModelRepository)(nil).Delete), ctx, id)
}

// MockNoteCommandOutputPort is a mock of port.NoteCommandOutputPort.
type MockNoteCommandOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockNoteCommandOutputPortMockRecorder
}

// MockNoteCommandOutputPortMockRecorder records invocations.
type MockNoteCommandOutputPortMockRecorder struct {
	mock *MockNoteCommandOutputPort
}

// NewMockNoteCommandOutputPort creates a new mock.
func NewMockNoteCommandOutputPort(ctrl *gomock.Controller) *MockNoteCommandOutputPort {
	mock := &MockNoteCommandOutputPort{ctrl: ctrl}
	mock.recorder = &MockNoteCommandOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockNoteCommandOutputPort) EXPECT() *MockNoteCommandOutputPortMockRecorder {
	return m.recorder
}

func (m *MockNoteCommandOutputPort) PresentNote(ctx context.Context, n *note.WithMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentNote", ctx, n)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteCommandOutputPortMockRecorder) PresentNote(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentNote", reflect.TypeOf((*MockNoteCommandOutputPort)(nil).PresentNote), ctx, n)
}

func (m *MockNoteCommandOutputPort) PresentNoteDeleted(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentNoteDeleted", ctx)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteCommandOutputPortMockRecorder) PresentNoteDeleted(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentNoteDeleted", reflect.TypeOf((*MockNoteCommandOutputPort)(nil).PresentNoteDeleted), ctx)
}
