// Code generated manually for gomock-based tests.
package mockusecase

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"immortal-architecture-notion/backend/internal/port"
)

// MockNotionClient is a mock of port.NotionClient.
type MockNotionClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotionClientMockRecorder
}

// MockNotionClientMockRecorder records invocations.
type MockNotionClientMockRecorder struct {
	mock *MockNotionClient
}

// NewMockNotionClient creates a new mock.
func NewMockNotionClient(ctrl *gomock.Controller) *MockNotionClient {
	mock := &MockNotionClient{ctrl: ctrl}
	mock.recorder = &MockNotionClientMockRecorder{mock: mock}
	return mock
}

// EXPECT returns the recorder.
func (m *MockNotionClient) EXPECT() *MockNotionClientMockRecorder {
	return m.recorder
}

func (m *MockNotionClient) CreatePage(ctx context.Context, parentPageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePage", ctx, parentPageID, title, sections)
	res0, _ := ret[0].(*port.NotionPage)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNotionClientMockRecorder) CreatePage(ctx, parentPageID, title, sections any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePage", reflect.TypeOf((*MockNotionClient)(nil).CreatePage), ctx, parentPageID, title, sections)
}

func (m *MockNotionClient) UpdatePage(ctx context.Context, pageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePage", ctx, pageID, title, sections)
	res0, _ := ret[0].(*port.NotionPage)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNotionClientMockRecorder) UpdatePage(ctx, pageID, title, sections any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePage", reflect.TypeOf((*MockNotionClient)(nil).UpdatePage), ctx, pageID, title, sections)
}

func (m *MockNotionClient) Trash(ctx context.Context, pageID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trash", ctx, pageID)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNotionClientMockRecorder) Trash(ctx, pageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trash", reflect.TypeOf((*MockNotionClient)(nil).Trash), ctx, pageID)
}

func (m *MockNotionClient) Restore(ctx context.Context, pageID string) (*port.NotionPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, pageID)
	res0, _ := ret[0].(*port.NotionPage)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNotionClientMockRecorder) Restore(ctx, pageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockNotionClient)(nil).Restore), ctx, pageID)
}

var _ port.NotionClient = (*MockNotionClient)(nil)
