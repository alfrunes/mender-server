// Copyright 2022 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/mendersoftware/mender-server/services/deviceauth/model"
)

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// AcceptDeviceAuth provides a mock function with given fields: ctx, dev_id, auth_id
func (_m *App) AcceptDeviceAuth(ctx context.Context, dev_id string, auth_id string) error {
	ret := _m.Called(ctx, dev_id, auth_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, dev_id, auth_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecommissionDevice provides a mock function with given fields: ctx, dev_id
func (_m *App) DecommissionDevice(ctx context.Context, dev_id string) error {
	ret := _m.Called(ctx, dev_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, dev_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAuthSet provides a mock function with given fields: ctx, dev_id, auth_id
func (_m *App) DeleteAuthSet(ctx context.Context, dev_id string, auth_id string) error {
	ret := _m.Called(ctx, dev_id, auth_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, dev_id, auth_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDevice provides a mock function with given fields: ctx, dev_id
func (_m *App) DeleteDevice(ctx context.Context, dev_id string) error {
	ret := _m.Called(ctx, dev_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, dev_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTenantLimit provides a mock function with given fields: ctx, tenant_id, limit
func (_m *App) DeleteTenantLimit(ctx context.Context, tenant_id string, limit string) error {
	ret := _m.Called(ctx, tenant_id, limit)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenant_id, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTokens provides a mock function with given fields: ctx, tenantID, deviceID
func (_m *App) DeleteTokens(ctx context.Context, tenantID string, deviceID string) error {
	ret := _m.Called(ctx, tenantID, deviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, deviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDevCountByStatus provides a mock function with given fields: ctx, status
func (_m *App) GetDevCountByStatus(ctx context.Context, status string) (int, error) {
	ret := _m.Called(ctx, status)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, status)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevice provides a mock function with given fields: ctx, dev_id
func (_m *App) GetDevice(ctx context.Context, dev_id string) (*model.Device, error) {
	ret := _m.Called(ctx, dev_id)

	var r0 *model.Device
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Device); ok {
		r0 = rf(ctx, dev_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Device)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, dev_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevices provides a mock function with given fields: ctx, skip, limit, filter
func (_m *App) GetDevices(ctx context.Context, skip uint, limit uint, filter model.DeviceFilter) ([]model.Device, error) {
	ret := _m.Called(ctx, skip, limit, filter)

	var r0 []model.Device
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint, model.DeviceFilter) []model.Device); ok {
		r0 = rf(ctx, skip, limit, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Device)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, uint, model.DeviceFilter) error); ok {
		r1 = rf(ctx, skip, limit, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLimit provides a mock function with given fields: ctx, name
func (_m *App) GetLimit(ctx context.Context, name string) (*model.Limit, error) {
	ret := _m.Called(ctx, name)

	var r0 *model.Limit
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Limit); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Limit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantDeviceStatus provides a mock function with given fields: ctx, tenantId, deviceId
func (_m *App) GetTenantDeviceStatus(ctx context.Context, tenantId string, deviceId string) (*model.Status, error) {
	ret := _m.Called(ctx, tenantId, deviceId)

	var r0 *model.Status
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Status); ok {
		r0 = rf(ctx, tenantId, deviceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Status)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantId, deviceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantLimit provides a mock function with given fields: ctx, name, tenant_id
func (_m *App) GetTenantLimit(ctx context.Context, name string, tenant_id string) (*model.Limit, error) {
	ret := _m.Called(ctx, name, tenant_id)

	var r0 *model.Limit
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Limit); ok {
		r0 = rf(ctx, name, tenant_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Limit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, tenant_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthCheck provides a mock function with given fields: ctx
func (_m *App) HealthCheck(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PreauthorizeDevice provides a mock function with given fields: ctx, req
func (_m *App) PreauthorizeDevice(ctx context.Context, req *model.PreAuthReq) (*model.Device, error) {
	ret := _m.Called(ctx, req)

	var r0 *model.Device
	if rf, ok := ret.Get(0).(func(context.Context, *model.PreAuthReq) *model.Device); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Device)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.PreAuthReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProvisionTenant provides a mock function with given fields: ctx, tenant_id
func (_m *App) ProvisionTenant(ctx context.Context, tenant_id string) error {
	ret := _m.Called(ctx, tenant_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, tenant_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RejectDeviceAuth provides a mock function with given fields: ctx, dev_id, auth_id
func (_m *App) RejectDeviceAuth(ctx context.Context, dev_id string, auth_id string) error {
	ret := _m.Called(ctx, dev_id, auth_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, dev_id, auth_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResetDeviceAuth provides a mock function with given fields: ctx, dev_id, auth_id
func (_m *App) ResetDeviceAuth(ctx context.Context, dev_id string, auth_id string) error {
	ret := _m.Called(ctx, dev_id, auth_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, dev_id, auth_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RevokeToken provides a mock function with given fields: ctx, tokenID
func (_m *App) RevokeToken(ctx context.Context, tokenID string) error {
	ret := _m.Called(ctx, tokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, tokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetTenantLimit provides a mock function with given fields: ctx, tenant_id, limit
func (_m *App) SetTenantLimit(ctx context.Context, tenant_id string, limit model.Limit) error {
	ret := _m.Called(ctx, tenant_id, limit)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.Limit) error); ok {
		r0 = rf(ctx, tenant_id, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitAuthRequest provides a mock function with given fields: ctx, r
func (_m *App) SubmitAuthRequest(ctx context.Context, r *model.AuthReq) (string, error) {
	ret := _m.Called(ctx, r)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *model.AuthReq) string); ok {
		r0 = rf(ctx, r)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.AuthReq) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyToken provides a mock function with given fields: ctx, token
func (_m *App) VerifyToken(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
