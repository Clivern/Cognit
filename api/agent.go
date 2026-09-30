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

// ListAgentsAction lists agents in a workspace.
func ListAgentsAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewAgentInstanceRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := am.ListAgents(db.Id(wid), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to list agents")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_agent_request"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"agents": result.Agents,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetAgentAction returns one agent by name.
func GetAgentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	agentName := chi.URLParam(r, "agentName")
	if lo.IsEmpty(agentName) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_name"),
		})
		return
	}

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewAgentInstanceRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	agent, err := am.GetAgent(db.Id(wid), agentName)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidAgentName):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_agent_name"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Str("agent", agentName).Msg("Failed to get agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_agent_request"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, agent)
}

// UpsertAgentAction registers or replaces an agent card.
func UpsertAgentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	agentName := chi.URLParam(r, "agentName")
	if lo.IsEmpty(agentName) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_name"),
		})
		return
	}

	var req module.UpsertAgentRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("agent", agentName).
		Msg("Upserting agent")

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewAgentInstanceRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	agent, created, err := am.UpsertAgent(db.Id(wid), agentName, &req)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrInvalidAgentName):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_agent_name"),
			})
		case errors.Is(err, module.ErrInvalidAgentCard):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_agent_card"),
			})
		case errors.Is(err, module.ErrAgentVersionRequired):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_agent_version"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Str("agent", agentName).Msg("Failed to upsert agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_agent_request"),
			})
		}
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	util.WriteJSON(w, status, agent)
}

// DeleteAgentAction removes an agent by name.
func DeleteAgentAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	agentName := chi.URLParam(r, "agentName")
	if lo.IsEmpty(agentName) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_name"),
		})
		return
	}

	log.Info().
		Str("workspaceId", wid).
		Str("agent", agentName).
		Msg("Deleting agent")

	am := module.NewAgent(
		db.NewAgentRepository(db.GetDB()),
		db.NewAgentInstanceRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	err := am.DeleteAgent(db.Id(wid), agentName)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrInvalidAgentName):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "invalid_agent_name"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Str("agent", agentName).Msg("Failed to delete agent")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_agent_request"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
