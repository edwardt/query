//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

type ResultTerms []*ResultTerm

type ResultTerm struct {
	star bool       `json:"star"`
	expr Expression `json:"expr"`
	as   string     `json:"as"`
}

func (this *ResultTerm) Star() bool {
	return this.star
}

func (this *ResultTerm) Expression() Expression {
	return this.expr
}

func (this *ResultTerm) As() string {
	return this.as
}

func (this *ResultTerm) Alias() string {
	if this.star {
		return ""
	} else if this.as != "" {
		return this.as
	} else {
		return this.expr.Alias()
	}
}