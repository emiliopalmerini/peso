-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_weights_user_id ON weights(user_id);
CREATE INDEX IF NOT EXISTS idx_weights_measured_at ON weights(measured_at);
CREATE INDEX IF NOT EXISTS idx_weights_user_measured ON weights(user_id, measured_at);
CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_active ON goals(active);
CREATE INDEX IF NOT EXISTS idx_goals_user_active ON goals(user_id, active)