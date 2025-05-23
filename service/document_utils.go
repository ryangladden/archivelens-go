package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// "github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
)

func (s *DocumentService) generateDocumentModel(request request.CreateDocumentRequest) *model.Document {

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for document titled \"%s\"", request.Title)
	}

	path := "documents/" + id.String()

	s3Key := s.storageManager.GenerateObjectKey(request.File.Filename, id, path)

	document := model.Document{
		Title:    request.Title,
		Location: request.Location,
		Date:     request.Date,
		Type:     request.Type,
		ID:       id,
		S3Key:    *s3Key,
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
	return authorships
}

func (s *DocumentService) generateListDocumentsFilter(request request.ListDocumentsRequest) *model.ListDocumentsFilter {
	titleMatch := ""
	if request.TitleMatch != nil {
		titleMatch = strings.ToLower(*request.TitleMatch)
	}
	filter := model.ListDocumentsFilter{
		UserID:       request.UserID,
		TitleMatch:   &titleMatch,
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

func (s *DocumentService) generateInlineDocument(documents []model.Document) []response.InlineDocument {
	var inlineDocuments []response.InlineDocument
	for _, document := range documents {
		inlineDocuments = append(inlineDocuments, response.InlineDocument{
			Title:  document.Title,
			Date:   document.Date,
			Author: s.generateInlinePerson(document.Author),
			Type:   document.Type,
		})
	}
	return inlineDocuments
}

// func (s *DocumentService) generateDocumentListResponse(db.DocumentPage)

func parseUUIDList(request *[]uuid.UUID) *string {
	if request != nil {
		var ids []string
		for _, id := range *request {
			ids = append(ids, id.String())
		}
		list := strings.Join(ids, ", ")
		return &list
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
			Name:         person.FirstName + person.LastName,
			PresignedURL: s.storageManager.GeneratePresignedURL(person.S3Key),
		}
	}
	return nil
}
