// Generated by sqlboiler-erg: DO NOT EDIT.
package erg

import (
	"github.com/wearepointers/tycho/example/models/dm"
	"time"
)

type Event struct {
	ID             string     `json:"id,omitempty" toml:"id" yaml:"id"`
	OrganizationID string     `json:"organizationId,omitempty" toml:"organization_id" yaml:"organization_id"`
	AccountID      string     `json:"accountId,omitempty" toml:"account_id" yaml:"account_id"`
	Name           string     `json:"name,omitempty" toml:"name" yaml:"name"`
	Description    string     `json:"description,omitempty" toml:"description" yaml:"description"`
	CreatedAt      time.Time  `json:"createdAt,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty" toml:"updated_at" yaml:"updated_at"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" toml:"deleted_at" yaml:"deleted_at"`

	Account      *Account      `json:"account,omitempty" toml:"account" yaml:"account"`
	Organization *Organization `json:"organization,omitempty" toml:"organization" yaml:"organization"`
	Comments     CommentSlice  `json:"comments,omitempty" toml:"comments" yaml:"comments"`

	CustomFields `json:"customFields,omitempty" toml:"custom_fields" yaml:"custom_fields"`
}

type EventSlice []*Event

func ToEvents(a dm.EventSlice, acf CustomFieldsSlice, exclude ...string) EventSlice {
	s := make(EventSlice, len(a))
	for i, d := range a {
		var cf CustomFields
		if acf != nil {
			if value, ok := acf[d.ID]; ok {
				cf = value
			}
		}

		s[i] = ToEvent(d, cf, exclude...)
	}
	return s
}

func ToEvent(a *dm.Event, customFields CustomFields, exclude ...string) *Event {
	p := Event{
		ID:             a.ID,
		OrganizationID: a.OrganizationID,
		AccountID:      a.AccountID,
		Name:           a.Name,
		Description:    a.Description,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		DeletedAt:      nullDotTimeToTimePtr(a.DeletedAt),
	}

	if a.R != nil {
		if a.R.Account != nil && doesNotContain(exclude, "event.account") {
			p.Account = ToAccount(a.R.Account, nil, append(exclude, "account.event")...)
		}
		if a.R.Organization != nil && doesNotContain(exclude, "event.organization") {
			p.Organization = ToOrganization(a.R.Organization, nil, append(exclude, "organization.event")...)
		}
		if a.R.Comments != nil && doesNotContain(exclude, "event.comment") {
			p.Comments = ToComments(a.R.Comments, nil, append(exclude, "comment.event")...)
		}
	}

	if customFields != nil {
		p.CustomFields = customFields
	}

	return &p
}
