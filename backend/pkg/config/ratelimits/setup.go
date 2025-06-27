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

package ratelimits

import (
	"fmt"

	"github.com/mendersoftware/mender-server/pkg/config"
	"github.com/mendersoftware/mender-server/pkg/log"
	"github.com/mendersoftware/mender-server/pkg/rate"
	"github.com/mendersoftware/mender-server/pkg/redis"
)

func init() {
	config.Config.SetDefault(SettingRatelimitsDefaultInterval, "1m")
	config.Config.SetDefault(SettingRatelimitsDefaultQuota, "120")
	config.Config.SetDefault(SettingRatelimitsDefaultEventExpression,
		"{{with .Identity}}{{.Subject}}{{end}}")
}

const (
	SettingRatelimits                       = "ratelimits"
	SettingRatelimitsEnable                 = SettingRatelimits + ".enable"
	SettingRatelimitsDefault                = SettingRatelimits + ".default"
	SettingRatelimitsDefaultQuota           = SettingRatelimitsDefault + ".quota"
	SettingRatelimitsDefaultInterval        = SettingRatelimitsDefault + ".interval"
	SettingRatelimitsDefaultEventExpression = SettingRatelimitsDefault + ".event_expression"
	SettingRatelimitsGroups                 = SettingRatelimits + ".groups"
	SettingRatelimitsMatch                  = SettingRatelimits + ".match"
)

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

func SetupRedisRateLimits(
	redisClient redis.Client,
	keyPrefix string,
	c config.Reader,
) (*rate.HTTPLimiter, error) {
	if !c.GetBool(SettingRatelimitsEnable) {
		return nil, nil
	}
	lims, err := LoadRatelimits(c)
	if err != nil {
		return nil, err
	}
	log.NewEmpty().Debugf("loaded rate limit configuration: %v", lims)
	defaultPrefix := fmt.Sprintf("%s:rate:default", keyPrefix)
	defaultLimiter := redis.NewFixedWindowRateLimiter(
		redisClient, defaultPrefix,
		lims.DefaultGroup.Interval, lims.DefaultGroup.Quota,
	)
	mux, err := rate.NewHTTPLimiter(defaultLimiter, c.GetString(SettingRatelimitsDefaultEventExpression))
	if err != nil {
		return nil, fmt.Errorf("error setting up rate limits: %w", err)
	}
	for _, group := range lims.RatelimitGroups {
		groupPrefix := fmt.Sprintf("%s:rate:g:%s", keyPrefix, group.Name)
		limiter := redis.NewFixedWindowRateLimiter(redisClient, groupPrefix, group.Interval, group.Quota)
		err = mux.AddRateLimitGroup(limiter, group.Name, group.EventExpression)
		if err != nil {
			return nil, fmt.Errorf("error setting up rate limit group %s: %w", group.Name, err)
		}
	}
	for _, expr := range lims.MatchExpressions {
		mux.MatchHTTPPattern(expr.APIPattern, expr.GroupExpression)
	}
	return mux, nil
}
