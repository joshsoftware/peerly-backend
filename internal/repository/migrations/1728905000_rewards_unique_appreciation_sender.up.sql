DELETE FROM rewards
WHERE id NOT IN (
    SELECT MIN(id)
    FROM rewards
    GROUP BY appreciation_id, sender
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_rewards_appreciation_sender ON rewards (appreciation_id, sender);
