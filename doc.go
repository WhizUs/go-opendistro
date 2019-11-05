// Copyright 2019 WhizUs GmbH. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package go-opendistro provides a simple client implementation to interact with the OpenDistro for Elasticsearch REST APIs.

Usage:

	import (
        "github.com/WhizUs/go-opendistro"
        "github.com/WhizUs/go-opendistro/security"
    )

A client can be instantiated by providing a client configuration containing the user, password, base URL and a TLS configuration (optionally). For example:

	clientConfig := &opendistro.ClientConfig{
		Username:  "vault",
		Password:  "vault",
		BaseURL:   "https://es.dev.whizus.net",
		TLSConfig: nil,
	}

	client, err := opendistro.NewClient(clientConfig)
	if err != nil {
		fmt.Printf("instantiate client: %s\n", err)
	}

	if err := client.Security.Users.Create(context.TODO(), "kirk", &security.UserCreate{
		Password: "kirkpass",
		BackendRoles: []string{
			"captains",
			"starfleet",
		},
		Attributes: map[string]string{
			"attribute1": "value1",
			"attribute2": "value2",
		},
	}); err != nil {
		fmt.Printf("create user: %s\n", err)
	}

Some code snippets are provided within the https://github.com/WhizUs/go-opendistro/tree/master/example directory.

Each of the resources is aimed to be implemented by a Go service object (f.e. opendistro.Security.UserService) which in turn
provides available methods of the resource.
*/
package opendistro
