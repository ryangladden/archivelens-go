package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func Init(db *pgx.Conn) {

	createUpdatedAtFunction(db)
	createDocumentTable(db)
	createPersonsTable(db)
	createUsersTable(db)
	createOwnershipTable(db)
	createAuthorshipTable(db)
	createTagsTable(db)
	createTaggingTable(db)
	createAuthTable(db)
	createUsersPersonsTable(db)
	createDocumentStatusTable(db)
}

func createDocumentTable(db *pgx.Conn) {

	_, err := db.Exec(context.Background(), `DO $$ BEGIN
		CREATE TYPE document_type AS ENUM 
			('letter', 'journal', 'audio', 'email', 'other');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)
	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create document_type enum")
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS documents (
		id uuid NOT NULL,
		title TEXT NOT NULL,
		date DATE,
		location TEXT,
		type document_type NOT NULL,
		original_filename TEXT NOT NULL,
		pages SMALLINT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create documents table")
	}
	createUpdatedAtTrigger(db, "documents")
}

func createPersonsTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS persons (
		id uuid NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		birth DATE,
		death DATE,
		summary TEXT,
		s3_key TEXT,
		metadata JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create persons table")
	}
	createUpdatedAtTrigger(db, "persons")
}

func createUsersTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id uuid NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password BYTEA NOT NULL,
		s3_key TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create users table")
	}

	createUpdatedAtTrigger(db, "users")
}

func createOwnershipTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `DO $$ BEGIN
		CREATE TYPE role_enum AS ENUM 
			('owner', 'editor', 'viewer');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create ownership role_enum enum")
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS ownership (
		user_id uuid NOT NULL,
		document_id uuid NOT NULL,
		role role_enum NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (user_id, document_id),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE
	)`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create ownership table")
	}
	createUpdatedAtTrigger(db, "ownership")
}

func createAuthorshipTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `DO $$ BEGIN
		CREATE TYPE authorship_enum AS ENUM
			('author', 'coauthor', 'mentioned', 'recipient');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)
	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create authorship_enum enum")
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS authorship (
		person_id uuid NOT NULL,
		document_id uuid NOT NULL,
		role authorship_enum NOT NULL,
		PRIMARY KEY (person_id, document_id),
		FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE,
		FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create authorship table")
	}
}

func createTagsTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS tags (
		id SERIAL NOT NULL,
		tag TEXT NOT NULL UNIQUE,
		PRIMARY KEY (id)
	)`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create tags table")
	}
}

func createTaggingTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS document_tags (
		document_id uuid NOT NULL,
		tag_id SERIAL NOT NULL,
		PRIMARY KEY (document_id, tag_id),
		FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
	)`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create tagging table")
	}
}

func createAuthTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS auth (
		token uuid NOT NULL,
		user_id uuid NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		expires_at TIMESTAMP WITH TIME ZONE DEFAULT now() + interval '1 day' * 180,
		PRIMARY KEY (token),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create auth table")
	}

	createUpdatedAtTrigger(db, "auth")
}

func createUsersPersonsTable(db *pgx.Conn) {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users_persons (
		user_id uuid NOT NULL,
		person_id uuid NOT NULL,
		role role_enum NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
		)`)
	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create users_persons table")
	}
}

func createUpdatedAtFunction(db *pgx.Conn) {
	_, err := db.Exec(context.Background(), `CREATE OR REPLACE FUNCTION
	update_updated_at_column()
	RETURNS TRIGGER AS
	$$ BEGIN
	NEW.updated_at = now();
	RETURN NEW;
	END; $$ language 'plpgsql';`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create updated_at function")
	}
}

func createUpdatedAtTrigger(db *pgx.Conn, table string) {
	_, err := db.Exec(context.Background(), `CREATE OR REPLACE TRIGGER
	update_updated_at
	BEFORE UPDATE ON `+table+`
	FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();`)

	if err != nil {
		log.Fatal().Err(err).Msgf("DB initialization failed to create updated_at trigger for %s table", table)
	}
}

func createDocumentStatusTable(db *pgx.Conn) {

	_, err := db.Exec(context.Background(), `DO $$ BEGIN
		CREATE TYPE job_status AS ENUM 
			('pending', 'processing', 'processed', 'failed');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)
	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create job_status enum")
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS document_status (
		document_id uuid NOT NULL,
		thumbnail job_status DEFAULT 'pending',
		preview job_status DEFAULT 'pending',
		FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE
		)`)
	if err != nil {
		log.Fatal().Err(err).Msg("DB initialization failed to create document_status table")
	}
}

// func createIndex(db *pgx.Conn, table string, column string) {
// 	insert := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS ON %s(%s)`, table, column)
// 	_, err := db.Exec(context.Background(), insert)
// 	if err != nil {
// 		log.Fatal().Err(err).Msgf("DB initialization failed to create index for column '%s' in table '%s'")
// 	}
// }
