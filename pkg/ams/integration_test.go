package ams

import (
	"context"
	"reflect"
	"sort"
	"testing"

	test "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/test/bundles"
)

func TestAuthorizationManagerWithMockBundleGateway(t *testing.T) {
	mock := test.NewBundleGatewayMock()
	defer mock.Close()

	a, err := newAuthorizationManagerForIAS(
		mock.GetAuthorizationBundleURL(),
		mock.GetAuthorizationInstanceID(),
		mock.GetHttpClient(),
		func(err error) {
			t.Fatalf("Error callback called: %v", err)
		},
	)
	a.Run(context.Background())

	if err != nil {
		t.Fatalf("Failed to create AuthorizationManager: %v", err)
	}

	t.Run("random user on entity1", func(t *testing.T) {
		authz := a.AuthorizationsForIdentity(identity{groups: []string{"g1", "g2"}})
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

}
