package cs3users

import (
	"context"
	"fmt"
	"strings"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	idxerrs "github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
	"google.golang.org/grpc/metadata"
)

func init() {
	registry.IndexConstructorRegistry["cs3users"]["non_unique"] = NewNonUniqueIndexWithOptions
}

// NonUnique are fields for an index of type non_unique.
type NonUnique struct {
	caseInsensitive bool
	indexBy         string
	typeName        string

	tokenManager token.Manager
	userProvider user.UserAPIClient

	cs3conf *Config
}

// NewNonUniqueIndexWithOptions instantiates a new NonUniqueIndex instance.
// for the cs3 users api these non unique indexes need to be supported:
// - IndexBy "Id"
// - IndexBy "DisplayName"
func NewNonUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &NonUnique{
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		cs3conf: &Config{
			ProviderAddr: opts.ProviderAddr,
			JWTSecret:    opts.JWTSecret,
			ServiceUser:  opts.ServiceUser,
		},
	}
}

// Init initializes a non_unique index.
func (idx *NonUnique) Init() error {
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": idx.cs3conf.JWTSecret,
	})

	if err != nil {
		return err
	}

	idx.tokenManager = tokenManager

	client, err := pool.GetUserProviderServiceClient(idx.cs3conf.ProviderAddr)
	if err != nil {
		return err
	}

	idx.userProvider = client

	return nil
}

// Lookup exact lookup by value.
func (idx *NonUnique) Lookup(v string) ([]string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	switch idx.indexBy {
	case "Id":
		return []string{v}, nil // TODO fixme
	case "DisplayName":
		res, err := idx.userProvider.FindUsers(ctx, &user.FindUsersRequest{
			Filter: v,
		})
		if err != nil {
			return nil, err
		}
		if res.Status.Code != rpcv1beta1.Code_CODE_OK {
			return []string{}, fmt.Errorf(res.Status.Message)
		}

		var matches = make([]string, 0)
		for i := range res.Users {
			// enforce (case sensitive) exact match
			if idx.caseInsensitive && strings.EqualFold(res.Users[i].DisplayName, v) {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			} else if res.Users[i].DisplayName == v {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			}
		}
		return matches, nil
	}
	return []string{}, nil // todo error ... insupported index?
}

// Add a new value to the index.
func (idx *NonUnique) Add(id, v string) (string, error) {
	return "", nil
}

// Remove a value v from an index.
func (idx *NonUnique) Remove(id string, v string) error {
	return &idxerrs.NotSupportedErr{}
}

// Update index from <oldV> to <newV>.
func (idx *NonUnique) Update(id, oldV, newV string) error {
	return &idxerrs.NotSupportedErr{}
}

// Search allows for glob search on the index.
func (idx *NonUnique) Search(pattern string) ([]string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	switch idx.indexBy {
	// TODO why do we need to support the Id based search?
	case "Id", "DisplayName":
		res, err := idx.userProvider.FindUsers(ctx, &user.FindUsersRequest{
			Filter: pattern,
		})
		if err != nil {
			return nil, err
		}
		if res.Status.Code != rpcv1beta1.Code_CODE_OK {
			return []string{}, fmt.Errorf(res.Status.Message)
		}
		var matches = make([]string, 0)
		for i := range res.Users {
			matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
		}
		return matches, nil
	}
	// TODO error ... insupported index?
	// for now igrnore other properties because the FindUsers already contains matches for username, display name, mail and opaqueid
	return []string{}, nil
}

// CaseInsensitive undocumented.
// TODO not called anywhere?
func (idx *NonUnique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
// used for the map
func (idx *NonUnique) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
// only used by error
func (idx *NonUnique) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
// TODO not called anywhere?
func (idx *NonUnique) FilesDir() string {
	return ""
}

func (idx *NonUnique) getAuthenticatedContext(ctx context.Context) (context.Context, error) {
	t, err := idx.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	return ctx, nil
}

// Delete deletes the index folder from its storage.
func (idx *NonUnique) Delete() error {
	return &idxerrs.NotSupportedErr{}
}

func (idx *NonUnique) authenticate(ctx context.Context) (token string, err error) {
	return storage.AuthenticateCS3(ctx, idx.cs3conf.ServiceUser, idx.tokenManager)
}
