/*
 * Copyright 2019 Jobteaser <opensource@jobteaser.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package circleci expose functions to manage environment variables in
// a CircleCI project.
package circleci // import "github.com/jobteaser-oss/circleci-env"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

const httpClientTimeout = 30 * time.Second

// Client contains anm `http.Client` and a URL.
type Client struct {
	http *http.Client
	url  *url.URL
}

// NewClient returns a pointer of `Client` struct with timeout and
// authentication properly configured.
func NewClient(token string) (*Client, error) {
	dialer := net.Dialer{Timeout: httpClientTimeout}

	transport := http.Transport{
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: httpClientTimeout,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   httpClientTimeout,
	}

	uri, err := url.Parse("https://circleci.com/api/v1.1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the CircleCI base URL")
	}

	params := url.Values{}
	params.Add("circle-token", token)
	uri.RawQuery = params.Encode()

	return &Client{
		http: &client,
		url:  uri,
	}, nil
}

// Env represent the CircleCI API response.
type Env struct {
	Key   string `json:"name"`
	Value string `json:"value"`
}

// ListEnv returns four 'x' characters plus the last four ASCII characters
// of the value, consistent with the display of environment variable values
// in the CircleCI website.
func (client *Client) ListEnv(vcsType, username, project string) ([]*Env, error) {
	uri := *client.url
	uri.Path = fmt.Sprintf("%s/project/%s/%s/%s/envvar",
		uri.Path,
		vcsType,
		username,
		project,
	)

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build the HTTP request")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request to CircleCI API")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	default:
		return nil, errors.New("unkown error failed to list keys")
	}

	var envs []*Env
	json.NewDecoder(resp.Body).Decode(&envs)
	return envs, nil
}

// GetEnv returns the hidden value of environment variable.
func (client *Client) GetEnv(vcsType, username, project, key string) (*Env, error) {
	uri := *client.url
	uri.Path = fmt.Sprintf("%s/project/%s/%s/%s/envvar/%s",
		uri.Path,
		vcsType,
		username,
		project,
		key,
	)

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build the HTTP request")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request to CircleCI API")

	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 404:
		return nil, fmt.Errorf("the key %q does not exist", key)
	default:
		return nil, fmt.Errorf("unkown error failed to get %q key", key)
	}

	var env Env
	json.NewDecoder(resp.Body).Decode(&env)
	return &env, nil
}

// SetEnv creates or updates a new environment variable.
func (client *Client) SetEnv(vcsType, username, project, key, value string) error {
	uri := *client.url
	uri.Path = fmt.Sprintf("%s/project/%s/%s/%s/envvar",
		uri.Path,
		vcsType,
		username,
		project,
	)

	env := Env{Key: key, Value: value}
	body, err := json.Marshal(env)
	if err != nil {
		return errors.Wrap(err, "failed to marshal the JSON payload")
	}

	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to build the HTTP request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute HTTP request to CircleCI API")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		break
	case 404:
		return fmt.Errorf("the key %q does not exist", key)
	default:
		return fmt.Errorf("unkown error failed to create %q key", key)
	}

	return nil
}

// DeleteEnv deletes the environment variable.
func (client *Client) DeleteEnv(vcsType, username, project, key string) error {
	uri := *client.url
	uri.Path = fmt.Sprintf("%s/project/%s/%s/%s/envvar/%s",
		uri.Path,
		vcsType,
		username,
		project,
		key,
	)

	req, err := http.NewRequest("DELETE", uri.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to build the HTTP request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute HTTP request to CircleCI API")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 404:
		return fmt.Errorf("the key %q does not exist", key)
	default:
		return fmt.Errorf("unkown error failed to delete %q key", key)
	}

	return nil
}
