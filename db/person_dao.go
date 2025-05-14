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

type PersonDAO struct {
	cm *ConnectionManager
}

func NewPersonDAO(cm *ConnectionManager) *PersonDAO {
	return &PersonDAO{
		cm: cm,
	}
}

type PersonPage struct {
	Persons      []model.Person
	TotalPersons int
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
		(id, first_name, last_name, s3_key, birth, death, summary)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		person.ID.String(), person.FirstName, person.LastName, person.S3Key,
		person.Birth, person.Death, person.Summary,
	)
	if err != nil {
		log.Error().Err(err).Msgf("Error inserting person %s %s into persons table", person.FirstName, person.LastName)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO users_persons
		(user_id, person_id, role)
		VALUES ($1, $2, $3)`,
		owner.String(), person.ID, "owner",
	)
	if err != nil {
		log.Error().Err(err).Msgf("Error adding owner %s to new person %s %s", owner.String(), person.FirstName, person.LastName)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		log.Error().Err(err).Msgf("Failed to create person %s %s", person.FirstName, person.LastName)
		return errs.ErrDB
	}
	return nil
}

func (dao *PersonDAO) GetPerson(userID uuid.UUID, personID uuid.UUID) (*model.Person, error) {

	var person model.Person
	var idHolder string
	row := dao.cm.DB.QueryRow(context.Background(),
		`SELECT person_id, first_name, last_name, birth, death, summary, s3_key, role, user_id
		FROM users_persons
		JOIN persons ON users_persons.person_id = persons.id
		WHERE person_id = $1 AND user_id = $2`,
		personID.String(), userID.String())

	err := row.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Birth, &person.Death, &person.Summary, &person.S3Key, &person.Role, &idHolder)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting row for person with ID %s owned by %s", personID.String(), userID.String())
		return nil, errs.ErrDB
	}
	return &person, nil
}

func (dao *PersonDAO) ListPersons(filter *model.ListPersonsFilter) (*PersonPage, error) {

	var personPage PersonPage
	var totalPersons int
	conditions := generateAndConditions(filter)
	countQuery := fmt.Sprintf(`SELECT COUNT(*)
		FROM persons
		JOIN users_persons ON persons.id = users_persons.person_id
		WHERE (user_id = $1%s)`, conditions)

	log.Debug().Msgf("Counting total matches based on query: %s", countQuery)
	dao.cm.DB.QueryRow(context.Background(), countQuery, filter.UserID).Scan(&totalPersons)
	personPage.TotalPersons = totalPersons
	log.Debug().Msgf("Total persons returned: %d", totalPersons)

	listQuery := fmt.Sprintf(`SELECT id, first_name, last_name, birth, death, summary, s3_key, role
		FROM persons
		JOIN users_persons ON persons.id = users_persons.person_id
		WHERE (user_id = $1%s)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, conditions, filter.SortBy, filter.Order)
	log.Debug().Msgf("Querying for persons: %s", listQuery)

	rows, err := dao.cm.DB.Query(context.Background(),
		listQuery, filter.UserID, filter.Limit, filter.Limit*filter.Page)

	if err != nil {
		log.Error().Err(err).Msg("Error getting persons list from DB")
		return nil, errs.ErrDB
	}
	personPage.Persons = readPersonListRows(rows)
	return &personPage, nil
}

func readPersonListRows(rows pgx.Rows) []model.Person {
	var persons []model.Person
	for rows.Next() {
		var person model.Person
		var s3key pgtype.Text
		if err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Birth, &person.Death, &person.Summary, &s3key, &person.Role); err != nil {
			log.Error().Err(err).Msgf("Failed to scan row")
			continue
		}
		log.Debug().Msgf("%s %s %s %s %s", person.FirstName, person.LastName, person.Birth, person.Death, person.ID)
		persons = append(persons, person)
	}
	return persons
}

func generateAndConditions(filter *model.ListPersonsFilter) string {
	var conditions strings.Builder
	if filter.BirthMax != nil {
		conditions.WriteString(" AND birth <= " + filter.BirthMax.String())
	}
	if filter.BirthMin != nil {
		conditions.WriteString(" AND birth >= " + filter.BirthMin.String())
	}
	if filter.DeathMax != nil {
		conditions.WriteString(" AND death <= " + filter.DeathMax.String())
	}
	if filter.DeathMin != nil {
		conditions.WriteString(" AND death >= " + filter.DeathMin.String())
	}
	if filter.NameMatch != nil {
		conditions.WriteString(" AND LOWER(CONCAT_WS(' ', first_name, last_name)) LIKE '%" + *filter.NameMatch + "%'")
	}
	if filter.ExcludeRoles != nil {
		conditions.WriteString(fmt.Sprintf(" AND role NOT IN (%s)", *filter.ExcludeRoles))
	}
	return conditions.String()
}
