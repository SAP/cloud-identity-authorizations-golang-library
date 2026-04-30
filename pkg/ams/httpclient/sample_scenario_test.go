package httpclient_test

import (
	// 	"context"
	// 	"fmt"
	// 	"testing"

	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/http/server"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/httpclient"
	// "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
)

type E1 struct {
	Size   int    `ams:"size"`
	Name   string `ams:"name"`
	Public bool   `ams:"public"`
	Group  string `ams:"group"`
}
type E2 struct {
	Name      string   `ams:"name"`
	Owners    []string `ams:"owners"`
	Subentity *struct {
		SubNumberField int `ams:"subNumberField"`
	} `ams:"subentity"`
}

type Schema struct {
	Entity1 *E1 `ams:"entity1"`
	Entity2 *E2 `ams:"entity2"`
}

type identity struct {
	appTID string
	scimID string
	email  string
	groups []string
}

func (i identity) AppTID() string {
	return i.appTID
}

func (i identity) ScimID() string {
	return i.scimID
}

func (i identity) UserUUID() string {
	return ""
}

func (i identity) Groups() []string {
	return i.groups
}
func (i identity) Email() string {
	return i.email
}

type crashLogger struct{}

func (l crashLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (l crashLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (l crashLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (l crashLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func TestSimpleScenario(t *testing.T) {

	aSrv := ams.NewAuthorizationManagerForFs("../test/scenarios/simple", crashLogger{})

	router := server.NewRouter(aSrv, crashLogger{})

	srv := httptest.NewServer(router.Mux())
	defer srv.Close()

	a := NewAuthorizationManager(srv.URL, srv.Client(), crashLogger{})
	<-a.WhenReady(context.Background())
	t.Run("random user on entity1", func(t *testing.T) {
		ctx := context.Background()
		authz := a.AuthorizationsForIdentity(
			ctx,
			identity{groups: []string{"g1", "g2"}})
		res, err := authz.GetResources(ctx)
		if err != nil {
			t.Fatalf("failed to get resources: %v", err)
		}
		sort.Strings(res)
		if !reflect.DeepEqual(res, []string{"r1", "r2"}) {
			t.Fatalf("expected resources to be [r1 r2], but was %v", res)
		}
		actions, err := authz.GetActions(ctx, "r1")
		if err != nil {
			t.Fatalf("failed to get actions: %v", err)
		}
		if !reflect.DeepEqual(actions, []string{"read"}) {
			t.Fatalf("expected actions to be [read], but was %v", actions)
		}
		actions, err = authz.GetActions(ctx, "r2")
		if err != nil {
			t.Fatalf("failed to get actions: %v", err)
		}
		if !reflect.DeepEqual(actions, []string{"read"}) {
			t.Fatalf("expected actions to be [read], but was %v", actions)
		}
		d, err := authz.Inquire(ctx, "write", "r1", nil)
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d.Condition())
		}
		d, err = authz.Inquire(ctx, "read", "r1", nil)
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if d.IsDenied() {
			t.Fatalf("expected access to be not denied, but was %s", d.Condition())
		}
		if d.IsGranted() {
			t.Fatalf("expected access to be not granted, but was %s", d.Condition())
		}
		// default policies should grant read when group is g1 or g2, or public is true
		d2, err := d.Inquire(
			ctx,
			Schema{
				Entity1: &E1{
					Group: "g1",
				},
			},
		)
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}
		d2, err = d.Inquire(ctx, Schema{
			Entity1: &E1{
				Group: "g2",
			},
		})
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}
		d2, err = d.Inquire(ctx, Schema{
			Entity1: &E1{
				Group: "g3",
			},
		})
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d2.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d2.Condition())
		}
		d2, err = d.Inquire(ctx, Schema{
			Entity1: &E1{
				Public: true,
			},
		})
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}

		authz.SetEnvInput(ams.DefaultEnvironmentInput{
			UserInfo: ams.UserInfo{
				Groups: []string{"g3"},
			},
		})
		d, err = authz.Inquire(ctx, "read", "r1", Schema{
			Entity1: &E1{
				Group: "g3",
			},
		})
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d.Condition())
		}
	})

	t.Run("nil identity always denied", func(t *testing.T) {
		authz := a.AuthorizationsForIdentity(context.Background(), nil)
		d, err := authz.Inquire(context.Background(), "", "", nil)
		if err != nil {
			t.Fatalf("failed to inquire: %v", err)
		}
		if !d.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d.Condition())
		}
	})

}
