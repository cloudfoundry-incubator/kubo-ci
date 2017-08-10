/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"context"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
	"github.com/cloudfoundry/bosh-gcscli/config"
)

const uaString = "gcscli"

// NewSDK returns context and client necessary to instantiate a client
// based off of the provided configuration.
func NewSDK(c config.GCSCli) (context.Context, *storage.Client, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error
	ua := option.WithUserAgent(uaString)
	if c.CredentialsSource == "" {
		client, err = storage.NewClient(ctx, ua)
	} else {
		client, err = storage.NewClient(ctx, ua,
			option.WithServiceAccountFile(c.CredentialsSource))
	}
	return ctx, client, err
}
