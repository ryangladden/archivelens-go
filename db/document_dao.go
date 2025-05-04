package db

import (
	"context"

	"github.com/google/uuid"
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
	dao.cm.DB.Exec(`
    INSERT INTO ownership
    (user_id, document_id, role)
    VALUES $1, $2, $3`, userId, documentId, role)
}

func (dao *DocumentDAO) AddAuthorship(personId uuid.UUID, documentId uuid.UUID, role string) {

}

func (dao *DocumentDAO) CreateDocument(userId uuid.UUID, document *model.Document) error {

	ctx := context.Background()

	tx, err := dao.cm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.ExecContext(ctx,
		`INSERT INTO documents
    	(id, title, location, date, )
		VALUES $1, $2, $3, $4 $5`,
		document.ID.String(), document.Title,
		document.Location, document.Date,
		document.S3Key)

	// tx.ExecContext(ctx,
	// 	`INSERT INTO authorship
	// 	(person_id, document_id, role_enum)
	// 	VALUES $1, $2, $3
	// 	`)
	return nil
}
