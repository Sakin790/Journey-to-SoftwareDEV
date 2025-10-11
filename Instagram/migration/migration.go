package migration

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS post_likes (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		post_id BIGINT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		
		-- Constraint to prevent double-liking
		CONSTRAINT unique_like UNIQUE (user_id, post_id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes (post_id);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully ✅")
}
