# usermgr

Managing administrative access to production systems is a pain. We all know that having accurate administrative access
is important for security. Create a accounts too slowly and you incentivize account sharing. Disable accounts too
slowly and you risk violating least-privilege (or worse). At the same time, centralized, online systems can represent a single
point of failure for your application.

Usermgr is a tool to turn access to production systems from a pain in the butt into ponies and rainbows.

<img src="http://img02.deviantart.net/34a9/i/2012/289/3/9/rainbow_dash_kissing_by_teiptr-d5i00jb.png" width="200">

Usermgr access rules are determined by a single signed file that 
can be cached locally so it doesn't create a dependency on any central 
system.

Problems it fixes:

 - Creating local unix accounts
 - Managing user's SSH authorized keys
 - Two factor authentication with Yubikey or TOTP
 - Shell logging
 - Managing the sudoers file

A web interface and command line tools are available for managing users, enrollment and user self-service.
Each node in the environment has access to the account database which is cryptographically signed and centrally managed.

## Getting started

### 0. Install 

Download an appropriate binary for your system:

    # curl -o /opt/usermgr/bin/usermgr https://github.com/crewjam/usermgr/releases/download/XXX/usermgr.$(uname -s).$(uname -m) 
    # chmod +x /opt/usermgr/bin/usermgr

From source:

    go install github.com/crewjam/usermgr

### 1. Generate a key pair

   The account database can be stored in Amazon S3 or on the local file system.
   (use the `file://` url-scheme). `usermgr` uses the username and password fields of the URL to store the public and private keys for the database, respectively.

    $  usermgr keygen
    admin key: Ulc7w67dHOagHVBWf18fmTAAOCs3dG0mql0NTTjDP2xQHNgZQjAo6Oy2aJie89TdOR10vg-cx-d0POwpm8tB5K-FMguXPr8b_zS3_fvTW1k16IMbs_aCoQ8u82eLcyB8A_CwAvsoRCVGmMzzBRtMJtquskeEMidS6AGMDvcteDc
    host key: Ulc7w67dHOagHVBWf18fmTAAOCs3dG0mql0NTTjDP2xQHNgZQjAo6Oy2aJie89TdOR10vg-cx-d0POwpm8tB5A

   This produces two keys. Posession of the admin key allows you to edit the file and resign it. Posession of the host key allows you to read the file, but not edit it. All your servers will have the host key, but generally only the one running the web interface will have the admin key.

### 2. Run the web interface (optional)

   The web interface requires web users to be authenticated with an external mechanism. You can use `oauth` or `header`. 

   For header authentication, another server (i.e. Apache or nginx) handles the authentication and places the user name in the `X-Remote-User` header.
   For OAuth authentication, you must provide parameters for the OAuth provider. 

     $ UM_ADMIN_KEY=Ulc7w67dHOagHVBWf18fmTAAOCs3dG0mql0NTTjDP2xQHNgZQjAo6Oy2aJie89TdOR10vg-cx-d0POwpm8tB5K-FMguXPr8b_zS3_fvTW1k16IMbs_aCoQ8u82eLcyB8A_CwAvsoRCVGmMzzBRtMJtquskeEMidS6AGMDvcteDc \
     UM_STORE=file:///var/usermgr/ \
     UM_AUTH=oauth://google?client_id=XXX.apps.googleusercontent.com&client_secret=xYxYxY&email_suffix=@example.com \
     UM_TOKEN_KEY=someRandomThing \
     UM_URL=https://users.example.com \
     usermgr web

   The web interface will automatically create users the first time they navigate to the web interface. Those newly created users will not be part of any groups, so they won't have access to any systems. (The first user is automatically added to the `usermgr-admin` group. This is the only special group. Users that are members of this group are allowed to create or destroy users, edit group membership, and modify other users besides themselves.)

   The only storage scheme supported is `file`. (Your contributions in this area are welcome!)

   The only auth scheme supported in `oauth`. (Your contributions in this area are welcome!). The following query parameters are supported for oauth:

   - `client_id` - The OAuth2 client id.
   - `client_secret` - The OAuth2 client secret.
   - `scope` - which OAuth2 scopes to request. Specify multiple times for multiple scopes. The default is `openid`, `profile` and `email`. 
   - `auth_url` - The URL where OAuth2 auth requests are sent
   - `token_url` - The URL where OAuth2 tokens are produced
   - `user_info_url` - An OpenID-compatible user info request URL.
   - `email_suffix` - If present require that the email address end with the specified suffix.
   
   If you specify `google` as the hostname in the auth URL, then the default `auth_url`, `token_url` and `user_info_url` are filled in for you.

## Setting up Managed Systems

On each system that you want to manage you'll need to install `usermgr` and configure the system.

### 1. Configure `usermgr`

   On managed systems you do not use the editing URL, instead use the read-only URL which contains only the public key. Create a file called `/etc/usermgr.conf` that looks something like this:

      URL = "https://users.example.com"
      HostKey = "m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg"
      LoginMFARequried = false

### 2. Configure the `sshd`

   Tell `sshd` to ignore the user's `~/.ssh/authorized_keys` file and instead invoke a command to determine which keys to use.

    echo "AuthorizedKeysCommand /opt/bin/usermgr.sshkeys" >> /etc/ssh/sshd_config
    echo "AuthorizedKeysCommandUser nobody" >> /etc/ssh/sshd_config
    ln -s /opt/bin/usermgr /opt/bin/usermgr.sshkeys

   Note: when you invoke `usermgr` as `usermgr.sshkeys` it is equivalent to invoking `usermgr sshkeys`. This mini-hack is needed because `AuthorizedKeysCommand` expects a single binary, not a command line.

