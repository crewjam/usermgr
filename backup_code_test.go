package usermgr

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestBackupCode{})

type TestBackupCode struct {
}

func (s *TestBackupCode) TestCanHashAndCompare(c *C) {
	bc := NewBackupCode("Wyja8OeSygqwT9v8")
	c.Assert(bc.Matches("Wyja8OeSygqwT9v8"), Equals, true)
	c.Assert(bc.Matches("xxxxxxxxxxxxxxxx"), Equals, false)

	bc = BackupCode{}
	bc.Salt, _ = hex.DecodeString("51e1153f96592e73ab8dbf93c12b20a79c1c19465e4cac05444fc3265fe9bb70")
	bc.Hash, _ = hex.DecodeString("4c06d54904fb520ea1e54c74edf02a6295d707a21fa65f287c3e028ce6e790ce64df39dad12af907bae4d32d260cadc8ba9c2e8cd1a3c2364299ed37b3868b0b")

	c.Assert(bc.Matches("Wyja8OeSygqwT9v8"), Equals, true)
	c.Assert(bc.Matches("xxxxxxxxxxxxxxxx"), Equals, false)

	bc = BackupCode{}
	bc.Salt, _ = hex.DecodeString("51")
	bc.Hash, _ = hex.DecodeString("4c06d54904fb520ea1e54c74edf02a6295d707a21fa65f287c3e028ce6e790ce64df39dad12af907bae4d32d260cadc8ba9c2e8cd1a3c2364299ed37b3868b0b")
	c.Assert(bc.Matches("Wyja8OeSygqwT9v8"), Equals, false)
}
