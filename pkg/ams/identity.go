package ams

// Provides information about the identity of the user or application
// this interfaces is implemented by github.com/sap/cloud-security-client-go/auth/Token

type Identity interface {
	AppTID() string
	ScimID() string
	UserUUID() string
	Groups() []string
	Email() string
}
