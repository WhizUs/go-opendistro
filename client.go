package opendistro

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-rootcerts"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	usersEndpoint        = "/_opendistro/_security/api/internalusers/"
	rolesEndpoint        = "/_opendistro/_security/api/roles/"
	rolesMappingEndpoint = "/_opendistro/_security/api/rolesmapping/"
	actiongroupEndpoint  = "/_opendistro/_security/api/actiongroups/"
	tenantEndpoint       = "/_opendistro/_security/api/tenants/"
)

type ClientConfig struct {
	Username, Password, BaseURL string

	TLSConfig *TLSConfig
}

type TLSConfig struct {
	CACert string
	CAPath string
	ClientCert string
	ClientKey string
	TLSServerName string
	Insecure bool
}

type service struct {
	client *Client
}

type Client struct {
	client *retryablehttp.Client

	Username, Password, BaseURL string

	common service

	Users        *UserService
	Roles        RoleServiceInterface
	Rolesmapping *RolesmappingService
	Actiongroups *ActiongroupService
}

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
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
		client:   rc,
		Username: config.Username,
		Password: config.Password,
		BaseURL:  config.BaseURL,
	}

	c.common.client = c
	c.Users = (*UserService)(&c.common)
	c.Roles = (*RoleService)(&c.common)
	c.Rolesmapping = (*RolesmappingService)(&c.common)
	c.Actiongroups = (*ActiongroupService)(&c.common)

	return c, nil
}

func (c *Client) Do(ctx context.Context, reqBytes interface{}, endpoint string, method string) ([]byte, error) {
	var req *http.Request
	var err error
	if reqBytes != nil {
		_reqBytes, err := json.Marshal(reqBytes)

		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, c.common.client.BaseURL+endpoint, bytes.NewReader(_reqBytes))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, c.common.client.BaseURL+endpoint, nil)
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

	resp, err := c.client.Do(retryableReq.WithContext(ctx))
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

	return body, nil
}

func (c *Client) get(ctx context.Context, path string, T interface{}) error {
	body, err := c.Do(ctx, nil, path, http.MethodGet)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &T)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) modify(ctx context.Context, path string, method string, reqBytes interface{}) (*StatusResponse, error) {
	body, err := c.Do(ctx, reqBytes, path, method)
	if err != nil {
		return nil, err
	}

	var sr *StatusResponse

	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

