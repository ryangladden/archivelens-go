package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ryangladden/archivelens-go/model"
)

type PersonDAO struct {
	cm *ConnectionManager
}

func NewPersonDAO(cm *ConnectionManager) *PersonDAO {
	return &PersonDAO{
		cm: cm,
	}
}

func (dao *PersonDAO) CreatePerson(person *model.Person, owner uuid.UUID) error {

	ctx := context.Background()

	tx, err := dao.cm.DB.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO persons
		(id, name, s3_key, birth, death, summary)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		person.ID.String(), person.Name, formatKey(person.S3Key),
		person.Birth, person.Death, person.Summary,
	)
	if err != nil {
		log.Error().Err(err).Msgf("Error inserting person %s into persons table", person.Name)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO users_persons
		(user_id, person_id, role)
		VALUES ($1, $2, $3)`,
		owner.String(), person.ID, "owner",
	)
	if err != nil {
		log.Error().Err(err).Msgf("Error adding owner %s to new person %s", owner.String(), person.Name)
		return err
	}
	return nil
}

func formatKey(key string) *string {
	if key == "" {
		return nil
	}
	return &key
}
