// Generated by sqlboiler-erg: DO NOT EDIT.
package erg

import (
	"github.com/wearepointers/tycho/example/models/dm"
	"time"
)

type OrganizationAccount struct {
	ID             string     `json:"id,omitempty" toml:"id" yaml:"id"`
	OrganizationID string     `json:"organizationId,omitempty" toml:"organization_id" yaml:"organization_id"`
	AccountID      string     `json:"accountId,omitempty" toml:"account_id" yaml:"account_id"`
	Role           Role       `json:"role,omitempty" toml:"role" yaml:"role"`
	CreatedAt      time.Time  `json:"createdAt,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty" toml:"updated_at" yaml:"updated_at"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" toml:"deleted_at" yaml:"deleted_at"`

	Account      *Account      `json:"account,omitempty" toml:"account" yaml:"account"`
	Organization *Organization `json:"organization,omitempty" toml:"organization" yaml:"organization"`

	CustomFields `json:"customFields,omitempty" toml:"custom_fields" yaml:"custom_fields"`
}

type OrganizationAccountSlice []*OrganizationAccount

func ToOrganizationAccounts(a dm.OrganizationAccountSlice, acf CustomFieldsSlice, exclude ...string) OrganizationAccountSlice {
	s := make(OrganizationAccountSlice, len(a))
	for i, d := range a {
		var cf CustomFields
		if acf != nil {
			if value, ok := acf[d.ID]; ok {
				cf = value
			}
		}

		s[i] = ToOrganizationAccount(d, cf, exclude...)
	}
	return s
}

func ToOrganizationAccount(a *dm.OrganizationAccount, customFields CustomFields, exclude ...string) *OrganizationAccount {
	p := OrganizationAccount{
		ID:             a.ID,
		OrganizationID: a.OrganizationID,
		AccountID:      a.AccountID,
		Role:           Role(a.Role),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		DeletedAt:      nullDotTimeToTimePtr(a.DeletedAt),
	}

	if a.R != nil {
		if a.R.Account != nil && doesNotContain(exclude, "organization_account.account") {
			p.Account = ToAccount(a.R.Account, nil, append(exclude, "account.organization_account")...)
		}
		if a.R.Organization != nil && doesNotContain(exclude, "organization_account.organization") {
			p.Organization = ToOrganization(a.R.Organization, nil, append(exclude, "organization.organization_account")...)
		}
	}

	if customFields != nil {
		p.CustomFields = customFields
	}

	return &p
}