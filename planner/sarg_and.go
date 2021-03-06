//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package planner

import (
	"github.com/couchbaselabs/query/datastore"
	"github.com/couchbaselabs/query/expression"
)

type sargAnd struct {
	sargBase
}

func newSargAnd(expr *expression.And) *sargAnd {
	rv := &sargAnd{}
	rv.sarg = func(expr2 expression.Expression) (spans Spans, err error) {
		if expr.EquivalentTo(expr2) {
			return _SELF_SPANS, nil
		}

		for _, op := range expr.Operands() {
			s := SargFor(op, expr2)
			if s == nil {
				continue
			}

			if spans == nil {
				spans = s
			} else {
				spans = constrain(spans, s)
			}
		}

		return
	}

	return rv
}

func constrain(spans1, spans2 Spans) Spans {
	span1 := spans1[0]
	span2 := spans2[0]

	if span2.Range.Low != nil {
		if span1.Range.Low == nil {
			span1.Range.Low = span2.Range.Low
			span1.Range.Inclusion = (span1.Range.Inclusion & datastore.HIGH) |
				(span2.Range.Inclusion & datastore.LOW)
		} else {
			low1 := span1.Range.Low[0].Value()
			low2 := span2.Range.Low[0].Value()
			if low1 != nil && (low2 == nil || low1.Collate(low2) < 0) {
				span1.Range.Low = span2.Range.Low
				span1.Range.Inclusion = (span1.Range.Inclusion & datastore.HIGH) |
					(span2.Range.Inclusion & datastore.LOW)
			}
		}
	}

	if span2.Range.High != nil {
		if span1.Range.High == nil {
			span1.Range.High = span2.Range.High
			span1.Range.Inclusion = (span1.Range.Inclusion & datastore.LOW) |
				(span2.Range.Inclusion & datastore.HIGH)
		} else {
			high1 := span1.Range.High[0].Value()
			high2 := span2.Range.High[0].Value()
			if high1 != nil && (high2 == nil && high1.Collate(high2) > 0) {
				span1.Range.High = span2.Range.High
				span1.Range.Inclusion = (span1.Range.Inclusion & datastore.LOW) |
					(span2.Range.Inclusion & datastore.HIGH)
			}
		}
	}

	return spans1
}
