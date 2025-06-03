package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
)

type DocumentDAO struct {
	cm *ConnectionManager
}

type DocumentPage struct {
	Documents      []InlineDocument
	TotalDocuments int
}

type InlineDocument struct {
	Document         model.Document
	DocumentMetadata DocumentMetadata
}

type DocumentMetadata struct {
	Tags    []model.Tag    `json:"tags"`
	Persons []model.Person `json:"persons"`
	Author  model.Person   `json:"author"`
}

type InlinePerson struct {
	ID        pgtype.UUID `json:"id"`
	FirstName pgtype.Text `json:"first_name"`
	LastName  pgtype.Text `json:"last_name"`
	Role      pgtype.Text `json:"role"`
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

	var page DocumentPage
	// var documents DocumentPage
	countQuery, query := dao.generateQuery(filter)

	log.Debug().Msgf("Used the following query: \n%s", query)
	log.Debug().Msgf("Counting rows with the following query: \n%s", countQuery)

	dao.cm.DB.QueryRow(context.Background(), countQuery, filter.UserID.String()).Scan(&page.TotalDocuments)
	log.Debug().Msgf("Total rows: %d", page.TotalDocuments)
	rows, err := dao.cm.DB.Query(context.Background(), query, filter.UserID.String(), filter.Limit, filter.Limit*filter.Page)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to list documents for user %s", filter.UserID.String())
		return nil, errs.ErrDB
	}
	page.Documents = readDocumentListRows(rows)
	return &page, nil
}

func (dao *DocumentDAO) GetDocument(userID uuid.UUID, documentID uuid.UUID) (*model.Document, error) {
	var document model.Document

	err := dao.cm.DB.QueryRow(context.Background(),
		`WITH users_documents AS (
   			SELECT document_id AS id, role
   			FROM ownership
   			WHERE user_id = $1
      			UNION
   			SELECT a.document_id AS id, up.role
   			FROM users_persons up
   			JOIN authorship a ON a.person_id = up.person_id
   			WHERE up.user_id = $1
  		)
		SELECT d.id, d.title, d.date, d.location, d.type, d.s3_key, MIN(ud.role) AS permissions
		FROM users_documents ud
		JOIN documents d ON ud.id = d.id
		WHERE ud.id = $2
		GROUP BY d.id`, userID.String(), documentID.String()).Scan(&document.ID, &document.Title, &document.Date, &document.Location, &document.Type, &document.S3Key, &document.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Msgf("Either document id %s does not exist or user %s does not have permissions to access it", documentID.String(), userID.String())
			return nil, errs.ErrNotFound
		}
		log.Error().Err(err).Msgf("Error finding document with id %s in database", documentID.String())
		return nil, errs.ErrDB
	}
	log.Debug().Msgf("Searching DB for persons associated with document ID %s", documentID.String())
	personsRows, err := dao.cm.DB.Query(context.Background(),
		`SELECT p.id, p.first_name, p.last_name, a.role, p.s3_key
		FROM persons p
		JOIN authorship a ON p.id = a.person_id AND a.document_id = $1`, documentID.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Msgf("No persons found associated with document %s", documentID.String())
		} else {
			log.Error().Err(err).Msgf("Error finding persons associated with document id %s in database", documentID.String())
			return nil, errs.ErrDB
		}
	} else {
		addDocumentAuthorship(personsRows, &document)
	}

	tagRows, err := dao.cm.DB.Query(context.Background(),
		`SELECT t.tag, t.id
		FROM tags t
		JOIN document_tags dt ON t.id = dt.tag_id AND dt.document_id = $1`, documentID.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Msgf("No tags found associated with document %s", documentID.String())
		} else {
			log.Error().Err(err).Msgf("Error finding tags associated with document id %s in database", documentID.String())
			return nil, errs.ErrDB
		}
	} else {
		document.Tags = readTagsRows(tagRows)
	}
	log.Debug().Msgf("Document of ID %s:\nTitle: %s\n", documentID.String(), document.Title)
	return &document, nil
}

func (dao *DocumentDAO) personsTagsCTE(filter *model.ListDocumentsFilter) (string, string, string, string) {
	if filter.Authors == nil && filter.IncludeTags == nil {
		return "", "JSONB_BUILD_OBJECT()", "", ""
	}
	tags := ""
	persons := ""
	union := "UNION ALL"
	if filter.IncludeTags != nil {
		tags = fmt.Sprintf(
			`SELECT dt.document_id AS id, JSONB_BUILD_OBJECT('id', t.id, 'tag', t.tag) AS tag, NULL AS persons
			FROM document_tags dt
			JOIN tags t on dt.tag_id = t.id
            WHERE dt.tag_id IN (%s)`, *filter.IncludeTags)
	}
	if filter.Authors != nil {
		persons = fmt.Sprintf(
			`SELECT a.document_id AS id, NULL as tag,
            JSONB_BUILD_OBJECT('id', p.id, 'name', p.first_name, p.last_name, 'role', a.role) AS persons
 			FROM authorship a
 			JOIN persons p ON a.person_id = p.id
            WHERE a.person_id IN (%s)`, *filter.Authors)
	}
	if filter.IncludeTags == nil || filter.Authors == nil {
		union = ""
	}
	return fmt.Sprintf(
		`,
		filter AS (
			WITH metadata AS (%s %s %s)
			SELECT id, ARRAY_AGG(tag) FILTER (WHERE tag IS NOT NULL) AS tags,
                 ARRAY_AGG(persons) FILTER (WHERE persons IS NOT NULL) AS persons
            FROM metadata GROUP BY id
		)`, tags, union, persons), "JSONB_BUILD_OBJECT('tags', f.tags, 'persons', f.persons)", "JOIN filter f on d.id = f.id", ", f.tags, f.persons"
}

