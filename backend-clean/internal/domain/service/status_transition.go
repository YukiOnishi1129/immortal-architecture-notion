// Package service contains domain services for cross-aggregate operations.
package service

import (
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
)

// CanPublish checks if the actor can publish the note.
// ルール: オーナーのみ、Draft -> Publish のみ。
func CanPublish(n note.Note, actorID string) error {
	if n.OwnerID != actorID || actorID == "" {
		return domainerr.ErrUnauthorized
	}
	if err := n.Status.Validate(); err != nil {
		return err
	}
	return note.CanChangeStatus(n.Status, note.StatusPublish)
}

// CanUnpublish checks if the actor can unpublish the note.
// ルール: オーナーのみ、Publish -> Draft のみ。
func CanUnpublish(n note.Note, actorID string) error {
	if n.OwnerID != actorID || actorID == "" {
		return domainerr.ErrUnauthorized
	}
	if err := n.Status.Validate(); err != nil {
		return err
	}
	return note.CanChangeStatus(n.Status, note.StatusDraft)
}
