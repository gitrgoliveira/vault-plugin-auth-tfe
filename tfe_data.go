package tfeauth

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

type RunInfo struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			Status string `json:"status"`
		} `json:"attributes"`
		Relationships struct {
			Workspace struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"workspace"`
		} `json:"relationships"`
	} `json:"data"`
}

func fetchRunInfo(login *tfeLogin, config *tfeConfig) (*RunInfo, error) {
	client := cleanhttp.DefaultClient()
	if len(config.CACert) > 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(config.CACert))

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
		}

		client.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	}

	url := fmt.Sprintf("https://%s/api/v2/runs/%s", strings.TrimSuffix(config.Host, "/"), login.RunID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bearer := fmt.Sprintf("Bearer %s", login.AtlasToken)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If the request was not a success log an error
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return nil, fmt.Errorf(string(body))
	}

	runResp := &RunInfo{}
	err = json.Unmarshal(body, runResp)
	if err != nil {
		return nil, err
	}

	return runResp, nil
}

type WorkspaceInfo struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name   string `json:"name"`
			Locked bool   `json:"locked"`
		} `json:"attributes"`
		Relationships struct {
			Organization struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"organization"`
			LockedBy struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"locked-by"`
			CurrentRun struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"current-run"`
			LatestRun struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"latest-run"`
		} `json:"relationships"`
	} `json:"data"`
}

func fetchWorkspaceInfo(login *tfeLogin, config *tfeConfig) (*WorkspaceInfo, error) {
	client := cleanhttp.DefaultClient()
	if len(config.CACert) > 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(config.CACert))

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
		}

		client.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	}

	url := fmt.Sprintf("https://%s/api/v2/organizations/%s/workspaces/%s",
		strings.TrimSuffix(config.Host, "/"), config.Organization, login.Workspace)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bearer := fmt.Sprintf("Bearer %s", login.AtlasToken)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If the request was not a success log an error
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return nil, fmt.Errorf(string(body))
	}

	workspaceInfo := &WorkspaceInfo{}
	err = json.Unmarshal(body, workspaceInfo)
	if err != nil {
		return nil, err
	}

	return workspaceInfo, nil
}

type AccountInfo struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Username         string `json:"username"`
			IsServiceAccount bool   `json:"is-service-account"`
		} `json:"attributes"`
		Relationships struct {
			AuthenticationTokens struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"authentication-tokens"`
		} `json:"relationships"`
	} `json:"data"`
}

func fetchAccountInfo(login *tfeLogin, config *tfeConfig) (*AccountInfo, error) {
	client := cleanhttp.DefaultClient()
	if len(config.CACert) > 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(config.CACert))

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
		}

		client.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	}

	url := fmt.Sprintf("https://%s/api/v2/account/details",
		strings.TrimSuffix(config.Host, "/"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bearer := fmt.Sprintf("Bearer %s", login.AtlasToken)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If the request was not a success log an error
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return nil, fmt.Errorf(string(body))
	}

	accountInfo := &AccountInfo{}
	err = json.Unmarshal(body, accountInfo)
	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}
