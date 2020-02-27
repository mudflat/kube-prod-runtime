/*
 * Bitnami Kubernetes Production Runtime - A collection of services that makes it
 * easy to run production workloads in Kubernetes.
 *
 * Copyright 2018-2019 Bitnami
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

package generic

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"

	"github.com/bitnami/kube-prod-runtime/kubeprod/tools"
)

func (conf *GenericConfig) Generate(ctx context.Context) error {
	flags := conf.flags

	// Leaks secrets to log!
	//log.Debugf("Input config: %#v", conf)

	if conf.ContactEmail == "" {
		email, err := flags.GetString(flagEmail)
		if err != nil {
			return err
		}
		conf.ContactEmail = email
	}

	if conf.DnsZone == "" {
		domain, err := flags.GetString(flagDNSSuffix)
		if err != nil {
			return err
		}
		conf.DnsZone = domain
	}

	// keycloak setup
	if conf.Keycloak.Password == "" {
		password, err := flags.GetString(flagKeycloakPassword)
		if err != nil {
			return err
		}
		conf.Keycloak.Password = password
	}

	if conf.Keycloak.ClientID == "" {
		conf.Keycloak.ClientID = "bkpr"
	}

	if conf.Keycloak.ClientSecret == "" {
		conf.Keycloak.ClientSecret = uuid.New().String()
	}

	//
	// powerdns setup
	//
	if conf.PowerDNS.ApiKey == "" {
		conf.PowerDNS.ApiKey = uuid.New().String()
	}

	//
	// oauth2-proxy setup
	//
	log.Debug("Starting oauth2-proxy setup")

	if conf.OauthProxy.CookieSecret == "" {
		// I Quote: cookie_secret must be 16, 24, or 32 bytes
		// to create an AES cipher when pass_access_token ==
		// true or cookie_refresh != 0
		secret, err := tools.Base64RandBytes(24)
		if err != nil {
			return err
		}
		conf.OauthProxy.CookieSecret = secret
	}

	if conf.OauthProxy.AuthzDomain == "" {
		domain, err := flags.GetString(flagAuthzDomain)
		if err != nil {
			return err
		}
		conf.OauthProxy.AuthzDomain = domain
	}

	return nil
}
