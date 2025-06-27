// Copyright 2025 Northern.tech AS
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

package rest

import (
	"net/http"
	"net/url"
)

// RewriteForwardedRequest makes a shallow clone of request and replaces
// the URL and Method with X-Forwarded-* headers.
func RewriteForwardedRequest(request *http.Request) *http.Request {
	newRequest := new(http.Request)
	*newRequest = *request
	newRequest.URL = new(url.URL)
	*newRequest.URL = *request.URL
	newRequest.URL.Host = request.Header.Get("X-Forwarded-Host")
	newRequest.URL.Path = request.Header.Get("X-Forwarded-Path")
	newRequest.URL.RawPath = ""
	newRequest.Method = request.Header.Get("X-Forwarded-Method")
	return newRequest
}
