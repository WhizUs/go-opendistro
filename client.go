// Copyright 2019 WhizUs GmbH. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package opendistro

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/WhizUs/go-opendistro/common"
	"github.com/WhizUs/go-opendistro/security"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-rootcerts"
	"io/ioutil"
	"log"
	"net/http"
)

type ClientConfig struct {
	Username, Password, BaseURL string

	TLSConfig *TLSConfig
}

type TLSConfig struct {
	CACert        string
	CAPath        string
	ClientCert    string
	ClientKey     string
	TLSServerName string
	Insecure      bool
}

type Client struct {
	Client *retryablehttp.Client

	Username, Password, BaseURL string

	common common.Service

	Security securityClient
}

type securityClient struct {
	Users        security.UserServiceInterface
	Roles        security.RoleServiceInterface
	Rolesmapping security.RolesmappingServiceInterface
	Actiongroups security.ActiongroupServiceInterface
	Tenants      security.TenantServiceInterface
	Health       security.HealthServiceInterface
}

func NewClient(config *ClientConfig) (*Client, error) {
	rc := retryablehttp.NewClient()

	if config.TLSConfig != nil {
		conf := &tls.Config{
			ServerName:         config.TLSConfig.TLSServerName,
			InsecureSkipVerify: config.TLSConfig.Insecure,
			MinVersion:         tls.VersionTLS12,
		}
		if config.TLSConfig.ClientCert != "" && config.TLSConfig.ClientKey != "" {
			clientCertificate, err := tls.LoadX509KeyPair(config.TLSConfig.ClientCert, config.TLSConfig.ClientKey)
			if err != nil {
				return nil, err
			}
			conf.Certificates = append(conf.Certificates, clientCertificate)
		}
		if config.TLSConfig.CACert != "" || config.TLSConfig.CAPath != "" {
			rootConfig := &rootcerts.Config{
				CAFile: config.TLSConfig.CACert,
				CAPath: config.TLSConfig.CAPath,
			}
			if err := rootcerts.ConfigureTLS(conf, rootConfig); err != nil {
				return nil, err
			}
		}

		rc.HTTPClient.Transport = &http.Transport{TLSClientConfig: conf}
	}

	c := &Client{
		Client:   rc,
		Username: config.Username,
		Password: config.Password,
		BaseURL:  config.BaseURL,
	}

	c.common.Client = c

	c.Security = securityClient{
		Users:        (*security.UserService)(&c.common),
		Roles:        (*security.RoleService)(&c.common),
		Rolesmapping: (*security.RolesmappingService)(&c.common),
		Actiongroups: (*security.ActiongroupService)(&c.common),
		Tenants:      (*security.TenantService)(&c.common),
		Health:       (*security.HealthService)(&c.common),
	}

	return c, nil
}

func (c *Client) GetBaseURL() string {
	return c.BaseURL
}

func (c *Client) Do(ctx context.Context, reqBytes interface{}, endpoint string, method string) ([]byte, error) {
	var req *http.Request
	var err error
	if reqBytes != nil {
		_reqBytes, err := json.Marshal(reqBytes)

		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, c.common.Client.GetBaseURL()+endpoint, bytes.NewReader(_reqBytes))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, c.common.Client.GetBaseURL()+endpoint, nil)
		if err != nil {
			return nil, err
		}
	}

	retryableReq, err := retryablehttp.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	retryableReq.Header.Add("Content-Type", "application/json")
	retryableReq.SetBasicAuth(c.Username, c.Password)

	resp, err := c.Client.Do(retryableReq.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusNotFound:

		var sr *common.StatusResponse

		err = json.Unmarshal(body, &sr)
		if err != nil {
			return nil, err
		}

		if sr.Status != nil && *sr.Status == string(common.Status.Error) {
			return nil, common.NewStatusError(*sr.Reason, *sr.InvalidKeys)
		}
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized: %s", resp.Body)
	}

	return body, nil
}

func (c *Client) Get(ctx context.Context, path string, T interface{}) error {
	body, err := c.Do(ctx, nil, path, http.MethodGet)
	if err != nil {
		return err
	}

	if body == nil {
		return nil
	}

	err = json.Unmarshal(body, &T)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Modify(ctx context.Context, path string, method string, reqBytes interface{}) error {
	body, err := c.Do(ctx, reqBytes, path, method)
	if err != nil {
		return err
	}

	var sr *common.StatusResponse

	err = json.Unmarshal(body, &sr)
	if err != nil {
		return err
	}

	if sr.Status != nil &&
		(*sr.Status == string(common.Status.Error) ||
			*sr.Status == string(common.Status.NotFound)) {

		if sr.InvalidKeys != nil {
			if sr.Reason != nil {
				return common.NewStatusError(*sr.Reason, *sr.InvalidKeys)
			}
			return common.NewStatusError(*sr.Message, *sr.InvalidKeys)
		}

		if sr.Reason != nil {
			return common.NewStatusError(*sr.Reason, nil)
		}

		return common.NewStatusError(*sr.Message, nil)
	}

	return nil
}
