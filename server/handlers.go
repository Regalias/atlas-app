package server

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) handleServeLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkPath := httprouter.ParamsFromContext(r.Context()).ByName("linkpath")

		if !checkValidPath(linkPath) {
			// Immediately throw a 404
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, http.StatusText(http.StatusNotFound))
			return
		}

		// Lookup from cache
		dest, err := s.cacheProvider.FetchLink(linkPath)
		if err != nil {
			if err.Error() == "NotFound" {
				// Query from DB
				s.logger.Debug().Msg("Link not found from cache, querying DB")
				lm, err := s.databaseProvider.GetLinkDetails(linkPath)
				if err != nil {
					if err.Error() == "NotFound" {
						// We don't have this link
						// TODO: cache 404?
						w.WriteHeader(http.StatusNotFound)
						io.WriteString(w, http.StatusText(http.StatusNotFound))
						return
					}
					s.handleOperationalError(w, err)
					return
				}

				// Stash into cache
				// TODO: some actually decent caching logic...
				if err := s.cacheProvider.UpsertLink(linkPath, lm.TargetURL); err != nil {
					// Log but don't break request as we already have response
					s.logger.Error().Str("Error", "OperationalError").Msg(err.Error())
				}

				dest = lm.TargetURL

			} else {
				// Something bad happened
				s.handleOperationalError(w, err)
				return
			}
		}

		http.Redirect(w, r, dest, http.StatusFound)
		return
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, s.indexSite, http.StatusFound)
	}
}

// Use a manual check here instead of regex to salvage a bit of performance
func checkValidPath(s string) bool {
	if len(s) < minPathLength || len(s) > maxPathLength {
		return false
	}
	for _, char := range s {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || (char == '-')) {
			return false
		}
	}
	return true
}

func (s *server) handleOperationalError(w http.ResponseWriter, err error) {
	w.WriteHeader(operationalErrorCode)
	io.WriteString(w, http.StatusText(operationalErrorCode))
	s.logger.Error().Str("Error", "OperationalError").Msg(err.Error())
}
