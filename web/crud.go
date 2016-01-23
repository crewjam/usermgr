package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"
	"golang.org/x/net/context"
)

func (s *Server) getUsersList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	remoteUser, err := s.RequireUser(ctx, w, r)
	if err != nil {
		return err
	}

	usersData, etag, err := s.loadData(ctx)
	if err != nil {
		return err
	}

	if !remoteUser.IsAdmin {
		user := usersData.GetUserByName(remoteUser.Name)
		usersData.Users = []usermgr.User{*user}
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(usersData)
	return nil
}

func (s *Server) getUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := s.requireUserOrAdmin(ctx, w, r)
	if err != nil {
		return err
	}

	usersData, etag, err := s.loadData(ctx)
	if err != nil {
		return err
	}

	user := usersData.GetUserByName(Param(ctx, "user"))
	if user == nil {
		return httperr.NotFound
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (s *Server) putUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	remoteUser, err := s.requireUserOrAdmin(ctx, w, r)
	if err != nil {
		return err
	}

	user := usermgr.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return httperr.BadRequest
	}
	if user.Name != "" && user.Name != Param(ctx, "user") {
		return httperr.Error{
			StatusCode:   http.StatusBadRequest,
			PrivateError: fmt.Errorf("user name in request body does not match URL"),
		}
	}
	user.Name = Param(ctx, "user")

	err = s.mutateUsersData(ctx, func(usersData *usermgr.UsersData) error {
		if !remoteUser.IsAdmin {
			existingUser := usersData.GetUserByName(Param(ctx, "user"))
			if existingUser == nil {
				return httperr.Forbidden
			}

			// forbid the user from adding a group
			for _, group := range user.Groups {
				if !existingUser.InGroup(group) {
					return httperr.Error{
						StatusCode:   http.StatusForbidden,
						PrivateError: fmt.Errorf("non-admin user cannot add groups"),
					}
				}
			}
		}
		usersData.Set(user)
		return nil
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) deleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	remoteUser, err := s.RequireUser(ctx, w, r)
	if err != nil {
		return err
	}
	if !remoteUser.IsAdmin {
		return httperr.Forbidden
	}
	err = s.mutateUsersData(ctx, func(usersData *usermgr.UsersData) error {
		usersData.Delete(Param(ctx, "user"))
		return nil
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil

}
