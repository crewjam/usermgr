package usermgr

import (
	"encoding/json"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestKey{})

type testRandomReader struct {
	Next byte
}

func (tr *testRandomReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = tr.Next
		tr.Next += 2
	}
	return len(p), nil
}

type TestKey struct {
}

func (s *TestKey) SetUpTest(c *C) {
	randReader = &testRandomReader{}
}

func (s *TestKey) TearDownTest(c *C) {
}

func (s *TestKey) TestAdminKey(c *C) {
	ak := AdminKey{}
	err := ak.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
	c.Assert(err, IsNil)
	c.Assert(ak, DeepEquals, AdminKey{
		HostKey: HostKey{
			AdminPublicKey: [32]uint8{0x9b, 0xf3, 0x62, 0xa8, 0xcc, 0x96, 0x92, 0x48, 0xe, 0x8b, 0x5b, 0x13, 0xe2, 0xe3, 0x2, 0x9e, 0x9e, 0x64, 0x62, 0xe3, 0x5a, 0x9d, 0xeb, 0x1c, 0x46, 0x44, 0x6b, 0xdc, 0x33, 0xf6, 0xf4, 0x55},
			HostPrivateKey: [32]uint8{0x0, 0x2, 0x4, 0x6, 0x8, 0xa, 0xc, 0xe, 0x10, 0x12, 0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e, 0x20, 0x22, 0x24, 0x26, 0x28, 0x2a, 0x2c, 0x2e, 0x30, 0x32, 0x34, 0x36, 0x38, 0x3a, 0x3c, 0x3e},
		},
		AdminPrivateKey: [32]uint8{0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54, 0x56, 0x58, 0x5a, 0x5c, 0x5e, 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 0x6c, 0x6e, 0x70, 0x72, 0x74, 0x76, 0x78, 0x7a, 0x7c, 0x7e},
		HostPublicKey:   [32]uint8{0xa2, 0x69, 0x90, 0x8f, 0x92, 0x89, 0xa1, 0xe1, 0xd1, 0x2e, 0x16, 0xc7, 0xc8, 0xd, 0x91, 0xcc, 0xd5, 0xc1, 0x78, 0x9f, 0xd7, 0xcf, 0x8a, 0x75, 0xbc, 0x95, 0x2c, 0xa3, 0x36, 0x73, 0x85, 0x7c},
	})

	text, err := ak.MarshalText()
	c.Assert(err, IsNil)
	c.Assert(string(text), Equals, "m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw")

	buf, err := json.Marshal(struct{ AK AdminKey }{AK: ak})
	c.Assert(err, IsNil)
	c.Assert(string(buf), Equals, "{\"AK\":\"m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw\"}")

	thing := struct{ AK AdminKey }{}
	err = json.Unmarshal(buf, &thing)
	c.Assert(err, IsNil)
	c.Assert(thing.AK, DeepEquals, ak)
}

func (s *TestKey) TestHostKey(c *C) {
	hk := HostKey{}
	err := hk.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg"))
	c.Assert(err, IsNil)
	c.Assert(hk, DeepEquals, HostKey{
		AdminPublicKey: [32]uint8{0x9b, 0xf3, 0x62, 0xa8, 0xcc, 0x96, 0x92, 0x48, 0xe, 0x8b, 0x5b, 0x13, 0xe2, 0xe3, 0x2, 0x9e, 0x9e, 0x64, 0x62, 0xe3, 0x5a, 0x9d, 0xeb, 0x1c, 0x46, 0x44, 0x6b, 0xdc, 0x33, 0xf6, 0xf4, 0x55},
		HostPrivateKey: [32]uint8{0x0, 0x2, 0x4, 0x6, 0x8, 0xa, 0xc, 0xe, 0x10, 0x12, 0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e, 0x20, 0x22, 0x24, 0x26, 0x28, 0x2a, 0x2c, 0x2e, 0x30, 0x32, 0x34, 0x36, 0x38, 0x3a, 0x3c, 0x3e},
	})

	text, err := hk.MarshalText()
	c.Assert(err, IsNil)
	c.Assert(string(text), Equals, "m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg")

	buf, err := json.Marshal(struct{ HK HostKey }{HK: hk})
	c.Assert(err, IsNil)
	c.Assert(string(buf), Equals, "{\"HK\":\"m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg\"}")

	thing := struct{ HK HostKey }{}
	err = json.Unmarshal(buf, &thing)
	c.Assert(err, IsNil)
	c.Assert(thing.HK, DeepEquals, hk)
}

func (s *TestKey) TestInvalidAdminKey(c *C) {
	ak := AdminKey{}
	err := ak.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg"))
	c.Assert(err, Equals, ErrIncorrectKeyFormat)

	err = ak.UnmarshalText([]byte("!!!!qMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
	c.Assert(err, Equals, ErrIncorrectKeyFormat)

	err = ak.UnmarshalJSON([]byte("!"))
	c.Assert(err, ErrorMatches, "invalid character.*")
}

func (s *TestKey) TestInvalidHostKey(c *C) {
	hk := HostKey{}
	err := hk.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
	c.Assert(err, Equals, ErrIncorrectKeyFormat)

	err = hk.UnmarshalText([]byte("!!!!qMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg"))
	c.Assert(err, Equals, ErrIncorrectKeyFormat)

	err = hk.UnmarshalJSON([]byte("!"))
	c.Assert(err, ErrorMatches, "invalid character.*")
}

func (s *TestKey) TestCanGenerateKey(c *C) {
	ak := GenerateKeyPair()
	c.Assert(ak.String(), Equals, "m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw")
	c.Assert(ak.HostKey.String(), Equals, "m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg")
}
