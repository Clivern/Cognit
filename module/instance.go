// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/clivern/cognit/db"
	"github.com/clivern/cognit/pkg/util"
)

var (
	ErrInstanceNotFound         = errors.New("instance not found")
	ErrInvalidInstanceId        = errors.New("invalid instance id")
	ErrFailedListInstances      = errors.New("failed list instances")
	ErrFailedGetInstance        = errors.New("failed get instance")
	ErrFailedRegisterInstance   = errors.New("failed register instance")
	ErrFailedRenewInstance      = errors.New("failed renew instance")
	ErrFailedDeregisterInstance = errors.New("failed deregister instance")
)

// DefaultLeaseTTL is the lease length in seconds when none is given.
const DefaultLeaseTTL = 30

var instanceIdPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`)

// Instance is the module for agent instance registration and leases.
type Instance struct {
	AgentRepository       db.AgentRepository
	InstanceRepository    db.AgentInstanceRepository
	HealthCheckRepository db.HealthCheckRepository
	WorkspaceRepository   db.WorkspaceRepository
}

// RegisterInstanceRequest is the body for registering or updating an instance.
type RegisterInstanceRequest struct {
	Address    string          `json:"address" validate:"required,max=255" label:"Address"`
	Port       int             `json:"port" validate:"required,min=1,max=65535" label:"Port"`
	Datacenter *string         `json:"datacenter" validate:"omitempty,max=60" label:"Datacenter"`
	Meta       json.RawMessage `json:"meta" label:"Meta"`
	LeaseTTL   int             `json:"leaseTtl" validate:"omitempty,min=5,max=86400" label:"Lease TTL"`
}

// HealthCheckResponse is one health check on an instance.
type HealthCheckResponse struct {
	CheckId      string     `json:"checkId"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	Output       *string    `json:"output,omitempty"`
	TTLExpiresAt *time.Time `json:"ttlExpiresAt,omitempty"`
	LastRunAt    *time.Time `json:"lastRunAt,omitempty"`
}