func (dao *DocumentDAO) generateQuery(filter *model.ListDocumentsFilter) (string, string) {
	personsTags, jsonBuild, joinFilter, groupBy := dao.personsTagsCTE(filter)
	where := dao.generateWhere(filter)
	cte := fmt.Sprintf(`
	WITH document_list AS (
		WITH users_documents AS (
			SELECT document_id AS id, role
			FROM ownership
			WHERE user_id = $1
			UNION
			SELECT a.document_id AS id, up.role
			FROM users_persons up
			JOIN authorship a ON a.person_id = up.person_id
			WHERE up.user_id = $1
		)%s -- personsTagsCTE(filter)
		SELECT d.id, d.title, d.date, d.type, %s AS metadata, MIN(ud.role) AS permissions -- JSONB_BUILD_OBJECT('tags', f.tags, 'persons', f.persons) or JSONB_BUILD_OBJECT()
		FROM documents d
		JOIN users_documents ud ON d.id = ud.id
		%s -- JOIN filter f on d.id = f.id
		GROUP BY d.id %s
		)
		`, personsTags, jsonBuild, joinFilter, groupBy)
	count := fmt.Sprintf(`SELECT COUNT(*) FROM document_list dl %s`, where)
	query := fmt.Sprintf(`
	SELECT dl.id, dl.title, dl.date, dl.type, dl.metadata, dl.permissions, p.id AS author_id, p.first_name AS author_first_name, p.last_name AS author_last_name
    FROM document_list dl -- order by, asc or desc
    LEFT JOIN authorship a ON dl.id = a.document_id AND a.role = 'author'
    LEFT JOIN persons p ON p.id = a.person_id
    %s
	ORDER BY %s %s -- order by, asc or desc
	LIMIT $2 OFFSET $3`,
		where, filter.SortBy, filter.Order)
	return cte + count, cte + query
}

func (dao *DocumentDAO) generateWhere(filter *model.ListDocumentsFilter) string {
	var conditions []string
	if filter.ExcludeRoles != nil {
		conditions = append(conditions, fmt.Sprintf("dl.permissions NOT IN (%s)", *filter.ExcludeRoles))
	}
	if filter.DateMax != nil {
		conditions = append(conditions, fmt.Sprintf("dl.date <= %s", *filter.DateMax))
	}
	if filter.DateMin != nil {
		conditions = append(conditions, fmt.Sprintf("dl.date >= %s", *filter.DateMin))
	}
	if filter.ExcludeType != nil {
		conditions = append(conditions, fmt.Sprintf("dl.type NOT IN (%s)", *filter.ExcludeType))
	}
	if filter.TitleMatch != nil {
		log.Debug().Msgf("Title match: %s", *filter.TitleMatch)
		conditions = append(conditions, fmt.Sprintf("dl.title ILIKE '%%%s%%'", *filter.TitleMatch))
	}
	where := strings.Join(conditions, " OR ")
	if where != "" {
		return "WHERE " + where
	}
	return ""
}

func readDocumentListRows(rows pgx.Rows) []InlineDocument {
	var documents []InlineDocument
	for rows.Next() {
		var document InlineDocument
		if err := rows.Scan(&document.Document.ID, &document.Document.Title, &document.Document.Date, &document.Document.Type, &document.DocumentMetadata, &document.Document.Role, &document.DocumentMetadata.Author.ID, &document.DocumentMetadata.Author.FirstName, &document.DocumentMetadata.Author.LastName); err != nil {
			log.Error().Err(err).Msg("Failed to scan row in document list")
			continue
		}
		log.Debug().Msgf("ID: %s, Title: %s, Date: %s, Type: %s, Role: %s, Author: %s", document.Document.ID.String(), document.Document.Title, document.Document.Date, document.Document.Type, document.Document.Role, document.DocumentMetadata.Author.FirstName, &document.DocumentMetadata.Author.LastName)
		documents = append(documents, document)
	}
	return documents
}

func readInlinePersonRows(rows pgx.Rows) []model.Person {
	var persons []model.Person
	for rows.Next() {
		var person model.Person
		var s3key pgtype.Text
		if err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Role, &s3key); err != nil {
			log.Error().Err(err).Msg("Failed to scan row in persons list")
			continue
		}
		if s3key.Status != pgtype.Null {
			person.S3Key = &s3key.String
		}
		log.Debug().Msgf("PERSON: ID: %s, Authorship Role: %s", person.ID.String(), *person.Role)
		persons = append(persons, person)
	}
	return persons
}

func addDocumentAuthorship(rows pgx.Rows, document *model.Document) {
	coauthors := []model.Person{}
	mentions := []model.Person{}
	persons := readInlinePersonRows(rows)
	for _, person := range persons {
		switch *person.Role {
		case "author":
			document.Author = &person
		case "coauthor":
			coauthors = append(coauthors, person)
		case "subject":
			mentions = append(mentions, person)
		case "recipient":
			document.Recipient = &person
		}
	}
	document.Coauthors = &coauthors
	document.Mentions = &mentions
}

func readTagsRows(rows pgx.Rows) *[]model.Tag {
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Tag); err != nil {
			log.Error().Err(err).Msg("Failed to scan row in tag list")
			continue
		}
		log.Debug().Msgf("TAG: ID: %d, Tag: %s", tag.ID, tag.Tag)
		tags = append(tags, tag)
	}
	return &tags
}
