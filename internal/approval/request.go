package approval

import (
	"encoding/json"
	"time"
)

type Type string

const (
	TypeSalaryChange   Type = "salary_change"
	TypePenalty        Type = "penalty"
	TypeResignation    Type = "resignation"
	TypeTermination    Type = "termination"
	TypeVacation       Type = "vacation"
	TypeSickLeave      Type = "sick_leave"
	TypeDayOff         Type = "day_off"
	TypePositionChange Type = "position_change"
	TypePromotion      Type = "promotion"
	TypeTransfer       Type = "transfer"
	TypeHire           Type = "hire"
)

type Request struct {
	ID                int64
	Type              Type
	Status            Status
	InitiatorID       int64
	TargetID          *int64
	Details           json.RawMessage
	CreatedAt         time.Time
	ExpiresAt         time.Time
	ResolvedBy        *int64
	ResolvedAt        *time.Time
	ResolutionComment string
}
