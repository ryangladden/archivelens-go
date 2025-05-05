package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
)

type DocumentDAO struct {
	cm *ConnectionManager
}

func NewDocumentDAO(cm *ConnectionManager) *DocumentDAO {
	return &DocumentDAO{
		cm: cm,
	}
}

func (dao *DocumentDAO) AddOwnership(userId uuid.UUID, documentId uuid.UUID, role string) {
	dao.cm.DB.Exec(context.Background(),
		`INSERT INTO ownership
    	(user_id, document_id, role)
    	VALUES $1, $2, $3`,
		userId, documentId, role)
}

func (dao *DocumentDAO) AddAuthorship(personId uuid.UUID, documentId uuid.UUID, role string) {

}

func (dao *DocumentDAO) CreateDocument(owner uuid.UUID, document *model.Document, authorships []model.Authorship) error {

	ctx := context.Background()

	tx, err := dao.cm.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO documents
    	(id, title, location, date, )
		VALUES $1, $2, $3, $4 $5`,
		document.ID.String(), document.Title,
		document.Location, document.Date,
		document.S3Key)
	if err != nil {
		return err
	}

	rows := [][]any{}
	for i, a := range authorships {
		rows[i] = []any{a.PersonID, a.DocumentID, a.Role}
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"authorship"},
		[]string{"person_id", "document_id", "role"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {

	}
	if int(copyCount) != len(authorships) {
		return errs.ErrDB
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO ownership
		(user_id, document_id, role)
		VALUES $1, $2, $3`,
		owner.String, document.ID, "author")
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
