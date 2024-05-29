--
--
--
-- UUID Extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
--
--
-- Account Table
CREATE TABLE account (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT account_pk PRIMARY KEY (id),
    CONSTRAINT account_ak_id UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT account_ak_email UNIQUE (email,deleted_at) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX account_idx_id ON account (id);
CREATE INDEX account_idx_email ON account (email);

--
--
--
-- Session Table
CREATE TABLE account_session (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    user_agent TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT session_pk PRIMARY KEY (id),
    CONSTRAINT session_fk_account_id FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT session_ak_token UNIQUE (refresh_token_hash,deleted_at) NOT DEFERRABLE INITIALLY IMMEDIATE
);

--
--
--
-- Organization Table
CREATE TABLE organization (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT organization_pk PRIMARY KEY (id),
    CONSTRAINT organization_ak_id UNIQUE (id) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX organization_idx_id ON organization (id);

--
--
--
-- Organization Account Table
CREATE TYPE ROLE AS ENUM ('SUPER_ADMIN', 'ADMIN', 'EMPLOYEE');
CREATE TABLE organization_account (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    account_id UUID NOT NULL,
    role ROLE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT organization_account_pk PRIMARY KEY (id),
    CONSTRAINT organization_account_fk_organization_id FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT organization_account_fk_account_id FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT organization_account_ak_organization_id_account_id UNIQUE (organization_id,account_id,deleted_at) NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX organization_account_idx_id ON organization_account (id);
CREATE INDEX organization_account_idx_organization_id ON organization_account (organization_id);
CREATE INDEX organization_account_idx_account_id ON organization_account (account_id);

--
--
--
-- Event
CREATE TABLE event (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    account_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT event_pk PRIMARY KEY (id),
    CONSTRAINT event_fk_organization_id FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT event_fk_account_id FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX event_idx_id ON event (id);
CREATE INDEX event_idx_organization_id ON event (organization_id);

--
--
--
-- Comment
CREATE TABLE comment (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    account_id UUID NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT event_comment_pk PRIMARY KEY (id),
    CONSTRAINT event_comment_fk_event_id FOREIGN KEY (event_id) REFERENCES event (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE,
    CONSTRAINT event_comment_fk_account_id FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE
);