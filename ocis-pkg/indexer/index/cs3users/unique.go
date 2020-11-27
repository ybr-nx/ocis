package cs3users

import (
	"context"
	"fmt"
	"strings"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	acccfg "github.com/owncloud/ocis/accounts/pkg/config"

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

// Unique are fields for an index of type non_unique.
type Unique struct {
	caseInsensitive bool
	indexBy         string
	typeName        string

	tokenManager token.Manager
	userProvider user.UserAPIClient

	cs3conf *Config
}

// Config represents cs3conf. Should be deprecated in favor of config.Config.
type Config struct {
	ProviderAddr string
	JWTSecret    string
	ServiceUser  acccfg.ServiceUser
}

func init() {
	registry.IndexConstructorRegistry["cs3users"]["unique"] = NewUniqueIndexWithOptions
}

// NewUniqueIndexWithOptions instantiates a new UniqueIndex instance.
// for the cs3 users api these non unique indexes need to be supported:
// - IndexBy "Mail"
// - IndexBy "OnPremisesSamAccountName"
// - IndexBy "PreferredName"
func NewUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	u := &Unique{
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		cs3conf: &Config{
			ProviderAddr: opts.ProviderAddr,
			JWTSecret:    opts.JWTSecret,
			ServiceUser:  opts.ServiceUser,
		},
	}

	return u
}

// Init initializes a unique index.
func (idx *Unique) Init() error {
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
func (idx *Unique) Lookup(v string) ([]string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	var claim string
	switch idx.indexBy {
	case "Mail":
		claim = "mail"
	case "OnPremisesSamAccountName", "PreferredName":
		claim = "username"
	case "UidNumber":
		claim = "uid"
	case "GidNumber":
		claim = "gid"
	case "Id":
		claim = "userid"
	default:
		return nil, fmt.Errorf("unsupported property %v", idx.indexBy)
	}

	res, err := idx.userProvider.GetUserByClaim(ctx, &user.GetUserByClaimRequest{
		Claim: claim,
		Value: v,
	})
	if err != nil {
		return nil, err
	}
	if res.Status.Code != rpcv1beta1.Code_CODE_OK {
		return []string{}, fmt.Errorf(res.Status.Message)
	}
	return []string{res.User.Id.OpaqueId + "@" + res.User.Id.Idp}, nil
	// TODO error ... unsupported index? what when we got more than one result?
}

// Add adds a value to the index, returns the path to the root-document
func (idx *Unique) Add(id, v string) (string, error) {
	return "", nil
}

// Remove a value v from an index.
func (idx *Unique) Remove(id string, v string) error {
	return &idxerrs.NotSupportedErr{}
}

// Update index from <oldV> to <newV>.
func (idx *Unique) Update(id, oldV, newV string) error {
	return &idxerrs.NotSupportedErr{}
}

// Search allows for glob search on the index.
// TODO cs3 really only allows
func (idx *Unique) Search(pattern string) ([]string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

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
		switch idx.indexBy {
		case "Id":
			if strings.Contains(res.Users[i].Id.OpaqueId, pattern) {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			}
		case "Mail":
			if strings.Contains(res.Users[i].Mail, pattern) {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			}
		case "OnPremisesSamAccountName", "PreferredName":
			if strings.Contains(res.Users[i].Username, pattern) {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			}
		case "DisplayName":
			if strings.Contains(res.Users[i].DisplayName, pattern) {
				matches = append(matches, res.Users[i].Id.OpaqueId+"@"+res.Users[i].Id.Idp)
			}
		case "UidNumber":
			// TODO check opaque property

		}
	}
	return matches, nil
	// TODO error ... insupported index?
	// for now igrnore other properties because the FindUsers already contains matches for username, display name, mail and opaqueid
}

// CaseInsensitive undocumented.
func (idx *Unique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *Unique) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *Unique) TypeName() string {
	return idx.typeName
}

// FilesDir undocumented.
func (idx *Unique) FilesDir() string {
	return ""
}

func (idx *Unique) authenticate(ctx context.Context) (token string, err error) {
	return storage.AuthenticateCS3(ctx, idx.cs3conf.ServiceUser, idx.tokenManager)
}

func (idx *Unique) getAuthenticatedContext(ctx context.Context) (context.Context, error) {
	t, err := idx.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	return ctx, nil
}

// Delete deletes the index folder from its storage.
func (idx *Unique) Delete() error {
	return &idxerrs.NotSupportedErr{}
}
