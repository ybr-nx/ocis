package storage

import (
	"context"
	"errors"
	"strconv"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis/ocis-pkg/log"
	"google.golang.org/grpc/metadata"
)

// CS3UsersRepo provides a cs3 users implementation of the Repo interface
// In contrast to the CS3Repo, which uses the storage provider, this implementation uses tdhge CS3 users api
type CS3UsersRepo struct {
	cfg        *config.Config
	tm         token.Manager
	userClient user.UserAPIClient
}

// NewCS3UsersRepo creates a new cs3 users repo
func NewCS3UsersRepo(cfg *config.Config) (Repo, error) {
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": cfg.TokenManager.JWTSecret,
	})

	if err != nil {
		return nil, err
	}

	client, err := pool.GetUserProviderServiceClient(cfg.Repo.CS3.UserProviderAddr)
	if err != nil {
		return nil, err
	}

	return CS3UsersRepo{
		cfg:        cfg,
		tm:         tokenManager,
		userClient: client,
	}, nil
}

// WriteAccount writes an account via cs3 and modifies the provided account (e.g. with a generated id).
func (r CS3UsersRepo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	return nil
}

// LoadAccount loads an account via cs3 by id and writes it to the provided account
func (r CS3UsersRepo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	// TODO only split at last @
	parts := strings.SplitN(id, "@", 2)
	if len(parts) < 2 {
		// make sure we always have 2 parts
		parts = append(parts, "")
	}
	resp, err := r.userClient.GetUser(ctx, &user.GetUserRequest{
		UserId: &user.UserId{
			OpaqueId: parts[0],
			Idp:      parts[1],
		},
	})
	switch {
	case err != nil:
		return err
	case resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND:
		return &notFoundErr{"account", id}
	case resp.Status.Code != v1beta11.Code_CODE_OK:
		return errors.New(resp.Status.Message)
	}

	return r.fillAccount(resp.User, a)
}

// LoadAccounts loads all the accounts from the cs3 api
func (r CS3UsersRepo) LoadAccounts(ctx context.Context, a *[]*proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	resp, err := r.userClient.FindUsers(ctx, &user.FindUsersRequest{})
	switch {
	case err != nil:
		return err
	case resp.Status.Code != v1beta11.Code_CODE_OK:
		return errors.New(resp.Status.Message)
	}

	// TODO get log from r
	log := olog.NewLogger(olog.Pretty(r.cfg.Log.Pretty), olog.Color(r.cfg.Log.Color), olog.Level(r.cfg.Log.Level))
	for i := range resp.Users {
		acc := &proto.Account{}
		if err = r.fillAccount(resp.Users[i], acc); err != nil {
			log.Err(err).Msg("could not load account")
			continue
		}
		*a = append(*a, acc)
	}
	return nil
}

func (r CS3UsersRepo) fillAccount(u *user.User, a *proto.Account) (err error) {
	// TODO Iss
	a.Id = u.Id.OpaqueId
	a.OnPremisesSamAccountName = u.Username
	a.PreferredName = u.Username
	a.Mail = u.Mail
	// TODO mail verified?
	a.DisplayName = u.DisplayName
	// TODO groups
	// map ids if available
	if u.Opaque != nil && u.Opaque.Map != nil {
		if uidObj, ok := u.Opaque.Map["uid"]; ok {
			if uidObj.Decoder == "plain" {
				if a.UidNumber, err = strconv.ParseInt(string(uidObj.Value), 10, 64); err != nil {
					return err
				}

			}
		}
		if gidObj, ok := u.Opaque.Map["gid"]; ok {
			if gidObj.Decoder == "plain" {
				if a.GidNumber, err = strconv.ParseInt(string(gidObj.Value), 10, 64); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// DeleteAccount deletes an account via cs3 users by id
func (r CS3UsersRepo) DeleteAccount(ctx context.Context, id string) (err error) {
	return &unsupportedErr{}
}

// WriteGroup writes a group via cs3 users and modifies the provided group (e.g. with a generated id).
func (r CS3UsersRepo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	return nil
}

// LoadGroup loads a group via cs3 users by id and writes it to the provided group
func (r CS3UsersRepo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	return &unsupportedErr{}
}

// LoadGroups loads all the groups from the cs3 users api
func (r CS3UsersRepo) LoadGroups(ctx context.Context, g *[]*proto.Group) (err error) {
	return &unsupportedErr{}
}

// DeleteGroup deletes a group via cs3 users by id
func (r CS3UsersRepo) DeleteGroup(ctx context.Context, id string) (err error) {
	return &unsupportedErr{}
}

func (r CS3UsersRepo) authenticate(ctx context.Context) (token string, err error) {
	return AuthenticateCS3(ctx, r.cfg.ServiceUser, r.tm)
}
