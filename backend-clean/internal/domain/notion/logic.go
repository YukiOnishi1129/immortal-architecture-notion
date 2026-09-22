package notion

import (
	"regexp"
	"strings"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
)

// pageIDPattern matches a Notion page id: 32 hex characters,
// with or without the UUID hyphens.
var pageIDPattern = regexp.MustCompile(`[0-9a-fA-F]{32}|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

// IsFirstSync reports whether the note has never been sent to Notion.
// A nil receiver means no sync record exists at all.
func (s *Sync) IsFirstSync() bool {
	return s == nil || s.PageID == nil || *s.PageID == ""
}

// PageIDOrEmpty returns the stored page id, or an empty string when unset.
func (s *Sync) PageIDOrEmpty() string {
	if s == nil || s.PageID == nil {
		return ""
	}
	return *s.PageID
}

// ExtractPageID pulls a Notion page id out of a page URL.
// Both hyphenated and non-hyphenated ids are accepted, and the
// returned value is always normalized to the non-hyphenated form.
//
//	https://notion.so/workspace/1429989fe8ac4effbc8f57f56486db54?v=x
//	-> 1429989fe8ac4effbc8f57f56486db54
func ExtractPageID(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", domainerr.ErrNotionParentNotSet
	}

	// Query strings and fragments may contain other ids (for example
	// ?v=<view id>), so they are stripped before scanning.
	if idx := strings.IndexAny(trimmed, "?#"); idx >= 0 {
		trimmed = trimmed[:idx]
	}

	matches := pageIDPattern.FindAllString(trimmed, -1)
	if len(matches) == 0 {
		return "", domainerr.ErrNotionParentNotSet
	}

	// The page id is the last path segment.
	last := matches[len(matches)-1]
	return strings.ToLower(strings.ReplaceAll(last, "-", "")), nil
}

// ValidateParentPageID checks that a stored parent page id is usable.
// The id must be a well-formed Notion page id, not merely non-empty.
func ValidateParentPageID(pageID string) error {
	_, err := ExtractPageID(pageID)
	return err
}
