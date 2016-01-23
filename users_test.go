package usermgr

import (
	"bytes"
	"encoding/pem"
	"time"

	"golang.org/x/crypto/nacl/box"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestUsersData{})

type TestUsersData struct {
	AdminKey AdminKey
}

func (s *TestUsersData) SetUpTest(c *C) {
	timeNow = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		return t
	}
	s.AdminKey.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
}

func (s *TestUsersData) TearDownTest(c *C) {
}

func (s *TestUsersData) TestCanGetAndSet(c *C) {
	ud := UsersData{}

	c.Assert(ud.GetUserByName("alice"), IsNil)
	ud.Set(User{Name: "alice", RealName: "Alice Smith"})
	ud.Set(User{Name: "bob", RealName: "Bob Smith"})

	u := ud.GetUserByName("alice")
	c.Assert(u, DeepEquals, &User{Name: "alice", RealName: "Alice Smith"})

	// replacement works
	ud.Set(User{Name: "alice", RealName: "Alice Smith II"})
	c.Assert(ud.GetUserByName("alice"), DeepEquals, &User{Name: "alice", RealName: "Alice Smith II"})

	// delete works
	ud.Delete("alice")
	u = ud.GetUserByName("alice")
	c.Assert(u, IsNil)
}

func (s *TestUsersData) TestCanLoadAndStore(c *C) {
	ud := &UsersData{}
	ud.Set(User{
		Name:     "alice",
		RealName: "Alice Smith",
		AuthorizedKeys: []string{
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm",
		},
		BackupCodes: []BackupCode{BackupCode{
			// "hellohellohelloo"
			CreateTime: timeNow(),
			Salt:       []byte{0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54, 0x56, 0x58, 0x5a, 0x5c, 0x5e, 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 0x6c, 0x6e, 0x70, 0x72, 0x74, 0x76, 0x78, 0x7a, 0x7c, 0x7e},
			Hash:       []byte{0x3f, 0xec, 0x91, 0x9a, 0xfa, 0x22, 0xec, 0x68, 0xf0, 0xd2, 0xfa, 0xdd, 0xa2, 0x30, 0x9d, 0x64, 0x99, 0x55, 0xaa, 0xbe, 0xa2, 0xe6, 0xf5, 0x17, 0xcd, 0x6b, 0x87, 0xc4, 0xd5, 0xdd, 0xed, 0xc3, 0x96, 0xbc, 0x75, 0x90, 0x1b, 0xf7, 0x9a, 0xf1, 0x6f, 0xdc, 0xb3, 0x76, 0x1f, 0x56, 0x17, 0x2e, 0x10, 0x29, 0xa8, 0xce, 0x68, 0x46, 0x3c, 0x76, 0x31, 0xfe, 0xe9, 0x2e, 0xf7, 0xf2, 0x26, 0xa7},
		}},
	})
	signedData, err := ud.SignedString(s.AdminKey)
	c.Assert(err, IsNil)

	ud, err = LoadUsersData(signedData, s.AdminKey.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud, DeepEquals,
		&UsersData{Users: []User{
			User{
				Name:           "alice",
				RealName:       "Alice Smith",
				AuthorizedKeys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm"},
				BackupCodes: []BackupCode{BackupCode{
					// "hellohellohelloo"
					CreateTime: timeNow(),
					Salt:       []byte{0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54, 0x56, 0x58, 0x5a, 0x5c, 0x5e, 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 0x6c, 0x6e, 0x70, 0x72, 0x74, 0x76, 0x78, 0x7a, 0x7c, 0x7e},
					Hash:       []byte{0x3f, 0xec, 0x91, 0x9a, 0xfa, 0x22, 0xec, 0x68, 0xf0, 0xd2, 0xfa, 0xdd, 0xa2, 0x30, 0x9d, 0x64, 0x99, 0x55, 0xaa, 0xbe, 0xa2, 0xe6, 0xf5, 0x17, 0xcd, 0x6b, 0x87, 0xc4, 0xd5, 0xdd, 0xed, 0xc3, 0x96, 0xbc, 0x75, 0x90, 0x1b, 0xf7, 0x9a, 0xf1, 0x6f, 0xdc, 0xb3, 0x76, 0x1f, 0x56, 0x17, 0x2e, 0x10, 0x29, 0xa8, 0xce, 0x68, 0x46, 0x3c, 0x76, 0x31, 0xfe, 0xe9, 0x2e, 0xf7, 0xf2, 0x26, 0xa7},
				}},
			},
		}})

	wrongSignedData := bytes.Replace(signedData, []byte{'Q'}, []byte{'X'}, -1)
	_, err = LoadUsersData(wrongSignedData, s.AdminKey.HostKey)
	c.Assert(err, ErrorMatches, "cannot decrypt user data. Wrong key\\?")

	wrongSignedData = bytes.Replace(signedData, []byte("USERMGR"), []byte("FROB"), -1)
	_, err = LoadUsersData(wrongSignedData, s.AdminKey.HostKey)
	c.Assert(err, ErrorMatches, "invalid encoding")

	wrongSignedData = func() []byte {
		nonce := [24]byte{}
		randReader.Read(nonce[:])
		ciphertext := nonce[:]
		ciphertext = box.Seal(ciphertext, []byte("invalid json"),
			&nonce, &s.AdminKey.HostPublicKey,
			&s.AdminKey.AdminPrivateKey)
		buf := bytes.NewBuffer(nil)
		err = pem.Encode(buf, &pem.Block{
			Type:  "USERMGR DATA",
			Bytes: ciphertext,
		})
		return buf.Bytes()
	}()
	_, err = LoadUsersData(wrongSignedData, s.AdminKey.HostKey)
	c.Assert(err, ErrorMatches, "invalid character 'i' looking for beginning of value")

	s.AdminKey.HostKey.AdminPublicKey[0] = 0x99
	_, err = LoadUsersData(signedData, s.AdminKey.HostKey)
	c.Assert(err, ErrorMatches, "cannot decrypt user data. Wrong key\\?")

}
