// Copyright 2023 Northern.tech AS
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
package store

import (
	"fmt"

	"github.com/mendersoftware/mender-server/services/inventory/model"
)

type ComparisonOperator int

func (op ComparisonOperator) String() string {
	switch op {
	case Eq:
		return "$eq"
	}
	panic(fmt.Sprintf("invalid comparison operator %s", op))
}

const (
	Eq ComparisonOperator = 1 << iota
)

type Filter struct {
	AttrName  string
	AttrScope string
	Value     any
	Operator  ComparisonOperator
}

type Filters []Filter

func (filters Filters) ToMongoFilter() map[string]any {
	fltrs := make([]map[string]any, 0, len(filters))
	for _, elem := range filters {
		scope := elem.AttrScope
		switch scope {
		case model.AttrScopeSystem:
			fltrs = append(fltrs, map[string]any{elem.AttrName: map[string]any{
				elem.Operator.String(): elem.Value,
			}})
		case model.AttrScopeIdentity:
			if elem.AttrName == "status" {
				fltrs = append(fltrs, map[string]any{elem.AttrName: map[string]any{
					elem.Operator.String(): elem.Value,
				}})
			} else {
				fltrs = append(fltrs, map[string]any{
					string(scope): map[string]any{"$elemMatch": map[string]any{
						"name": elem.AttrName,
						"value": map[string]any{
							elem.Operator.String(): elem.Value,
						},
					}},
				})
			}

		case "":
			scope = model.AttrScopeInventory
			fallthrough
		default:
			fltrs = append(fltrs, map[string]any{
				string(scope): map[string]any{"$elemMatch": map[string]any{
					"name": elem.AttrName,
					"value": map[string]any{
						elem.Operator.String(): elem.Value,
					},
				}},
			})
		}
	}
	return map[string]any{"$and": fltrs}
}

type Sort struct {
	AttrName  string
	AttrScope string
	Ascending bool
}

type ListQuery struct {
	Skip      int
	Limit     int
	Filters   Filters
	Sort      *Sort
	HasGroup  *bool
	GroupName string
}