// InstanceResponse is an instance shaped for API responses.
// Status is passing while the lease is live, critical once it expires.
type InstanceResponse struct {
	Id             db.Id                  `json:"id"`
	AgentId        db.Id                  `json:"agentId"`
	InstanceId     string                 `json:"instanceId"`
	Address        string                 `json:"address"`
	Port           int                    `json:"port"`
	Datacenter     *string                `json:"datacenter,omitempty"`
	Meta           json.RawMessage        `json:"meta,omitempty"`
	Status         string                 `json:"status"`
	LeaseExpiresAt *time.Time             `json:"leaseExpiresAt,omitempty"`
	Checks         []*HealthCheckResponse `json:"checks"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

// LeaseDefinition is stored as the TTL check definition so renew knows the lease length.
type LeaseDefinition struct {
	TTL int `json:"ttl"`
}

// NewInstance creates an instance module with the given repositories.
func NewInstance(agents db.AgentRepository, instances db.AgentInstanceRepository, checks db.HealthCheckRepository, workspaces db.WorkspaceRepository) *Instance {
	return &Instance{
		AgentRepository:       agents,
		InstanceRepository:    instances,
		HealthCheckRepository: checks,
		WorkspaceRepository:   workspaces,
	}
}

// ListInstances returns the instances of an agent. With passing set, only
// instances still in the discovery pool are returned.
func (i *Instance) ListInstances(workspaceId db.Id, agentName string, passing bool) ([]*InstanceResponse, error) {
	if !agentNamePattern.MatchString(agentName) {
		return nil, ErrInvalidAgentName
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := i.AgentRepository.GetByWorkspaceAndName(workspaceId, agentName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInstances, err)
	}
	if agent == nil {
		return nil, ErrAgentNotFound
	}

	now := time.Now().UTC()

	var instances []*db.AgentInstance
	if passing {
		instances, err = i.InstanceRepository.ListLiveByAgentId(agent.Id, now)
	} else {
		instances, err = i.InstanceRepository.ListByAgentId(agent.Id)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInstances, err)
	}

	list := make([]*InstanceResponse, 0, len(instances))
	for _, instance := range instances {
		checks, err := i.HealthCheckRepository.ListByAgentInstanceId(instance.Id)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedListInstances, err)
		}

		item := &InstanceResponse{
			Id:         instance.Id,
			AgentId:    instance.AgentId,
			InstanceId: instance.InstanceId,
			Address:    instance.Address,
			Port:       instance.Port,
			Datacenter: instance.Datacenter,
			Meta:       util.JSONRawFromString(instance.Meta),
			Status:     db.AgentInstanceStatusCritical,
			Checks:     make([]*HealthCheckResponse, 0, len(checks)),
			CreatedAt:  instance.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:  instance.UpdatedAt.UTC().Format(time.RFC3339),
		}
		for _, check := range checks {
			if check.Type == db.HealthCheckTypeTTL {
				item.LeaseExpiresAt = check.TTLExpiresAt
				if check.TTLExpiresAt.After(now) {
					item.Status = db.AgentInstanceStatusPassing
				}
			}
			item.Checks = append(item.Checks, &HealthCheckResponse{
				CheckId:      check.CheckId,
				Name:         check.Name,
				Type:         check.Type,
				Status:       check.Status,
				Output:       check.Output,
				TTLExpiresAt: check.TTLExpiresAt,
				LastRunAt:    check.LastRunAt,
			})
		}

		list = append(list, item)
	}
	return list, nil
}

// GetInstance returns one instance of an agent.
func (i *Instance) GetInstance(workspaceId db.Id, agentName, instanceId string) (*InstanceResponse, error) {
	if !agentNamePattern.MatchString(agentName) {
		return nil, ErrInvalidAgentName
	}
	if !instanceIdPattern.MatchString(instanceId) {
		return nil, ErrInvalidInstanceId
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := i.AgentRepository.GetByWorkspaceAndName(workspaceId, agentName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInstance, err)
	}
	if agent == nil {
		return nil, ErrAgentNotFound
	}

	instance, err := i.InstanceRepository.GetByAgentAndInstanceId(agent.Id, instanceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInstance, err)
	}
	if instance == nil {
		return nil, ErrInstanceNotFound
	}

	checks, err := i.HealthCheckRepository.ListByAgentInstanceId(instance.Id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInstance, err)
	}

	now := time.Now().UTC()
	item := &InstanceResponse{
		Id:         instance.Id,
		AgentId:    instance.AgentId,
		InstanceId: instance.InstanceId,
		Address:    instance.Address,
		Port:       instance.Port,
		Datacenter: instance.Datacenter,
		Meta:       util.JSONRawFromString(instance.Meta),
		Status:     db.AgentInstanceStatusCritical,
		Checks:     make([]*HealthCheckResponse, 0, len(checks)),
		CreatedAt:  instance.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  instance.UpdatedAt.UTC().Format(time.RFC3339),
	}
	for _, check := range checks {
		if check.Type == db.HealthCheckTypeTTL {
			item.LeaseExpiresAt = check.TTLExpiresAt
			if check.TTLExpiresAt.After(now) {
				item.Status = db.AgentInstanceStatusPassing
			}
		}
		item.Checks = append(item.Checks, &HealthCheckResponse{
			CheckId:      check.CheckId,
			Name:         check.Name,
			Type:         check.Type,
			Status:       check.Status,
			Output:       check.Output,
			TTLExpiresAt: check.TTLExpiresAt,
			LastRunAt:    check.LastRunAt,
		})
	}
	return item, nil
}

// RegisterInstance creates or updates an instance and starts its lease.
func (i *Instance) RegisterInstance(workspaceId db.Id, agentName, instanceId string, req *RegisterInstanceRequest) (*InstanceResponse, bool, error) {
	if !agentNamePattern.MatchString(agentName) {
		return nil, false, ErrInvalidAgentName
	}
	if !instanceIdPattern.MatchString(instanceId) {
		return nil, false, ErrInvalidInstanceId
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, false, err
	}
	if workspace == nil {
		return nil, false, ErrWorkspaceNotFound
	}

	agent, err := i.AgentRepository.GetByWorkspaceAndName(workspaceId, agentName)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}
	if agent == nil {
		return nil, false, ErrAgentNotFound
	}

	instance, err := i.InstanceRepository.GetByAgentAndInstanceId(agent.Id, instanceId)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}

	created := instance == nil
	if created {
		instance = &db.AgentInstance{AgentId: agent.Id, InstanceId: instanceId}
	}
	instance.Address = req.Address
	instance.Port = req.Port
	instance.Datacenter = req.Datacenter
	instance.Meta = nil
	if len(req.Meta) > 0 {
		meta := string(req.Meta)
		instance.Meta = &meta
	}
	instance.Status = db.AgentInstanceStatusPassing

	if created {
		err = i.InstanceRepository.Create(instance)
	} else {
		err = i.InstanceRepository.Update(instance)
	}
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}

	ttl := req.LeaseTTL
	if ttl == 0 {
		ttl = DefaultLeaseTTL
	}
	raw, _ := json.Marshal(LeaseDefinition{TTL: ttl})
	definition := string(raw)
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(ttl) * time.Second)

	lease, err := i.HealthCheckRepository.GetByInstanceAndCheckId(instance.Id, db.HealthCheckIDTTL)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}
	if lease == nil {
		lease = &db.HealthCheck{
			AgentInstanceId: instance.Id,
			CheckId:         db.HealthCheckIDTTL,
			Name:            "Lease",
			Type:            db.HealthCheckTypeTTL,
			Status:          db.HealthCheckStatusPassing,
			Definition:      &definition,
			TTLExpiresAt:    &expiresAt,
			LastRunAt:       &now,
		}
		err = i.HealthCheckRepository.Create(lease)
	} else {
		lease.Definition = &definition
		err = i.HealthCheckRepository.Update(lease)
		if err == nil {
			err = i.HealthCheckRepository.Pass(lease.Id, "", &expiresAt)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}

	checks, err := i.HealthCheckRepository.ListByAgentInstanceId(instance.Id)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrFailedRegisterInstance, err)
	}

	item := &InstanceResponse{
		Id:         instance.Id,
		AgentId:    instance.AgentId,
		InstanceId: instance.InstanceId,
		Address:    instance.Address,
		Port:       instance.Port,
		Datacenter: instance.Datacenter,
		Meta:       util.JSONRawFromString(instance.Meta),
		Status:     db.AgentInstanceStatusPassing,
		Checks:     make([]*HealthCheckResponse, 0, len(checks)),
		CreatedAt:  instance.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  now.Format(time.RFC3339),
	}
	for _, check := range checks {
		if check.Type == db.HealthCheckTypeTTL {
			item.LeaseExpiresAt = check.TTLExpiresAt
		}
		item.Checks = append(item.Checks, &HealthCheckResponse{
			CheckId:      check.CheckId,
			Name:         check.Name,
			Type:         check.Type,
			Status:       check.Status,
			Output:       check.Output,
			TTLExpiresAt: check.TTLExpiresAt,
			LastRunAt:    check.LastRunAt,
		})
	}
	return item, created, nil
}

// RenewInstance extends an instance lease by its registered TTL.
func (i *Instance) RenewInstance(workspaceId db.Id, agentName, instanceId string) (*InstanceResponse, error) {
	if !agentNamePattern.MatchString(agentName) {
		return nil, ErrInvalidAgentName
	}
	if !instanceIdPattern.MatchString(instanceId) {
		return nil, ErrInvalidInstanceId
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	agent, err := i.AgentRepository.GetByWorkspaceAndName(workspaceId, agentName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRenewInstance, err)
	}
	if agent == nil {
		return nil, ErrAgentNotFound
	}

	instance, err := i.InstanceRepository.GetByAgentAndInstanceId(agent.Id, instanceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRenewInstance, err)
	}
	if instance == nil {
		return nil, ErrInstanceNotFound
	}

	// Registering always creates the lease, so it is there to extend.
	lease, err := i.HealthCheckRepository.GetByInstanceAndCheckId(instance.Id, db.HealthCheckIDTTL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRenewInstance, err)
	}
	if lease == nil {
		return nil, ErrInstanceNotFound
	}

	var def LeaseDefinition
	json.Unmarshal([]byte(*lease.Definition), &def)
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(def.TTL) * time.Second)

	err = i.HealthCheckRepository.Pass(lease.Id, "", &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRenewInstance, err)
	}

	checks, err := i.HealthCheckRepository.ListByAgentInstanceId(instance.Id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRenewInstance, err)
	}

	item := &InstanceResponse{
		Id:         instance.Id,
		AgentId:    instance.AgentId,
		InstanceId: instance.InstanceId,
		Address:    instance.Address,
		Port:       instance.Port,
		Datacenter: instance.Datacenter,
		Meta:       util.JSONRawFromString(instance.Meta),
		Status:     db.AgentInstanceStatusPassing,
		Checks:     make([]*HealthCheckResponse, 0, len(checks)),
		CreatedAt:  instance.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  instance.UpdatedAt.UTC().Format(time.RFC3339),
	}
	for _, check := range checks {
		if check.Type == db.HealthCheckTypeTTL {
			item.LeaseExpiresAt = check.TTLExpiresAt
		}
		item.Checks = append(item.Checks, &HealthCheckResponse{
			CheckId:      check.CheckId,
			Name:         check.Name,
			Type:         check.Type,
			Status:       check.Status,
			Output:       check.Output,
			TTLExpiresAt: check.TTLExpiresAt,
			LastRunAt:    check.LastRunAt,
		})
	}
	return item, nil
}

// DeregisterInstance removes an instance and its health checks.
func (i *Instance) DeregisterInstance(workspaceId db.Id, agentName, instanceId string) error {
	if !agentNamePattern.MatchString(agentName) {
		return ErrInvalidAgentName
	}
	if !instanceIdPattern.MatchString(instanceId) {
		return ErrInvalidInstanceId
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	agent, err := i.AgentRepository.GetByWorkspaceAndName(workspaceId, agentName)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeregisterInstance, err)
	}
	if agent == nil {
		return ErrAgentNotFound
	}

	instance, err := i.InstanceRepository.GetByAgentAndInstanceId(agent.Id, instanceId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeregisterInstance, err)
	}
	if instance == nil {
		return ErrInstanceNotFound
	}

	err = i.InstanceRepository.Delete(instance.Id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeregisterInstance, err)
	}
	return nil
}
