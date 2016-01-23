package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/crewjam/usermgr"
	. "gopkg.in/check.v1"
)

func (suite *TestWeb) TestCanGenerateTOTP(c *C) {
	suite.FakeAuth.Err = fmt.Errorf("not reached")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/_totp", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)

	var response struct {
		Secret string             `json:"secret"`
		Device usermgr.TOTPDevice `json:"device"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	c.Assert(err, IsNil)
	c.Assert(response.Secret, HasLen, 16)
	c.Assert(len(response.Device.Codes), Equals, 260)
}

func (suite *TestWeb) TestCanGenerateBackupCode(c *C) {
	suite.FakeAuth.Err = fmt.Errorf("not reached")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/_backup_code", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)

	var response struct {
		PlaintextCode string             `json:"plaintext_code"`
		BackupCode    usermgr.BackupCode `json:"backup_code"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	c.Assert(err, IsNil)
	c.Assert(response.PlaintextCode, HasLen, 16)
	c.Assert(len(response.BackupCode.Salt), Equals, 32)
}
