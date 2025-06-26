// Copyright 2025 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package config

import (
	"fmt"
	"time"

	"github.com/mendersoftware/mender-server/pkg/config"
)

type RatelimitConfig struct {
	ParamGroups map[string]RatelimitParam `json:"parameter_groups"`
	Rules       []Rules                   `json:"rules"`
}

type RatelimitParam struct {
	Tokens   int           `json:"tokens"`
	Interval time.Duration `json:"interval"`
}

type Rules struct {
	// APIPattern matches method and path of the incoming request using pattern
	// from Go standard library ServeMux.
	// https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	APIPattern string `mapstructure:"pattern" json:"pattern"`
	ParamGroup string `mapstructure:"group" json:"group"`
}

// ratelimits:
//   parameter_groups:
//     default:
//     	tokens: 100
//     	interval: 60s
//     	group_by: {{ .Identity.Tenant }}:{{ .Identity.Subject }}
//     billing_upgrade:
//     	tokens: 2
//     	interval: 24h
//     	group_by: {{ .Identity.Tenant }}
//		rules:
//			- pattern: /
//				group: default
//			- pattern: POST /api/management/v2/tenantadm/billing/subscriptions
//				group: billing_upgrade

func LoadRatelimits(c config.Reader) (*RatelimitConfig, error) {
	var ratelimitConfig RatelimitConfig
	ratelimitConfig.Default.Interval = c.GetDuration(SettingRatelimitsDefaultInterval)
	ratelimitConfig.Default.Tokens = c.GetInt(SettingRatelimitsDefaultTokens)
	err := config.UnmarshalSliceSetting(c,
		SettingRatelimitsOverrides,
		&ratelimitConfig.Override,
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing rate limit overrides: %w", err)
	}
	return &ratelimitConfig, nil
}