### 3. Configure usermgr to be a login shell.
  
   If you want to use multi-factor authentication or logging, you can replace each user's login shell with `usermgr`
    
    (
      echo 'Defaults log_output'
      echo 'Defaults!/usr/bin/sudoreplay !log_output'
      echo 'Defaults!/sbin/reboot !log_output'
      echo 'Defaults env_keep += "SSH_CLIENT SSH_CONNECTION"'
    ) > /etc/sudoers.d/replay
    ln -s /opt/bin/usermgr /opt/bin/usermgr.shell
    chsh -s /opt/bin/usermgr.shell bob

   `usermgr shell` enforces the second authentication factor. It also enforces shell-logging by invoking `bash` via `sudo` to the current user. (Note: invoking `sudo` in this context doesn't change the user's privilege level, it only allows terminal logging to happen.)

### 4. Arrange for `usermgr sync` to run every couple of minutes.

   This keeps the local copy of the account database 
   up to date, creates or removes local users, and keeps the sudoers file in sync.

    echo '*/5 * * * /opt/usermgr/bin/usermgr.sync' > /etc/cron.d/usermgr

## Multi-factor Authentication

Usermgr supports [yubikey](https://www.yubico.com/products/yubikey-hardware/), and [Google Authenticator](https://support.google.com/accounts/answer/1066447?hl=en) (TOTP) for multi-factor authentication. You can (should!) also create backup codes that allow you to connect in the event of a failure of the authentication service or your multi-factor device.

### Yubikey

We support yubikeys in the (default) Yubikey-OTP mode.

Register your yubikey by trying it out at https://demo.yubico.com/. If it isn't use the "Upload to Yubico" button in the Yubico Personalization Tool (or some other means, as documented) so that yubicloud knows about your key.

In the web interface, click "Enroll Yubikey" and press the button on your key. 

Yubikeys are validated online by each server at login time. The Yubico client ID and secret must be entered in the web interface and are then available to 
each host.

### Google Authenticator / TOTP

Install the Google Authenticator app for [Android](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en) or [iOS](https://itunes.apple.com/us/app/google-authenticator/id388497605?mt=8). 

In the web interface, click "Enroll Smartphone" and scan the QR code with the Google Authenticator app. 

TOTP requires a secret be shared between the system doing the authenticating and the device (mobile phone, etc.). We'd prefer **not** to share the TOTP secret with every host you log in to, only the auth server. At the same time we **must not** require that the central server be online to authenticate. So instead, on the auth server, usermgr pregenerates a bunch of TOTP auth codes for every user and store the resulting hashes in `users.pem` which is distributed to all the hosts. Thus the hosts can authenticate TOTP codes by comparing hashes even if the auth server is offline for a short time. If the auth server is offline long enough that the host runs out of codes, then TOTP authentication will fail. By default, every hour, two hours worth of codes are generated.

### Backup Codes

Click "Generate Backup Code". The web interface will display a code and give you a chance to copy it down. Put it in a safe place, this is the only chance you'll have to use the backup code.

Backup codes are hashed before they are stored in `users.pem`.

# Configuration Reference

Here is a commented example configuration file:

    # The URL where the account database is stored. This URL can also
    # point to a storage service, i.e. "https://s3.amazonaws.com/example/users.pem"
    URL = "https://users.example.com/users.pem"
    
    # Specifies the host key used to decrypt and verify the database
    HostKey = ""

    # Specifies the path where a local copy of the account database is stored.
    # (Default: /var/lib/usermgr)
    CacheDir = "/var/lib/usermgr"

    # Specifies which groups a user must be part of in order to enable their 
    # account. Comma separated list. (Default: users)
    LoginGroups = "myapp-admin,myapp-user"

    # Specifies which groups a user must be part of in order to enable them
    # to sudo to root. Comma separated list. (Default: wheel)
    SudoGroups = "myapp-admin"

    # If true then all remote users must specify an MFA token to login.
    # (Default: false)
    LoginMFARequried = true

# FAQ

## How should I secure `users.pem`?

The file is encrypted such that only holders of the host key or admin key can read it. You should probably treat this file a bit like you would treat 
`/etc/shadow` -- private, but not secret. Even if an attacker had posession of the file and the host key, it does not contain any directly usable credentials, only one-way hashes of credentials. 

## How should I secure my host key.

It is private but not secret, like `/etc/shadow`. The host key unlocks access to `users.pem` which would give an attacker access to a list of your users, their privileges and public keys. But it does not contain directly usable credentials and thus if it fell into the wrong hands would not allow an attacker additional access.

## How should I secure my admin key.

This is the key to your kingdom, so you should protect it well. The admin key allows editing the user database, which could be used to add new user accounts or grant additional privileges to user accounts. The admin key also allows access to TOTP secrets. 

## What cryptographic algorithms are used?

* To sign the user database, [ED25519](http://ed25519.cr.yp.to/)
* To encrypt the secrets in the user database, we use [NaCL](https://godoc.org/golang.org/x/crypto/nacl/secretbox) which uses XSalsa20 and Poly1305.
* To hash the backup keys, [Scrypt](http://www.tarsnap.com/scrypt/scrypt.pdf).

## What happens if the web server does down?

You cannot use it to modify the accounts database any more. Hosts that have `users.pem` cached locally will continue to use it.

Eventually, the pregenerated TOTP codes will expire and TOTP authentication will stop working. Yubikey authentication requires access to yubikey's servers so it will continue to work. Backup codes continue to work.
