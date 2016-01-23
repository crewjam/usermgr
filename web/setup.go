package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"

	"golang.org/x/net/context"
)

func (s *Server) getSetup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// important: we cannot use s.AdminKey.HostKey here because the
	// remote user has not been authenticated. They have to prove they
	// know the host key by giving it to us in the URL
	hostKey := Param(ctx, "key")
	if hostKey != "" {
		hk := usermgr.HostKey{}
		if err := hk.UnmarshalText([]byte(hostKey)); err != nil {
			return httperr.Error{
				StatusCode:   http.StatusBadRequest,
				PrivateError: err,
			}
		}
	}

	adminURL := url.URL{
		Scheme: "https",
		Host:   r.Host,
		Path:   "/users.pem",
	}
	fmt.Fprintf(w, `#!/bin/sh
set -ex

# fetch the binary
[ -d /opt/usermgr/bin ] || mkdir -p /opt/usermgr/bin
curl -o /opt/usermgr/bin/usermgr %s.$(uname -s).$(uname -m)
chmod +x /opt/usermgr/bin/usermgr

# write the configuration file
(
	echo 'HostKey = "%s"'
	echo 'URL = "%s"' >> /etc/usermgr.conf
) > /etc/usermgr.conf

# do the initial sync and install a cronjob to keep it
# updated.
/opt/usermgr/bin/usermgr sync
echo '*/5 * * * /opt/usermgr/bin/usermgr sync' > /etc/cron.d/usermgr

# configure sshd
echo "AuthorizedKeysCommand /opt/usermgr/bin/usermgr.authorized-keys" >> /etc/ssh/sshd_config
echo "AuthorizedKeysCommandUser nobody" >> /etc/ssh/sshd_config
ln -s /opt/usermgr/bin/usermgr /opt/usermgr/bin/usermgr.authorized-keys

# configure sudo
(
  echo 'Defaults log_output'
  echo 'Defaults!/usr/bin/sudoreplay !log_output'
  echo 'Defaults!/sbin/reboot !log_output'
  echo 'Defaults env_keep += "SSH_CLIENT SSH_CONNECTION"'
) > /etc/sudoers.d/replay
if ! grep /opt/usermgr/bin/usermgr.shell /etc/shells >/dev/null ; then
	echo /opt/usermgr/bin/usermgr.shell >> /etc/shells
fi
ln -s /opt/usermgr/bin/usermgr /opt/usermgr/bin/usermgr.shell

`, s.DownloadURL, hostKey, adminURL.String())
	return nil
}
