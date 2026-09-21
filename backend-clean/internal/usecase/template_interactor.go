// Package usecase contains application use case implementations.
package usecase

import (
	"context"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/template"
	"immortal-architecture-cqrs/backend/internal/port"
)

// TemplateInteractor handles template use cases.
type TemplateInteractor struct {
	repo   port.TemplateRepository
	tx     port.TxManager
	output port.TemplateOutputPort
}

var _ port.TemplateInputPort = (*TemplateInteractor)(nil)

// NewTemplateInteractor creates TemplateInteractor.
func NewTemplateInteractor(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) *TemplateInteractor {
	return &TemplateInteractor{repo: repo, tx: tx, output: output}
}

// List returns templates by filters.
func (u *TemplateInteractor) List(ctx context.Context, filters template.Filters) error {
	templates, err := u.repo.List(ctx, filters)
	if err != nil {
		return err
	}
	return u.output.PresentTemplateList(ctx, templates)
}

// Get returns template by ID.
func (u *TemplateInteractor) Get(ctx context.Context, id string) error {
	tpl, err := u.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	return u.output.PresentTemplate(ctx, tpl)
}

// Create creates a template.
func (u *TemplateInteractor) Create(ctx context.Context, input port.TemplateCreateInput) error {
	if err := template.ValidateTemplate(template.Template{
		Name:    input.Name,
		OwnerID: input.OwnerID,
		Fields:  input.Fields,
	}); err != nil {
		return err
	}

	var createdID string
	err := u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		tpl, err := u.repo.Create(txCtx, template.Template{
			Name:    input.Name,
			OwnerID: input.OwnerID,
		})
		if err != nil {
			return err
		}
		createdID = tpl.ID
		if len(input.Fields) > 0 {
			if err := u.repo.ReplaceFields(txCtx, tpl.ID, input.Fields); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	tpl, err := u.repo.Get(ctx, createdID)
	if err != nil {
		return err
	}
	return u.output.PresentTemplate(ctx, tpl)
}

// Update updates a template.
func (u *TemplateInteractor) Update(ctx context.Context, input port.TemplateUpdateInput) error {
	current, err := u.repo.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	if err := template.ValidateTemplateOwnership(current.Template.OwnerID, input.OwnerID); err != nil {
		return err
	}
	if input.Fields != nil {
		if err := template.ValidateTemplate(template.Template{
			ID:      input.ID,
			Name:    input.Name,
			Fields:  input.Fields,
			OwnerID: input.OwnerID,
		}); err != nil {
			return err
		}
	}
	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, err := u.repo.Update(txCtx, template.Template{
			ID:   input.ID,
			Name: input.Name,
		})
		if err != nil {
			return err
		}
		if input.Fields != nil {
			if len(input.Fields) == 0 {
				return domainerr.ErrInvalidTemplateField
			}
			if err := u.repo.ReplaceFields(txCtx, input.ID, input.Fields); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	tpl, err := u.repo.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	return u.output.PresentTemplate(ctx, tpl)
}

// Delete deletes a template.
func (u *TemplateInteractor) Delete(ctx context.Context, id, ownerID string) error {
	tpl, err := u.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := template.ValidateTemplateOwnership(tpl.Template.OwnerID, ownerID); err != nil {
		return err
	}
	if err := template.CanDeleteTemplate(tpl.IsUsed); err != nil {
		return err
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}
	return u.output.PresentTemplateDeleted(ctx)
}
