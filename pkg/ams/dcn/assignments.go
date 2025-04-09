package dcn

type TenantID string
type UserID string

type UserAssignments map[string][]string

type Assignments map[string]UserAssignments
type AssignmentsContainer struct {
	Assignments Assignments `json:"principal2policies"`
}
