package cmd

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/crewjam/usermgr"
)

// DefaultConfigPath is the default path to the configuration file.
var DefaultConfigPath string = "/etc/usermgr.conf"

// Config describes the host configuration
type Config struct {
	// The URL where the account database is stored. This URL can also
	// point to a storage service, i.e. "https://s3.amazonaws.com/example/users.pem"
	URL string

	// Specifies the host key used to decrypt and verify the database
	HostKey usermgr.HostKey

	// Specifies the path where a local copy of the account database is stored.
	// (Default: /var/lib/usermgr)
	CacheDir string

	// Specifies which groups a user must be part of in order to enable their
	// account.
	LoginGroups []string

	// Specifies which groups a user must be part of in order to enable them
	// to sudo to root.
	SudoGroups []string

	// If true then all remote users must specify an MFA token to login.
	LoginMFARequried bool
}

// LoadConfig returns a new config object by reading the file at path.
func LoadConfig(path string) (*Config, error) {
	config := Config{
		CacheDir:    "/var/lib/usermgr",
		LoginGroups: []string{"users"},
		SudoGroups:  []string{"wheel"},
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
