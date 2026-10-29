package org

import "time"

type Role string

const (
	RoleCEO      Role = "ceo"
	RoleDirector Role = "director"
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type Position string

const (
	PositionSeller   Position = "seller"
	PositionSecurity Position = "security"
	PositionHelper   Position = "helper"
)

type EmploymentStatus string

const (
	StatusActive    EmploymentStatus = "active"
	StatusFired     EmploymentStatus = "fired"
	StatusProbation EmploymentStatus = "probation"
)

type Employee struct {
	ID        int64
	Name      string
	Role      Role
	Position  Position // only with Employee Role
	ManagerID *int64
	Salary    int64 // in Rubles
	Status    EmploymentStatus
	ChangedAt time.Time
	HiredAt   time.Time
}
