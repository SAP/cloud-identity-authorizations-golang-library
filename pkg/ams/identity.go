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

type UserInfo struct {
	Email    string   `ams:"email"`
	Groups   []string `ams:"groups"`
	UserUUID string   `ams:"user_uuid"`
}

type DefaultEnvironmentInput struct {
	UserInfo UserInfo `ams:"$user"`
}
type IASConfig interface {
	GetAuthorizationBundleURL() string
	GetAuthorizationInstanceID() string
	GetCertifcate() string
	GetKey() string
}
