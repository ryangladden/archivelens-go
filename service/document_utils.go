package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
)

func (s *DocumentService) generateDocumentModel(request request.CreateDocumentRequest) *model.Document {

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for document titled \"%s\"", request.Title)
	}
	// s3Key := fmt.Sprintf("documents/%s/document%s", id.String(), filepath.Ext(request.File.Filename))
	// s3Key := storage.GenerateObjectKey("documents", id, request.File.Filename, "original")

	// s3Key := s.storageManager.GenerateObjectKey(request.File.Filename, id, path)
	original := request.File.Filename
	document := model.Document{
		Title:            request.Title,
		Location:         request.Location,
		Date:             request.Date,
		Type:             request.Type,
		ID:               id,
		OriginalFilename: original,
	}
	return &document
}

func createAuthorship(personIds []string, documentId string, role string) []model.Authorship {
	var authorships []model.Authorship
	for _, id := range personIds {
		authorships = append(authorships, model.Authorship{
			PersonID:   id,
			DocumentID: documentId,
			Role:       role,
		})
	}
	return authorships
}

func generateAuthorshipArray(documentId string, request request.CreateDocumentRequest) []model.Authorship {
	var authorships []model.Authorship
	if request.Author != nil {
		authorships = append(authorships, createAuthorship([]string{*request.Author}, documentId, "author")...)
	}
	if request.Coauthors != nil {
		authorships = append(authorships, createAuthorship(*request.Coauthors, documentId, "coauthor")...)
	}
	if request.Mentions != nil {
		authorships = append(authorships, createAuthorship(*request.Mentions, documentId, "subject")...)
	}
	if request.Recipient != nil {
		authorships = append(authorships, createAuthorship([]string{*request.Recipient}, documentId, "recipient")...)
	}
	log.Debug().Interface("authors", authorships)
	return authorships
}

func (s *DocumentService) generateListDocumentsFilter(request request.ListDocumentsRequest) *model.ListDocumentsFilter {
	filter := model.ListDocumentsFilter{
		UserID:       request.UserID,
		TitleMatch:   request.TitleMatch,
		DateMin:      request.DateMin,
		DateMax:      request.DateMax,
		ExcludeRoles: parseExcludeRoles(request.ExcludeRoles),
		SortBy:       parseSortBy(request.SortBy, []string{"title", "date", "last_name"}, "title"),
		Order:        parseOrder(request.Order),
		Authors:      parseUUIDList(request.Authors),
		IncludeTags:  parseTags(request.IncludeTags),
	}
	if request.Limit == nil {
		filter.Limit = 20
	} else {
		filter.Limit = *request.Limit
	}
	if request.Page == nil {
		filter.Page = 0
	} else {
		filter.Page = *request.Page - 1
	}
	return &filter
}

// func (s *DocumentService) generateInlineDocument(documents []model.Document) []response.InlineDocument {
// 	var inlineDocuments []response.InlineDocument
// 	for _, document := range documents {
// 		inlineDocuments = append(inlineDocuments, response.InlineDocument{
// 			Title:  document.Title,
// 			Date:   document.Date,
// 			Author: s.generateInlinePerson(document.Author),
// 			Type:   document.Type,
// 		})
// 	}
// 	return inlineDocuments
// }

// func (s *DocumentService) generateDocumentListResponse(db.DocumentPage)

func parseUUIDList(request *[]string) *string {
	if request != nil {
		var ids []string
		for _, id := range *request {
			if err := uuid.Validate(id); err == nil {
				ids = append(ids, fmt.Sprintf("'%s'", id))
			}
		}
		if len(ids) != 0 {
			list := strings.Join(ids, ", ")
			return &list
		}
	}
	return nil
}

func parseTags(request *[]string) *string {
	if request != nil {
		tags := strings.Join(*request, ", ")
		return &tags
	}
	return nil
}

func (s *DocumentService) generateInlinePerson(person *model.Person) *response.InlinePerson {
	if person != nil {
		return &response.InlinePerson{
			ID:           person.ID,
			FirstName:    person.FirstName,
			LastName:     person.LastName,
			PresignedURL: s.storageManager.GeneratePresignedURL(person.S3Key),
		}
	}
	return nil
}

func (s *DocumentService) generateInlinePersonList(persons *[]model.Person) *[]response.InlinePerson {
	if persons != nil {
		var inlinePersons []response.InlinePerson
		for _, person := range *persons {
			inlinePersons = append(inlinePersons, *s.generateInlinePerson(&person))
		}
		return &inlinePersons
	}
	return nil
}

func (s *DocumentService) generateListDocumentsResponse(page *db.DocumentPage) *response.ListDocumentsResponse {
	var listResponse response.ListDocumentsResponse
	for _, document := range page.Documents {
		s3key := fmt.Sprintf("documents/%s/thumb.webp", document.Document.ID)
		thumb := s.storageManager.GeneratePresignedURL(&s3key)
		log.Debug().Msg(*thumb)
		inlineDocument := response.InlineDocument{
			ID:    document.Document.ID,
			Title: document.Document.Title,
			Date:  document.Document.Date,
			Type:  document.Document.Type,
			Author: &response.InlinePerson{
				ID:        document.DocumentMetadata.Author.ID,
				FirstName: document.DocumentMetadata.Author.FirstName,
				LastName:  document.DocumentMetadata.Author.LastName,
			},
			Role:      document.Document.Role,
			Thumbnail: *thumb,
		}
		inlineDocument.Persons, inlineDocument.Tags = s.parseSearchMetadata(document)
		listResponse.Documents = append(listResponse.Documents, inlineDocument)
	}
	return &listResponse
}

func (s *DocumentService) parseSearchMetadata(document db.InlineDocument) (*[]response.InlinePerson, *[]response.Tag) {
	var persons []response.InlinePerson
	for _, personData := range document.DocumentMetadata.Persons {
		person := response.InlinePerson{
			ID:        personData.ID,
			FirstName: personData.FirstName,
			LastName:  personData.LastName,
			Role:      personData.Role,
		}
		persons = append(persons, person)
	}
	var tags []response.Tag
	for _, tagData := range document.DocumentMetadata.Tags {
		tag := response.Tag{
			ID:  tagData.ID,
			Tag: tagData.Tag,
		}
		tags = append(tags, tag)
	}
	return &persons, &tags
}

// func (s *DocumentService) getS3Key(uuid uuid.UUID, name string)

// func (s *DocumentService) convertFile(file *multipart.FileHeader) {
// }
