// Package template holds template domain models.
package template

// Filters for listing templates.
type Filters struct {
	Query   *string
	OwnerID *string
}

// Owner holds minimal owner info for embedding.
type Owner struct {
	ID        string
	FirstName string
	LastName  string
	Thumbnail *string
}

// WithUsage is a template with usage metadata.
type WithUsage struct {
	Template Template
	IsUsed   bool
	Owner    Owner
}
