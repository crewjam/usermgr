package usermgr

import "time"

// User represents a single user
type User struct {
	Name           string          `json:"name"`
	RealName       string          `json:"real_name,omitempty"`
	Email          string          `json:"email,omitempty"`
	Groups         []string        `json:"groups,omitempty"`
	AuthorizedKeys []string        `json:"authorized_keys,omitempty"`
	Yubikeys       []YubikeyDevice `json:"yubikeys,omitempty"`
	BackupCodes    []BackupCode    `json:"backup_codes,omitempty"`
	TOTPDevices    []TOTPDevice    `json:"totp_devices,omitempty"`
}

// InGroup returns true if the user is a member of the specified group
func (u User) InGroup(groupName string) bool {
	for _, g := range u.Groups {
		if g == groupName {
			return true
		}
	}
	return false
}

// InAnyGroup returns true if the user is a member of the specified group
func (u User) InAnyGroup(groupNames []string) bool {
	for _, groupName := range groupNames {
		if u.InGroup(groupName) {
			return true
		}
	}
	return false
}

type YubikeyDevice struct {
	Name       string    `json:"name,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	DeviceID   string    `json:"device_id,omitempty"`
}
