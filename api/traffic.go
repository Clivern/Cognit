// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/cognit/db"
	"github.com/clivern/cognit/locale"
	"github.com/clivern/cognit/module"
	"github.com/clivern/cognit/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// ListTrafficAction lists recorded gateway calls for a workspace.
func ListTrafficAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	tm := module.NewTraffic(
		db.NewTrafficRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := tm.ListTraffic(db.Id(wid), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to list traffic")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_traffic_request"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"calls": result.Calls,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetTrafficAction returns one recorded gateway call.
func GetTrafficAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	trafficId := chi.URLParam(r, "trafficId")
	if lo.IsEmpty(trafficId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_traffic_id"),
		})
		return
	}

	tm := module.NewTraffic(
		db.NewTrafficRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	call, err := tm.GetTraffic(db.Id(wid), db.Id(trafficId))
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrTrafficNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "traffic_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Str("trafficId", trafficId).Msg("Failed to get traffic")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_traffic_request"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, call)
}
