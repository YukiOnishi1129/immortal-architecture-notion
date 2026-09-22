//go:build e2e

package e2e

import (
	"context"
	"errors"
	"sync"

	"immortal-architecture-notion/backend/internal/port"
)

// errPageNotFound mirrors what the real API returns for an unknown page id.
var errPageNotFound = errors.New("notion: page not found")

// fakeNotion is an in-memory stand-in for the Notion API.
//
// It records what was asked of it so tests can assert on the calls, and it
// models the one rule that matters for this feature: trashing keeps the page,
// so restoring returns the same id and URL.
type fakeNotion struct {
	mu sync.Mutex

	pages   map[string]*fakePage
	nextID  int
	calls   []string
	failNow error
}

type fakePage struct {
	id       string
	parentID string
	title    string
	sections []port.NotionSection
	trashed  bool
}

func newFakeNotion() *fakeNotion {
	return &fakeNotion{pages: map[string]*fakePage{}}
}

var _ port.NotionClient = (*fakeNotion)(nil)

// failWith makes every subsequent call fail, simulating a Notion outage.
func (f *fakeNotion) failWith(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failNow = err
}

func (f *fakeNotion) recover() {
	f.failWith(nil)
}

// resetCalls forgets the calls made so far, so a test can assert on only the
// calls of the step it is exercising.
func (f *fakeNotion) resetCalls() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = nil
}

func (f *fakeNotion) callNames() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.calls))
	copy(out, f.calls)
	return out
}

func (f *fakeNotion) page(id string) *fakePage {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.pages[id]
	if !ok {
		return nil
	}
	cp := *p
	return &cp
}

func (f *fakeNotion) pageCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.pages)
}

func (f *fakeNotion) CreatePage(_ context.Context, parentPageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "CreatePage")
	if f.failNow != nil {
		return nil, f.failNow
	}

	f.nextID++
	id := fakePageID(f.nextID)
	f.pages[id] = &fakePage{id: id, parentID: parentPageID, title: title, sections: sections}
	return &port.NotionPage{PageID: id, URL: "https://www.notion.so/" + id}, nil
}

func (f *fakeNotion) UpdatePage(_ context.Context, pageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "UpdatePage")
	if f.failNow != nil {
		return nil, f.failNow
	}

	p, ok := f.pages[pageID]
	if !ok {
		return nil, errPageNotFound
	}
	p.title = title
	p.sections = sections
	return &port.NotionPage{PageID: p.id, URL: "https://www.notion.so/" + p.id}, nil
}

func (f *fakeNotion) Trash(_ context.Context, pageID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "Trash")
	if f.failNow != nil {
		return f.failNow
	}

	p, ok := f.pages[pageID]
	if !ok {
		return errPageNotFound
	}
	// Notion has no hard delete: the page stays, only flagged.
	p.trashed = true
	return nil
}

func (f *fakeNotion) Restore(_ context.Context, pageID string) (*port.NotionPage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "Restore")
	if f.failNow != nil {
		return nil, f.failNow
	}

	p, ok := f.pages[pageID]
	if !ok {
		return nil, errPageNotFound
	}
	p.trashed = false
	return &port.NotionPage{PageID: p.id, URL: "https://www.notion.so/" + p.id}, nil
}

func fakePageID(n int) string {
	const base = "beef0000000000000000000000000000"
	suffix := string(rune('0' + n%10))
	return base[:len(base)-1] + suffix
}
