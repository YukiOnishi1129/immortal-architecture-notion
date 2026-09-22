package port

import "context"

// NotionPage is the result of creating or updating a Notion page.
type NotionPage struct {
	PageID string
	URL    string
}

// NotionSection is one block of note content sent to Notion.
// It becomes a heading followed by a paragraph.
type NotionSection struct {
	Label   string
	Content string
}

// NotionClient abstracts the Notion API.
// The UseCase layer depends on this interface, not on HTTP.
type NotionClient interface {
	// CreatePage creates a page under parentPageID and returns its id and url.
	CreatePage(ctx context.Context, parentPageID, title string, sections []NotionSection) (*NotionPage, error)

	// UpdatePage replaces the content of an existing page.
	UpdatePage(ctx context.Context, pageID, title string, sections []NotionSection) (*NotionPage, error)

	// Trash moves a page to the Notion trash. The page keeps its id.
	Trash(ctx context.Context, pageID string) error

	// Restore brings a trashed page back.
	Restore(ctx context.Context, pageID string) (*NotionPage, error)
}
