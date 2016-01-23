package usermgr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// UpdateLocalCache fetches the user data from upstreamURL if it is unchanged.
// If the response is valid, the cache files in path are replaced and the
// new data are returned.
func UpdateLocalCache(path string, upstreamURL string, hostKey HostKey) (*UsersData, error) {
	var userData *UsersData

	etag := ""
	etagBuf, err := ioutil.ReadFile(filepath.Join(path, "users.pem.etag"))
	if os.IsNotExist(err) {
		etag = ""
	} else if err != nil {
		return nil, err
	} else {
		// Check that the existing data are valid. If not, then we set the
		// etag to "" so we can do an unconditional fetch.
		dataBuf, err := ioutil.ReadFile(filepath.Join(path, "users.pem"))
		if err == nil {
			userData, err = LoadUsersData(dataBuf, hostKey)
			if err == nil {
				etag = string(etagBuf)
			}
		}
	}

	// make an HTTP request to fetch the new data
	req, err := http.NewRequest("GET", upstreamURL, nil)
	if err != nil {
		return nil, err
	}
	if etag != "" {
		req.Header.Add("If-None-Match", etag)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotModified && userData != nil {
		// The response is that the file is unchanged, so we just return
		// the parsed, cached data.
		return userData, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	// Verify that the response contains a valid message
	dataBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userData, err = LoadUsersData(dataBuf, hostKey)
	if err != nil {
		return nil, err
	}

	// make sure the output directory exists. The error can be safely ignored,
	// because it will be either because the directory already exists or a
	// subsequent write will fail.
	os.MkdirAll(path, 0755)

	// Write the response (very carefully) to the cache location
	if err := ioutil.WriteFile(filepath.Join(path, "users.pem~"), dataBuf, 0644); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(path, "users.pem.etag"),
		[]byte(resp.Header.Get("ETag")), 0644); err != nil {
		os.Remove(filepath.Join(path, "users.pem~"))
		return nil, err
	}
	if err := os.Rename(filepath.Join(path, "users.pem~"),
		filepath.Join(path, "users.pem")); err != nil {
		os.Remove(filepath.Join(path, "users.pem~"))
		os.Remove(filepath.Join(path, "users.pem.etag"))
		return nil, err
	}
	return userData, nil
}

// GetLocalCache returns the local cached data in path if it is valid. It
// does not attempt to update the cache.
func GetLocalCache(path string, hostKey HostKey) (*UsersData, error) {
	dataBuf, err := ioutil.ReadFile(filepath.Join(path, "users.pem"))
	if err != nil {
		return nil, fmt.Errorf("Cannot read users data: %s", err)
	}

	userData, err := LoadUsersData(dataBuf, hostKey)
	if err != nil {
		return nil, err
	}

	return userData, nil
}
