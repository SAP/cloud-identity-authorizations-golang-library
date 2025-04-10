package dcn

type (
	TenantID string
	UserID   string
)

type UserAssignments map[string][]string

type (
	Assignments          map[string]UserAssignments
	AssignmentsContainer struct {
		Assignments Assignments `json:"principal2policies"`
	}
)
