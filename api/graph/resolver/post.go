package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
)

// Post represents the resolver of the identically named type
type Post struct {
	root     *Resolver
	uid      store.UID
	id       store.ID
	creation time.Time
	title    string
	contents string
}

// Id resolves Post.id
func (rsv *Post) Id() store.ID {
	return rsv.id
}

// Author resolves Post.author
func (rsv *Post) Author(ctx context.Context) (*User, error) {
	var query struct {
		Post []struct {
			Author []dbmod.User `json:"Post.author"`
		} `json:"post"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Author($nodeId: string) {
			post(func: uid($nodeId)) {
				Post.author {
					uid
					User.id
					User.creation
					User.email
					User.displayName
					User.posts
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid.NodeID,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	author := query.Post[0].Author[0]
	return &User{
		root:        rsv.root,
		uid:         store.UID{NodeID: author.UID},
		id:          author.ID,
		creation:    author.Creation,
		email:       author.Email,
		displayName: author.DisplayName,
	}, nil
}

// Creation resolves Post.creation
func (rsv *Post) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Title resolves Post.title
func (rsv *Post) Title() string {
	return rsv.title
}

// Contents resolves Post.contents
func (rsv *Post) Contents() string {
	return rsv.contents
}

// Reactions resolves Post.reactions
func (rsv *Post) Reactions() ([]*Reaction, error) {
	return nil, nil
}