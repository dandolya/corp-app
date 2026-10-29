package approval

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"pending", StatusPending, true},
		{"expired", StatusExpired, true},
		{"typo", "approvd", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.status); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{"approve", StatusPending, StatusApproved, true},
		{"reject", StatusPending, StatusRejected, true},
		{"cancel", StatusPending, StatusCancelled, true},
		{"expire", StatusPending, StatusExpired, true},
		{"reverse", StatusApproved, StatusPending, false},
		{"resolve twice", StatusApproved, StatusRejected, false},
		{"same status", StatusPending, StatusPending, false},
		{"typo", "pending", "approvd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Errorf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
