package httpclient_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/sap/cloud-identity-authorizations-golang-library/http/server"
// 	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
// )

// type crashLogger struct{}

// func (l crashLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
// func (l crashLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
// func (l crashLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
// func (l crashLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
// 	panic(fmt.Sprintf(format, args...))
// }

// func TestXxx(t *testing.T) {

// 	a := ams.NewAuthorizationManagerForFs("test/scenarios/simple", crashLogger{})

// 	router := server.NewRouter(a, crashLogger{})

// }
