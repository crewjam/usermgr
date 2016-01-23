package usermgr

import (
	"time"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestTOTP{})

type TestTOTP struct {
	AdminKey AdminKey
}

func (s *TestTOTP) SetUpTest(c *C) {
	timeNow = func() time.Time {
		return time.Unix(1436983488, 0)
	}
	s.AdminKey.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
}

func (s *TestTOTP) TestValidate(c *C) {
	device := TOTPDevice{
		Name: "thing",
	}
	device.SetSecret(s.AdminKey, "WNZNSH53RSVUJWH2")
	err := device.GenerateCodes(timeNow().Add(-1*time.Minute),
		timeNow().Add(time.Minute), s.AdminKey)
	c.Assert(err, IsNil)

	// otp success
	err = device.VerifyCode(timeNow(), 30*time.Second, "066243")
	c.Assert(err, IsNil)

	// wrong code
	err = device.VerifyCode(timeNow(), 30*time.Second, "123456")
	c.Assert(err, Equals, ErrIncorrectCode)

	// expired
	timeNow = func() time.Time { return time.Unix(1000000000, 0) }
	err = device.VerifyCode(timeNow(), 30*time.Second, "066243")
	c.Assert(err, Equals, ErrIncorrectCode)
}

func (s *TestTOTP) TestCanDecryptKey(c *C) {
	device := TOTPDevice{
		Name: "thing",
	}
	device.SetSecret(s.AdminKey, "WNZNSH53RSVUJWH2")

	secret, err := device.Secret(s.AdminKey)
	c.Assert(err, IsNil)
	c.Assert(secret, Equals, "WNZNSH53RSVUJWH2")

	s.AdminKey.UnmarshalText([]byte("n_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXy"))
	secret, err = device.Secret(s.AdminKey)
	c.Assert(err, ErrorMatches, "cannot decrypt secret. Wrong key\\?")
	c.Assert(secret, Equals, "")

	err = device.GenerateCodes(timeNow().Add(-1*time.Minute),
		timeNow().Add(time.Minute), s.AdminKey)
	c.Assert(err, ErrorMatches, "cannot decrypt secret. Wrong key\\?")

	device.SetSecret(s.AdminKey, "wrongsecretencoding")
	err = device.GenerateCodes(timeNow().Add(-1*time.Minute),
		timeNow().Add(time.Minute), s.AdminKey)
	c.Assert(err, ErrorMatches, "invalid secret: illegal base32 data at input byte 0")
}
