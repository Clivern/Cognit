CREATE TABLE agent_protection_tokens (
	id UUID PRIMARY KEY,
	protection_id UUID NOT NULL REFERENCES agent_protection(id) ON DELETE CASCADE,
	token_hash VARCHAR(128) NOT NULL UNIQUE,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP NULL,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_agent_protection_tokens_protection_id ON agent_protection_tokens(protection_id);
CREATE INDEX idx_agent_protection_tokens_expires_at ON agent_protection_tokens(expires_at);
