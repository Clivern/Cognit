CREATE TABLE traffic (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	source_instance VARCHAR(100) NOT NULL,
	destination_instance VARCHAR(100) NOT NULL,
	skill VARCHAR(120) NOT NULL,
	verb VARCHAR(60) NOT NULL,
	decision VARCHAR(20) NOT NULL,
	result VARCHAR(40) NOT NULL,
	latency_ms INTEGER NOT NULL,
	task_id VARCHAR(100),
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_traffic_workspace_id_created_at
	ON traffic (workspace_id, created_at DESC);
CREATE INDEX idx_traffic_workspace_id_instances
	ON traffic (workspace_id, source_instance, destination_instance);
