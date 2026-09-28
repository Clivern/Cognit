CREATE TABLE agent_instances_meta (
	id UUID PRIMARY KEY,
	agent_instance_id UUID NOT NULL REFERENCES agent_instances(id) ON DELETE CASCADE,
	key VARCHAR(60) NOT NULL,
	value JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (agent_instance_id, key)
);
