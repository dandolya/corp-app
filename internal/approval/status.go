package approval

type Status string

const (
	StatusPending   Status = "pending"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
)

func IsValid(status Status) bool {
	switch status {
	case StatusPending, StatusApproved, StatusRejected, StatusCancelled, StatusExpired:
		return true
	default:
		return false
	}
}

func IsTerminal(status Status) bool {
	return IsValid(status) && status != StatusPending
}

func CanTransition(from, to Status) bool {
	return from == StatusPending && IsTerminal(to)
}
