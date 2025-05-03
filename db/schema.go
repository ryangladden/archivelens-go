package db

import (
	"fmt"
)

func (cm *ConnectionManager) Init() error {

	err := createUpdatedAtFunction(cm)
	if err != nil {
		return fmt.Errorf("error creating updated_at function: %w", err)
	}

	err = createDocumentTable(cm)
	if err != nil {
		return fmt.Errorf("error creating documents table: %w", err)
	}

	err = createPersonsTable(cm)
	if err != nil {
		return fmt.Errorf("error creating persons table: %w", err)
	}

	err = createUsersTable(cm)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	err = createOwnershipTable(cm)
	if err != nil {
		return fmt.Errorf("error creating ownership table: %w", err)
	}

	err = createAuthorshipTable(cm)
	if err != nil {
		return fmt.Errorf("error creating authorship table: %w", err)
	}

	err = createTagsTable(cm)
	if err != nil {
		return fmt.Errorf("error creating tags table: %w", err)
	}

	err = createTaggingTable(cm)
	if err != nil {
		return fmt.Errorf("error creating tagging table: %w", err)
	}

	err = createAuthTable(cm)
	if err != nil {
		return fmt.Errorf("error creating auth table: %w", err)
	}

	return nil
}

func createDocumentTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS documents (
		id uuid NOT NULL,
		title TEXT NOT NULL,
		date TEXT,
		location TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		return fmt.Errorf("error creating documents table: %w", err)
	}

	err = createUpdatedAtTrigger(cm, "documents")

	return err
}

func createPersonsTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS persons (
		id uuid NOT NULL,
		name TEXT NOT NULL,
		metadata JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		return fmt.Errorf("error creating persons table: %w", err)
	}

	err = createUpdatedAtTrigger(cm, "persons")

	return err
}

func createUsersTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id uuid NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password BYTEA NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (id)
	)`)

	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	err = createUpdatedAtTrigger(cm, "users")

	return err
}

func createOwnershipTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`DO $$ BEGIN
			CREATE TYPE role_enum  AS ENUM 
			('owner', 'editor', 'viewer');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)

	if err != nil {
		return fmt.Errorf("error creating role_enum: %w", err)
	}

	_, err = cm.DB.Exec(`CREATE TABLE IF NOT EXISTS ownership (
		user_id uuid NOT NULL,
		document_id uuid NOT NULL,
		role role_enum NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		PRIMARY KEY (user_id, document_id),
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (document_id) REFERENCES documents (id)
	)`)

	if err != nil {
		return fmt.Errorf("error creating ownership table: %w", err)
	}

	err = createUpdatedAtTrigger(cm, "ownership")

	return err
}

func createAuthorshipTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`DO $$ BEGIN
		CREATE TYPE
			authorship_enum AS ENUM
			('author', 'subject', 'recipient');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)
	if err != nil {
		return fmt.Errorf("error creating authorship_enum: %w", err)
	}

	_, err = cm.DB.Exec(`CREATE TABLE IF NOT EXISTS authorship (
		person_id uuid NOT NULL,
		document_id uuid NOT NULL,
		role authorship_enum NOT NULL,
		PRIMARY KEY (person_id, document_id),
		FOREIGN KEY (person_id) REFERENCES persons (id),
		FOREIGN KEY (document_id) REFERENCES documents (id)
	)`)
	return err
}

func createTagsTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS tags (
		id SERIAL NOT NULL,
		tag TEXT NOT NULL UNIQUE,
		PRIMARY KEY (id)
	)`)

	return err
}

func createTaggingTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS document_tags (
		document_id uuid NOT NULL,
		tag_id SERIAL NOT NULL,
		PRIMARY KEY (document_id, tag_id),
		FOREIGN KEY (document_id) REFERENCES documents (id),
		FOREIGN KEY (tag_id) REFERENCES tags (id)
	)`)

	return err
}

func createAuthTable(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE TABLE IF NOT EXISTS auth (
		token uuid NOT NULL,
		user_id uuid NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		expires_at TIMESTAMP WITH TIME ZONE DEFAULT now() + interval '1 day' * 180,
		PRIMARY KEY (token),
		FOREIGN KEY (user_id) REFERENCES users (id)
		)`)

	if err != nil {
		return fmt.Errorf("error creating documents table: %w", err)
	}

	err = createUpdatedAtTrigger(cm, "documents")

	return err
}

func createUpdatedAtFunction(cm *ConnectionManager) error {
	_, err := cm.DB.Exec(`CREATE OR REPLACE FUNCTION
	update_updated_at_column()
	RETURNS TRIGGER AS
	$$ BEGIN
	NEW.updated_at = now();
	RETURN NEW;
	END; $$ language 'plpgsql';`)

	return err
}

func createUpdatedAtTrigger(cm *ConnectionManager, table string) error {
	_, err := cm.DB.Exec(`CREATE  OR REPLACE TRIGGER
	supdate_updated_at
	BEFORE UPDATE ON ` + table + `
	FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();`)

	if err != nil {
		return fmt.Errorf("error creating updated_at trigger: %w", err)
	}
	return nil
}
