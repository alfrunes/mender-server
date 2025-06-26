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
	DefaultGroup     RatelimitParams        `json:"default"`
	RatelimitGroups  []RatelimitGroupParams `json:"groups"`
	MatchExpressions []MatchGroup           `json:"match"`
}

type RatelimitGroupParams struct {
	Name string `json:"name"`
	RatelimitParams
}

type RatelimitParams struct {
	// Quota is the number of requests that can be made within Interval
	Quota int64 `json:"quota"`

	// Interval is the time for the rate limit algorithm to reset.
	Interval time.Duration `json:"interval"`

	// EventExpression specifies a Go template for grouping events (requests)
	// when invoking the rate limiter. For example:
	// {{.Identity.Subject}}{{/* Group by JWT subject (user ID) */}}
	// {{.Identity.Tenant}}{{/* Group by tenant ID (shared quota) */}}
	EventExpression string `json:"event_expression"`
}

type MatchGroup struct {
	// APIPattern matches method and path of the incoming request using pattern
	// from Go standard library ServeMux.
	// https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	APIPattern string `json:"api_pattern"`

	// GroupExpression is a template string for selecting rate limit group.
	GroupExpression string `json:"group_expression,omitempty"`
}

// ratelimits:
//   api_groups:
//   - unfair
//   api_group_template: {{ .Identity.plan }}
//   parameters:
//    - quota: 100
//      interval: 60s
//      event_tamplate: {{ .Identity.Tenant }}:{{ .Identity.Subject }}
//      api_pattern: /
//    - tokens: 2
//      interval: 24h
//      event_template: {{ .Identity.Tenant }}
//      api_pattern: POST /api/management/v2/tenantadm/billing/subscriptions
//    - tokens: 1
//      interval: 168h
//      event_template: {{ "boohoo" }}
//      api_pattern: /

func LoadRatelimits(c config.Reader) (*RatelimitConfig, error) {
	if !c.GetBool(SettingRatelimitsEnable) {
		return nil, nil
	}
	ratelimitConfig := RatelimitConfig{
		DefaultGroup: RatelimitParams{
			Quota:           int64(c.GetInt(SettingRatelimitsDefaultQuota)),
			Interval:        c.GetDuration(SettingRatelimitsDefaultInterval),
			EventExpression: c.GetString(SettingRatelimitsDefaultEventExpression),
		},
	}
	err := config.UnmarshalSliceSetting(c,
		SettingRatelimitsGroups,
		&ratelimitConfig.RatelimitGroups,
	)
	if err != nil {
		return nil, fmt.Errorf("error loading rate limit groups: %w", err)
	}

	err = config.UnmarshalSliceSetting(c,
		SettingRatelimitsMatch,
		&ratelimitConfig.MatchExpressions,
	)
	if err != nil {
		return nil, fmt.Errorf("error loading rate limit match expressions: %w", err)
	}
	return &ratelimitConfig, nil
}
