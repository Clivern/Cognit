// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/cognit/pkg/util"

	"github.com/rs/zerolog/log"
)

// HealthAction returns a simple health check (status ok).
func (a *API) HealthAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Health check")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
