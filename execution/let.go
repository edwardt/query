//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"github.com/couchbaselabs/query/errors"
	"github.com/couchbaselabs/query/plan"
	"github.com/couchbaselabs/query/value"
)

type Let struct {
	base
	plan *plan.Let
}

func NewLet(plan *plan.Let) *Let {
	rv := &Let{
		base: newBase(),
		plan: plan,
	}

	rv.output = rv
	return rv
}

func (this *Let) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLet(this)
}

func (this *Let) Copy() Operator {
	return &Let{this.base.copy(), this.plan}
}

func (this *Let) RunOnce(context *Context, parent value.Value) {
	this.runConsumer(this, context, parent)
}

func (this *Let) processItem(item value.AnnotatedValue, context *Context) bool {
	n := len(this.plan.Bindings())
	cv := value.NewScopeValue(make(map[string]interface{}, n), item)
	lv := value.NewAnnotatedValue(cv)
	lv.SetAttachments(item.Attachments())

	for _, b := range this.plan.Bindings() {
		v, e := b.Expression().Evaluate(item, context)
		if e != nil {
			context.Error(errors.NewError(e, "Error evaluating LET."))
			return false
		}

		lv.SetField(b.Variable(), v)
	}

	return this.sendItem(lv)
}
