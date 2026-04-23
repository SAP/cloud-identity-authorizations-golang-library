package ams

import (
	"context"
	"reflect"
	"sort"
	"testing"
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

func (l crashLogger) Debug(ctx context.Context, msg string) {

}
func (l crashLogger) Info(ctx context.Context, msg string) {

}
func (l crashLogger) Warn(ctx context.Context, msg string) {

}
func (l crashLogger) Error(ctx context.Context, msg string) {
	panic(msg)
}

func TestSimpleScenario(t *testing.T) {
	a := NewAuthorizationManagerForFs("test/scenarios/simple", crashLogger{})

	<-a.WhenReady()
	t.Run("random user on entity1", func(t *testing.T) {
		authz := a.AuthorizationsForIdentity(
			context.Background(),
			identity{groups: []string{"g1", "g2"}})
		res := authz.GetResources()
		sort.Strings(res)
		if !reflect.DeepEqual(res, []string{"r1", "r2"}) {
			t.Fatalf("expected resources to be [r1 r2], but was %v", res)
		}
		actions := authz.GetActions("r1")
		if !reflect.DeepEqual(actions, []string{"read"}) {
			t.Fatalf("expected actions to be [read], but was %v", actions)
		}
		actions = authz.GetActions("r2")
		if !reflect.DeepEqual(actions, []string{"read"}) {
			t.Fatalf("expected actions to be [read], but was %v", actions)
		}
		d := authz.Inquire("write", "r1", nil)
		if !d.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d.Condition())
		}
		d = authz.Inquire("read", "r1", nil)
		if d.IsDenied() {
			t.Fatalf("expected access to be not denied, but was %s", d.Condition())
		}
		if d.IsGranted() {
			t.Fatalf("expected access to be not granted, but was %s", d.Condition())
		}
		// default policies should grant read when group is g1 or g2, or public is true
		d2 := d.Inquire(Schema{
			Entity1: &E1{
				Group: "g1",
			},
		})
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}
		d2 = d.Inquire(Schema{
			Entity1: &E1{
				Group: "g2",
			},
		})
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}
		d2 = d.Inquire(Schema{
			Entity1: &E1{
				Group: "g3",
			},
		})
		if !d2.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d2.Condition())
		}
		d2 = d.Inquire(Schema{
			Entity1: &E1{
				Public: true,
			},
		})
		if !d2.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d2.Condition())
		}

		authz.SetEnvInput(DefaultEnvironmentInput{
			UserInfo: UserInfo{
				Groups: []string{"g3"},
			},
		})
		d = authz.Inquire("read", "r1", Schema{
			Entity1: &E1{
				Group: "g3",
			},
		})
		if !d.IsGranted() {
			t.Fatalf("expected access to be granted, but was %s", d.Condition())
		}
	})

	t.Run("nil identity always denied", func(t *testing.T) {
		authz := a.AuthorizationsForIdentity(context.Background(), nil)
		d := authz.Inquire("", "", nil)
		if !d.IsDenied() {
			t.Fatalf("expected access to be denied, but was %s", d.Condition())
		}
	})
}
