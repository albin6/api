package domain

type Role string

const (
	RoleHead   Role = "HEAD"
	RoleLead   Role = "LEAD"
	RoleMember Role = "MEMBER"
)

func (r Role) String() string {
	return string(r)
}
