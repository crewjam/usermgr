package usermgr

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

// UsersData is the document that contains all the authoritative
// data about users, their keys, group membership, etc.
type UsersData struct {
	Users               []User `json:"users"`
	YubikeyClientID     string `json:"yubikey_client_id,omitempty"`
	YubikeyClientSecret string `json:"yubikey_client_secret,omitempty"`
}

// GetUserByName returns the user having the specified name or
// nil if no such user exists.
func (ud UsersData) GetUserByName(name string) *User {
	for _, user := range ud.Users {
		if user.Name == name {
			return &user
		}
	}
	return nil
}

// Set adds or replaces `user` to the list of users.
func (ud *UsersData) Set(user User) {
	found := false
	newUsers := []User{}
	for _, u := range ud.Users {
		if u.Name == user.Name {
			newUsers = append(newUsers, user)
			found = true
		} else {
			newUsers = append(newUsers, u)
		}
	}
	if !found {
		ud.Users = append(ud.Users, user)
	} else {
		ud.Users = newUsers
	}
}

// Delete removes a user from the list of users
func (ud *UsersData) Delete(userName string) {
	newUsers := []User{}
	for _, u := range ud.Users {
		if u.Name == userName {
			continue
		}
		newUsers = append(newUsers, u)
	}
	ud.Users = newUsers
}

// SignedString returns a serialized version of the user database
// signed with the given key.
func (ud UsersData) SignedString(adminKey AdminKey) ([]byte, error) {
	unsignedData, err := json.Marshal(ud)
	if err != nil {
		return nil, err
	}

	nonce := [24]byte{}
	randReader.Read(nonce[:])
	ciphertext := nonce[:]
	ciphertext = box.Seal(ciphertext, unsignedData, &nonce, &adminKey.HostPublicKey,
		&adminKey.AdminPrivateKey)

	buf := bytes.NewBuffer(nil)
	err = pem.Encode(buf, &pem.Block{
		Type:  "USERMGR DATA",
		Bytes: ciphertext,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// LoadUsersData reads signedData, verifies that it was signed by
// the specified publicKey and if all is well, returns a new
// instance of UserData.
func LoadUsersData(data []byte, hostKey HostKey) (*UsersData, error) {
	dataBlock, _ := pem.Decode(data)
	if dataBlock == nil || dataBlock.Type != "USERMGR DATA" {
		return nil, fmt.Errorf("invalid encoding")
	}

	nonce := [24]byte{}
	copy(nonce[:], dataBlock.Bytes[:24])

	plaintext, ok := box.Open(nil, dataBlock.Bytes[24:], &nonce,
		&hostKey.AdminPublicKey, &hostKey.HostPrivateKey)
	if !ok {
		return nil, fmt.Errorf("cannot decrypt user data. Wrong key?")
	}

	ud := UsersData{}
	if err := json.Unmarshal(plaintext, &ud); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return &ud, nil
}
