package service

import (
	"fmt"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func parseExcludeRoles(request *[]string) *string {
	if request != nil {
		roleList := *request
		if len(roleList) == 0 {
			return nil
		}
		roles := make([]string, 0, len(roleList))
		for _, role := range roleList {
			if slices.Contains([]string{"editor", "owner", "viewer"}, role) {
				formattedRole := fmt.Sprintf("'%s'", role)
				roles = append(roles, formattedRole)
			}
		}
		formattedString := strings.Join(roles, ", ")
		return &formattedString
	}
	return nil
}

func parseSortBy(request *string, options []string, fallback string) string {
	if request != nil {
		log.Debug().Msgf("Parsed sort by, user requested sortby %s", *request)
		if slices.Contains(options, *request) {
			log.Debug().Msgf("Sorting query by %s", *request)
			return *request
		}
	}
	log.Debug().Msgf("Sorting query by %s", fallback)
	return fallback
}

func parseOrder(request *string) string {
	order := "ASC"
	if request != nil {
		log.Debug().Msgf("Parsed order, user requested order %s", *request)
		if *request == "descending" {
			order = "DESC"
			return order
		}
	}
	return order
}
