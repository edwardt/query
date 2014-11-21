//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"encoding/json"

	"github.com/couchbaselabs/query/expression"
)

type Offset struct {
	readonly
	expr expression.Expression
}

func NewOffset(expr expression.Expression) *Offset {
	return &Offset{
		expr: expr,
	}
}

func (this *Offset) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitOffset(this)
}

func (this *Offset) Expression() expression.Expression {
	return this.expr
}

func (this *Offset) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"Type": "offset"}
	exprStr := expression.NewStringer().Visit(this.expr)
	r["Expr"] = exprStr
	return json.Marshal(r)
}

type Limit struct {
	readonly
	expr expression.Expression
}

func NewLimit(expr expression.Expression) *Limit {
	return &Limit{
		expr: expr,
	}
}

func (this *Limit) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLimit(this)
}

func (this *Limit) Expression() expression.Expression {
	return this.expr
}

func (this *Limit) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"Type": "limit"}
	exprStr := expression.NewStringer().Visit(this.expr)
	r["Expr"] = exprStr
	return json.Marshal(r)
}
