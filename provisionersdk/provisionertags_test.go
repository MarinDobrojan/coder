package provisionersdk_test

import (
	"testing"

	"github.com/coder/coder/v2/provisionersdk"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMutateTags(t *testing.T) {
	t.Parallel()

	testUserID := uuid.New()

	for _, tt := range []struct {
		name   string
		userID uuid.UUID
		tags   []map[string]string
		want   map[string]string
	}{
		{
			name:   "nil tags",
			userID: uuid.Nil,
			tags:   mapslice(nil),
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name:   "empty tags",
			userID: uuid.Nil,
			tags:   mapslice(map[string]string{}),
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "user scope",
			tags: mapslice(
				map[string]string{provisionersdk.TagScope: provisionersdk.ScopeUser},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeUser,
				provisionersdk.TagOwner: testUserID.String(),
			},
		},
		{
			name: "organization scope",
			tags: mapslice(
				map[string]string{provisionersdk.TagScope: provisionersdk.ScopeOrganization},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "organization scope with owner",
			tags: mapslice(
				map[string]string{
					provisionersdk.TagScope: provisionersdk.ScopeOrganization,
					provisionersdk.TagOwner: testUserID.String(),
				},
			),
			userID: uuid.Nil,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "owner tag with no other context",
			tags: mapslice(
				map[string]string{
					provisionersdk.TagOwner: testUserID.String(),
				},
			),
			userID: uuid.Nil,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "invalid scope",
			tags: mapslice(
				map[string]string{provisionersdk.TagScope: "360noscope"},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "merge two empty maps",
			tags: mapslice(
				map[string]string{},
				map[string]string{},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
			},
		},
		{
			name: "merge empty map with non-empty map",
			tags: mapslice(
				map[string]string{},
				map[string]string{"foo": "bar"},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
				"foo":                   "bar",
			},
		},
		{
			name: "merge non-empty map with empty map",
			tags: mapslice(
				map[string]string{"foo": "bar"},
				map[string]string{},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
				"foo":                   "bar",
			},
		},
		{
			name: "merge map with same map",
			tags: mapslice(
				map[string]string{"foo": "bar"},
				map[string]string{"foo": "bar"},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
				"foo":                   "bar",
			},
		},
		{
			name: "merge map with override",
			tags: mapslice(
				map[string]string{"foo": "bar"},
				map[string]string{"foo": "baz"},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
				"foo":                   "baz",
			},
		},
		{
			name: "do not override empty in second map",
			tags: mapslice(
				map[string]string{"foo": "bar"},
				map[string]string{"foo": ""},
			),
			userID: testUserID,
			want: map[string]string{
				provisionersdk.TagScope: provisionersdk.ScopeOrganization,
				provisionersdk.TagOwner: "",
				"foo":                   "bar",
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := provisionersdk.MutateTags(tt.userID, tt.tags...)
			require.Equal(t, tt.want, got)
		})
	}
}

func mapslice(m ...map[string]string) []map[string]string {
	return m
}
