package notion

import (
	"strings"
	"testing"

	"immortal-architecture-notion/backend/internal/port"
)

func TestToBlocks(t *testing.T) {
	tests := []struct {
		name      string
		sections  []port.NotionSection
		wantCount int
		wantTypes []string
	}{
		{
			name:      "[Empty] no sections",
			sections:  nil,
			wantCount: 0,
		},
		{
			name: "[Single] one section becomes heading and paragraph",
			sections: []port.NotionSection{
				{Label: "背景", Content: "なぜ作るか"},
			},
			wantCount: 2,
			wantTypes: []string{"heading_2", "paragraph"},
		},
		{
			name: "[Multiple] two sections keep their order",
			sections: []port.NotionSection{
				{Label: "背景", Content: "なぜ"},
				{Label: "対策", Content: "どうする"},
			},
			wantCount: 4,
			wantTypes: []string{"heading_2", "paragraph", "heading_2", "paragraph"},
		},
		{
			name: "[Empty content] still produces a paragraph",
			sections: []port.NotionSection{
				{Label: "任意項目", Content: ""},
			},
			wantCount: 2,
			wantTypes: []string{"heading_2", "paragraph"},
		},
		{
			name: "[No label] only a paragraph",
			sections: []port.NotionSection{
				{Label: "", Content: "本文だけ"},
			},
			wantCount: 1,
			wantTypes: []string{"paragraph"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := toBlocks(tt.sections)

			if len(blocks) != tt.wantCount {
				t.Fatalf("block count = %d, want %d", len(blocks), tt.wantCount)
			}
			for i, want := range tt.wantTypes {
				if blocks[i].Type != want {
					t.Errorf("block[%d].Type = %q, want %q", i, blocks[i].Type, want)
				}
			}
		})
	}
}

func TestToBlocks_LongContentIsSplit(t *testing.T) {
	long := strings.Repeat("あ", maxTextLength+500)

	blocks := toBlocks([]port.NotionSection{{Label: "長文", Content: long}})

	// heading + two paragraphs
	if len(blocks) != 3 {
		t.Fatalf("block count = %d, want 3", len(blocks))
	}
	for i := 1; i < len(blocks); i++ {
		content := blocks[i].Paragraph.RichText[0].Text.Content
		if len([]rune(content)) > maxTextLength {
			t.Errorf("block[%d] length = %d, exceeds %d", i, len([]rune(content)), maxTextLength)
		}
	}
}

func TestToBlocks_RespectsBlockLimit(t *testing.T) {
	sections := make([]port.NotionSection, 80) // 80 sections -> 160 blocks
	for i := range sections {
		sections[i] = port.NotionSection{Label: "項目", Content: "内容"}
	}

	blocks := toBlocks(sections)

	if len(blocks) != maxBlocksPerReq {
		t.Fatalf("block count = %d, want %d", len(blocks), maxBlocksPerReq)
	}
}

func TestSplitText(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantParts int
	}{
		{name: "[Empty] yields one empty chunk", content: "", wantParts: 1},
		{name: "[Short] stays as one chunk", content: "短い文章", wantParts: 1},
		{
			name:      "[Exactly at limit] stays as one chunk",
			content:   strings.Repeat("a", maxTextLength),
			wantParts: 1,
		},
		{
			name:      "[Over limit] is split",
			content:   strings.Repeat("a", maxTextLength+1),
			wantParts: 2,
		},
		{
			name:      "[Multibyte] counts runes not bytes",
			content:   strings.Repeat("あ", maxTextLength),
			wantParts: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := splitText(tt.content)

			if len(parts) != tt.wantParts {
				t.Fatalf("parts = %d, want %d", len(parts), tt.wantParts)
			}
			for i, p := range parts {
				if len([]rune(p)) > maxTextLength {
					t.Errorf("part[%d] length = %d, exceeds %d", i, len([]rune(p)), maxTextLength)
				}
			}
		})
	}
}

func TestTruncateTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  int
	}{
		{name: "[Short] unchanged", title: "ノート", want: 3},
		{
			name:  "[Over limit] truncated",
			title: strings.Repeat("あ", maxTextLength+100),
			want:  maxTextLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateTitle(tt.title)

			if len([]rune(got)) != tt.want {
				t.Fatalf("length = %d, want %d", len([]rune(got)), tt.want)
			}
		})
	}
}
