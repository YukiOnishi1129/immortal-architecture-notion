package notion

import "immortal-architecture-notion/backend/internal/port"

// Notion limits how much text one rich_text object can hold,
// and how many blocks a single request may contain.
const (
	maxTextLength   = 2000
	maxBlocksPerReq = 100
)

// richText is a Notion rich text object.
type richText struct {
	Type string   `json:"type"`
	Text textBody `json:"text"`
}

type textBody struct {
	Content string `json:"content"`
}

// block is a Notion block object. Only the field matching Type is set.
type block struct {
	Object    string     `json:"object"`
	Type      string     `json:"type"`
	Heading2  *blockText `json:"heading_2,omitempty"`
	Paragraph *blockText `json:"paragraph,omitempty"`
}

type blockText struct {
	RichText []richText `json:"rich_text"`
}

// toBlocks converts note sections into Notion blocks.
// Each section becomes a heading followed by one or more paragraphs.
//
// Empty content still produces an empty paragraph so that the
// section headings stay aligned with the note structure.
func toBlocks(sections []port.NotionSection) []block {
	blocks := make([]block, 0, len(sections)*2)

	for _, s := range sections {
		if s.Label != "" {
			blocks = append(blocks, headingBlock(s.Label))
		}
		for _, chunk := range splitText(s.Content) {
			blocks = append(blocks, paragraphBlock(chunk))
		}
	}

	if len(blocks) > maxBlocksPerReq {
		blocks = blocks[:maxBlocksPerReq]
	}
	return blocks
}

func headingBlock(text string) block {
	return block{
		Object:   "block",
		Type:     "heading_2",
		Heading2: &blockText{RichText: []richText{{Type: "text", Text: textBody{Content: text}}}},
	}
}

func paragraphBlock(text string) block {
	return block{
		Object:    "block",
		Type:      "paragraph",
		Paragraph: &blockText{RichText: []richText{{Type: "text", Text: textBody{Content: text}}}},
	}
}

// splitText breaks content into chunks Notion accepts.
// An empty string yields a single empty chunk so the block is still created.
func splitText(content string) []string {
	if content == "" {
		return []string{""}
	}

	runes := []rune(content)
	if len(runes) <= maxTextLength {
		return []string{content}
	}

	chunks := make([]string, 0, (len(runes)/maxTextLength)+1)
	for start := 0; start < len(runes); start += maxTextLength {
		end := start + maxTextLength
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

// truncateTitle keeps the title within Notion's limit.
func truncateTitle(title string) string {
	runes := []rune(title)
	if len(runes) <= maxTextLength {
		return title
	}
	return string(runes[:maxTextLength])
}
