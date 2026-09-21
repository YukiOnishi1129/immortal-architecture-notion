package presenter

import (
	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/domain/note"
)

func toNoteResponse(n note.WithMeta) openapi.ModelsNoteResponse {
	sections := make([]openapi.ModelsSection, 0, len(n.Sections))
	for _, s := range n.Sections {
		sections = append(sections, openapi.ModelsSection{
			Id:         s.Section.ID,
			FieldId:    s.Section.FieldID,
			FieldLabel: s.FieldLabel,
			Content:    s.Section.Content,
			IsRequired: s.IsRequired,
		})
	}
	return openapi.ModelsNoteResponse{
		Id:           n.Note.ID,
		Title:        n.Note.Title,
		TemplateId:   n.Note.TemplateID,
		TemplateName: n.TemplateName,
		OwnerId:      n.Note.OwnerID,
		Owner: openapi.ModelsAccountSummary{
			Id:        n.Note.OwnerID,
			FirstName: n.OwnerFirstName,
			LastName:  n.OwnerLastName,
			Thumbnail: n.OwnerThumbnail,
		},
		Status:    openapi.ModelsNoteStatus(n.Note.Status),
		Sections:  sections,
		CreatedAt: n.Note.CreatedAt,
		UpdatedAt: n.Note.UpdatedAt,
	}
}
