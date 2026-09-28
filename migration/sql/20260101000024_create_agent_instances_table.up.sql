CREATE TABLE agent_instances (
	id UUID PRIMARY KEY,
	agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
	instance_id VARCHAR(100) NOT NULL,
	address VARCHAR(255) NOT NULL,
	port INTEGER NOT NULL,
	datacenter VARCHAR(60),
	meta JSONB,
	status VARCHAR(20) NOT NULL DEFAULT 'passing',
	lease_expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (agent_id, instance_id)
);
CREATE INDEX idx_agent_instances_lease_expires_at ON agent_instances(lease_expires_at);
