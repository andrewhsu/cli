package shared

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/andrewhsu/cli/v2/api"
	"github.com/andrewhsu/cli/v2/internal/ghrepo"
	"github.com/andrewhsu/cli/v2/pkg/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_listURLWithQuery(t *testing.T) {
	type args struct {
		listURL string
		options FilterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "blank",
			args: args{
				listURL: "https://example.com/path?a=b",
				options: FilterOptions{
					Entity: "issue",
					State:  "open",
				},
			},
			want:    "https://example.com/path?a=b&q=is%3Aissue+is%3Aopen",
			wantErr: false,
		},
		{
			name: "all",
			args: args{
				listURL: "https://example.com/path",
				options: FilterOptions{
					Entity:     "issue",
					State:      "open",
					Assignee:   "bo",
					Author:     "ka",
					BaseBranch: "trunk",
					Mention:    "nu",
				},
			},
			want:    "https://example.com/path?q=is%3Aissue+is%3Aopen+assignee%3Abo+author%3Aka+mentions%3Anu+base%3Atrunk",
			wantErr: false,
		},
		{
			name: "spaces in values",
			args: args{
				listURL: "https://example.com/path",
				options: FilterOptions{
					Entity:    "pr",
					State:     "open",
					Labels:    []string{"docs", "help wanted"},
					Milestone: `Codename "What Was Missing"`,
				},
			},
			want:    "https://example.com/path?q=is%3Apr+is%3Aopen+label%3Adocs+label%3A%22help+wanted%22+milestone%3A%22Codename+%5C%22What+Was+Missing%5C%22%22",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListURLWithQuery(tt.args.listURL, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("listURLWithQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("listURLWithQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeReplacer_Replace(t *testing.T) {
	rtSuccess := &httpmock.Registry{}
	rtSuccess.Register(
		httpmock.GraphQL(`query UserCurrent\b`),
		httpmock.StringResponse(`
		{ "data": {
			"viewer": { "login": "ResolvedLogin" }
		} }
		`),
	)

	rtFailure := &httpmock.Registry{}
	rtFailure.Register(
		httpmock.GraphQL(`query UserCurrent\b`),
		httpmock.StatusStringResponse(500, `
		{ "data": {
			"viewer": { }
		} }
		`),
	)

	type args struct {
		logins []string
		client *api.Client
		repo   ghrepo.Interface
	}
	tests := []struct {
		name    string
		args    args
		verify  func(t httpmock.Testing)
		want    []string
		wantErr bool
	}{
		{
			name: "succeeds resolving the userlogin",
			args: args{
				client: api.NewClientFromHTTP(&http.Client{Transport: rtSuccess}),
				repo:   ghrepo.New("OWNER", "REPO"),
				logins: []string{"some", "@me", "other"},
			},
			verify: rtSuccess.Verify,
			want:   []string{"some", "ResolvedLogin", "other"},
		},
		{
			name: "fails resolving the userlogin",
			args: args{
				client: api.NewClientFromHTTP(&http.Client{Transport: rtFailure}),
				repo:   ghrepo.New("OWNER", "REPO"),
				logins: []string{"some", "@me", "other"},
			},
			verify:  rtFailure.Verify,
			want:    []string(nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := NewMeReplacer(tt.args.client, tt.args.repo.RepoHost())
			got, err := me.ReplaceSlice(tt.args.logins)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceAtMeLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceAtMeLogin() = %v, want %v", got, tt.want)
			}

			if tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}

func Test_QueryHasStateClause(t *testing.T) {
	tests := []struct {
		searchQuery string
		hasState    bool
	}{
		{
			searchQuery: "is:closed is:merged",
			hasState:    true,
		},
		{
			searchQuery: "author:mislav",
			hasState:    false,
		},
		{
			searchQuery: "assignee:g14a mentions:vilmibm",
			hasState:    false,
		},
		{
			searchQuery: "merged:>2021-05-20",
			hasState:    true,
		},
		{
			searchQuery: "state:merged state:open",
			hasState:    true,
		},
		{
			searchQuery: "assignee:g14a is:closed",
			hasState:    true,
		},
		{
			searchQuery: "state:closed label:bug",
			hasState:    true,
		},
	}
	for _, tt := range tests {
		gotState := QueryHasStateClause(tt.searchQuery)
		assert.Equal(t, tt.hasState, gotState)
	}
}

func Test_WithPrAndIssueQueryParams(t *testing.T) {
	type args struct {
		baseURL string
		state   IssueMetadataState
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "blank",
			args: args{
				baseURL: "",
				state:   IssueMetadataState{},
			},
			want: "?body=",
		},
		{
			name: "no values",
			args: args{
				baseURL: "http://example.com/hey",
				state:   IssueMetadataState{},
			},
			want: "http://example.com/hey?body=",
		},
		{
			name: "title and body",
			args: args{
				baseURL: "http://example.com/hey",
				state: IssueMetadataState{
					Title: "my title",
					Body:  "my bodeh",
				},
			},
			want: "http://example.com/hey?body=my+bodeh&title=my+title",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WithPrAndIssueQueryParams(nil, nil, tt.args.baseURL, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithPrAndIssueQueryParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WithPrAndIssueQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
