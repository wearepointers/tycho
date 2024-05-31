// Generated by sqlboiler-erg: DO NOT EDIT.
package erg

import (
	"github.com/wearepointers/tycho/example/models/dm"
	"time"
)

type AccountSession struct {
	ID        string     `json:"id,omitempty" toml:"id" yaml:"id"`
	AccountID string     `json:"accountId,omitempty" toml:"account_id" yaml:"account_id"`
	IPAddress string     `json:"ipAddress,omitempty" toml:"ip_address" yaml:"ip_address"`
	UserAgent string     `json:"userAgent,omitempty" toml:"user_agent" yaml:"user_agent"`
	ExpiresAt time.Time  `json:"expiresAt,omitempty" toml:"expires_at" yaml:"expires_at"`
	CreatedAt time.Time  `json:"createdAt,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" toml:"updated_at" yaml:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" toml:"deleted_at" yaml:"deleted_at"`

	Account *Account `json:"account,omitempty" toml:"account" yaml:"account"`

	CustomFields `json:"customFields,omitempty" toml:"custom_fields" yaml:"custom_fields"`
}

type AccountSessionSlice []*AccountSession

func ToAccountSessions(a dm.AccountSessionSlice, acf CustomFieldsSlice, exclude ...string) AccountSessionSlice {
	s := make(AccountSessionSlice, len(a))
	for i, d := range a {
		var cf CustomFields
		if acf != nil {
			if value, ok := acf[d.ID]; ok {
				cf = value
			}
		}

		s[i] = ToAccountSession(d, cf, exclude...)
	}
	return s
}

func ToAccountSession(a *dm.AccountSession, customFields CustomFields, exclude ...string) *AccountSession {
	p := AccountSession{
		ID:        a.ID,
		AccountID: a.AccountID,
		IPAddress: a.IPAddress,
		UserAgent: a.UserAgent,
		ExpiresAt: a.ExpiresAt,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: nullDotTimeToTimePtr(a.DeletedAt),
	}

	if a.R != nil {
		if a.R.Account != nil && doesNotContain(exclude, "account_session.account") {
			p.Account = ToAccount(a.R.Account, nil, append(exclude, "account.account_session")...)
		}
	}

	if customFields != nil {
		p.CustomFields = customFields
	}

	return &p
}