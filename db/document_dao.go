package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
)

type DocumentDAO struct {
	cm *ConnectionManager
}

type DocumentPage struct {
	Documents      []model.Document
	TotalDocuments int
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
		log.Error().Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO documents
    	(id, title, location, date, s3_key, type)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		document.ID.String(), document.Title,
		document.Location, document.Date,
		document.S3Key, document.Type)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert document into documents table")
		return err
	}

	rows := [][]any{}
	for _, a := range authorships {
		log.Debug().Msgf("authorship docId: %s, personId: %s, role: %s", a.DocumentID, a.PersonID, a.Role)
		rows = append(rows, []any{a.PersonID, a.DocumentID, a.Role})
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"authorship"},
		[]string{"person_id", "document_id", "role"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update authorship table")
		return err
	}
	if int(copyCount) != len(authorships) {
		return errs.ErrDB
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO ownership
		(user_id, document_id, role)
		VALUES ($1, $2, $3)`,
		owner.String(), document.ID.String(), "owner")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update ownership table")
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (dao *DocumentDAO) ListDocuments(filter *model.ListDocumentsFilter) (*DocumentPage, error) {

	// 	var documents DocumentPage

	//	dao.cm.DB.Exec(context.Background(),
	//		`SELECT ownership.document_id, documents.title,
	//			documents.s3key, documents.date, documents.type,
	//			author, role FROM ownership
	//		JOIN documents ON ownership.document_id = documents.id
	//		JOIN authorship ON ownership.document_id = authorship.document_id
	//		WHERE (user_id = $1 AND )`
	//
	// )
	return nil, nil
}
