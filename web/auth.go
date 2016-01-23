package web

import (
	"net/http"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"

	"golang.org/x/net/context"
)

type Auth interface {
	RequireUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error)
}

func (s *Server) RequireUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (*RemoteUser, error) {
	remoteUserName, err := s.Auth.RequireUser(ctx, w, r)
	if err != nil {
		return nil, err
	}

	usersData, _, err := s.loadData(ctx)
	if err != nil {
		return nil, err
	}

	user := usersData.GetUserByName(remoteUserName)
	if user == nil {
		// auto create the user
		user = &usermgr.User{
			Name: remoteUserName,
		}
		if len(usersData.Users) == 0 {
			// First user is automatically an admin
			user.Groups = []string{"usermgr-admin"}
		}
		usersData.Set(*user)

		if err := s.storeData(ctx, usersData); err != nil {
			return nil, err
		}
	}

	return &RemoteUser{Name: remoteUserName, IsAdmin: user.InGroup("usermgr-admin")}, nil
}

func (s *Server) requireUserOrAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*RemoteUser, error) {
	remoteUser, err := s.RequireUser(ctx, w, r)
	if err != nil {
		return nil, err
	}
	if remoteUser.Name != Param(ctx, "user") && !remoteUser.IsAdmin {
		return nil, httperr.Forbidden
	}
	return remoteUser, nil
}
