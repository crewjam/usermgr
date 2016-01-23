package web

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/crewjam/httperr"

	. "gopkg.in/check.v1"
)

func (suite *TestWeb) TestCanGetUsers(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)
	c.Assert(w.Header(), DeepEquals, http.Header{
		"Etag":         []string{"34cfb4411e1ed35d1183e544202c0608e3c91c0c"},
		"Content-Type": []string{"application/json"},
	})
	c.Assert(string(w.Body.Bytes()), Equals, "{\"users\":[{\"name\":\"alice\",\"real_name\":\"Alice Smith\",\"groups\":[\"usermgr-admin\"],\"authorized_keys\":[\"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm\"],\"backup_codes\":[{\"create_time\":\"2006-01-02T15:04:05Z\",\"salt\":\"QEJERkhKTE5QUlRWWFpcXmBiZGZoamxucHJ0dnh6fH4=\",\"hash\":\"P+yRmvoi7Gjw0vrdojCdZJlVqr6i5vUXzWuHxNXd7cOWvHWQG/ea8W/cs3YfVhcuECmozmhGPHYx/uku9/Impw==\"}]},{\"name\":\"bob\",\"groups\":[\"wheel\"]}]}\n")
}

func (suite *TestWeb) TestUsersNonAdmin(c *C) {
	suite.FakeAuth.User = "bob"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)
	c.Assert(string(w.Body.Bytes()), Equals, "{\"users\":[{\"name\":\"bob\",\"groups\":[\"wheel\"]}]}\n")
}

func (suite *TestWeb) TestUsersRequiresAuth(c *C) {
	suite.FakeAuth.User = ""
	suite.FakeAuth.Err = httperr.Forbidden
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)
}

func (suite *TestWeb) TestUsersStorageFail(c *C) {
	suite.FakeStorage.Data = nil
	suite.FakeStorage.Etag = ""
	suite.FakeStorage.Err = fmt.Errorf("cannot frob the grob")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusInternalServerError)
	c.Assert(string(w.Body.Bytes()), Equals, "Internal Server Error\n")
}

func (suite *TestWeb) TestGetUserNonAdmin(c *C) {
	suite.FakeAuth.User = "bob"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)
	c.Assert(string(w.Body.Bytes()), Equals, "{\"name\":\"bob\",\"groups\":[\"wheel\"]}\n")
}

func (suite *TestWeb) TestGetOtherUserNonAdmin(c *C) {
	suite.FakeAuth.User = "bob"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/alice", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)
}

func (suite *TestWeb) TestGetUserAdmin(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)
	c.Assert(string(w.Body.Bytes()), Equals, "{\"name\":\"bob\",\"groups\":[\"wheel\"]}\n")
}

func (suite *TestWeb) TestGetUserAdminNotFound(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/charlie", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusNotFound)
}

func (suite *TestWeb) TestGetUserNonAdminNotFound(c *C) {
	suite.FakeAuth.User = "bob"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/charlie", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)
}

func (suite *TestWeb) TestPutUserNonAdmin(c *C) {
	suite.FakeAuth.User = "bob"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{\"name\":\"bob\",\"real_name\": \"Bob Smith\", \"groups\":[\"wheel\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	log.Printf("%v", w.Header())
	c.Assert(w.Code, Equals, 204)
	c.Assert(string(w.Body.Bytes()), Equals, "")

	// I can't change my user name
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{\"name\":\"charlie\",\"groups\":[\"wheel\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	log.Printf("%v", w.Header())
	c.Assert(w.Code, Equals, http.StatusBadRequest)

	// bad content
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{"))
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusBadRequest)

	// I can't add groups
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{\"name\":\"bob\",\"groups\":[\"wheel\",\"usermgr-admin\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)

	// I can't edit other users
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/alice",
		strings.NewReader("{\"name\":\"alice\",\"groups\":[\"wheel\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)
}

func (suite *TestWeb) TestPutUserAdmin(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{\"name\":\"bob\",\"real_name\": \"Bob Smith\", \"groups\":[\"wheel\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	log.Printf("%v", w.Header())
	c.Assert(w.Code, Equals, 204)
	c.Assert(string(w.Body.Bytes()), Equals, "")

	// I can add groups, and edit other users
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/bob",
		strings.NewReader("{\"name\":\"bob\",\"groups\":[\"wheel\",\"usermgr-admin\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 204)

	// I can create users
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("PUT", "/users/charlie",
		strings.NewReader("{\"name\":\"charlie\",\"groups\":[\"wheel\",\"usermgr-admin\"]}"))
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 204)
}

func (suite *TestWeb) TestDeleteUserAdmin(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	log.Printf("%v", w.Header())
	c.Assert(w.Code, Equals, 204)
	c.Assert(string(w.Body.Bytes()), Equals, "")
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 404)

	// non-admin cannot delete

	suite.FakeAuth.User = "bob"
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("DELETE", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)

	// requires auth
	suite.FakeAuth.User = ""
	suite.FakeAuth.Err = httperr.Unauthorized
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("DELETE", "/users/bob", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusUnauthorized)
}
