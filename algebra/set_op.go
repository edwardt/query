//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

import (
	"github.com/couchbaselabs/query/expression"
	"github.com/couchbaselabs/query/value"
)

/*
This is the base class for the set operations. Type setOp
is a struct that contains two fields that represent
results from multiple subselect statements.
*/
type setOp struct {
	first  Subresult `json:"first"`
	second Subresult `json:"second"`
}

/*
Returns the shape of the first subresult.
*/
func (this *setOp) Signature() value.Value {
	return this.first.Signature()
}

/*
Returns true if either of the subresults are correlated.
*/
func (this *setOp) IsCorrelated() bool {
	return this.first.IsCorrelated() || this.second.IsCorrelated()
}

/*
Applies mapper to the two subresults.
*/
func (this *setOp) MapExpressions(mapper expression.Mapper) (err error) {
	err = this.first.MapExpressions(mapper)
	if err != nil {
		return
	}

	return this.second.MapExpressions(mapper)
}

/*
   Returns all contained Expressions.
*/
func (this *setOp) Expressions() expression.Expressions {
	return append(this.first.Expressions(), this.second.Expressions()...)
}

/*
Fully qualifies the identifiers in the first and second sub-result
using the input parent.
*/
func (this *setOp) Formalize(parent *expression.Formalizer) (f *expression.Formalizer, err error) {
	_, err = this.first.Formalize(parent)
	if err != nil {
		return nil, err
	}

	_, err = this.second.Formalize(parent)
	if err != nil {
		return nil, err
	}

	return expression.NewFormalizer(), nil
}

/*
Returns the first sub result.
*/
func (this *setOp) First() Subresult {
	return this.first
}

/*
Returns the second sub result.
*/
func (this *setOp) Second() Subresult {
	return this.second
}

/*
Represents the result of the UNION set operation. Type
unionSubresult inherits from setOp.
*/
type unionSubresult struct {
	setOp
}

/*
Returns the shape of the union result. If the two sub results
are equal return the first value. If either of the inputs
to the union setop are not objects then return the _JSON_SIGNATURE.
Range through the two objects and check for equality and return
object value.
*/
func (this *unionSubresult) Signature() value.Value {
	first := this.first.Signature()
	second := this.second.Signature()

	if first.Equals(second) {
		return first
	}

	if first.Type() != value.OBJECT ||
		second.Type() != value.OBJECT {
		return _JSON_SIGNATURE
	}

	rv := first.Copy()
	sa := second.Actual().(map[string]interface{})
	for k, v := range sa {
		cv, ok := rv.Field(k)
		if ok {
			if !value.NewValue(cv).Equals(value.NewValue(v)) {
				rv.SetField(k, _JSON_SIGNATURE)
			}
		} else {
			rv.SetField(k, v)
		}
	}

	return rv
}

/*
New JSON string value.
*/
var _JSON_SIGNATURE = value.NewValue(value.JSON.String())

/*
Represents the Union set operation used to combine results
from multiple subselects. UNION returns values from the
first sub-select and from the second sub-select. Type Union
is a struct that inherits from the unionSubresult.
*/
type Union struct {
	unionSubresult
}

/*
The function NewUnion returns a pointer to the
Union struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewUnion(first, second Subresult) Subresult {
	return &Union{
		unionSubresult{
			setOp{
				first:  first,
				second: second,
			},
		},
	}
}

/*
It calls the VisitUnion method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *Union) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitUnion(this)
}

/*
Represents the Union all set operation used to combine results
from multiple subselects. It returns all applicable values,
including duplicates. These forms are faster, because they do
not compute distinct results. Type UnionAll is a struct that
inherits from the unionSubresult.
*/
type UnionAll struct {
	unionSubresult
}

/*
The function NewUnionAll returns a pointer to the
UnionAll struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewUnionAll(first, second Subresult) Subresult {
	return &UnionAll{
		unionSubresult{
			setOp{
				first:  first,
				second: second,
			},
		},
	}
}

/*
It calls the VisitUnionAll method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *UnionAll) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitUnionAll(this)
}

/*
Represents the Intersect set operation used to combine results
from multiple subselects. INTERSECT returns values from the
first sub-select that are present in the second sub-select.
Type Intersect is a struct that inherits from setOp.
*/
type Intersect struct {
	setOp
}

/*
The function NewIntersect returns a pointer to the
Intersect struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewIntersect(first, second Subresult) Subresult {
	return &Intersect{
		setOp{
			first:  first,
			second: second,
		},
	}
}

/*
It calls the VisitIntersect method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *Intersect) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitIntersect(this)
}

/*
Represents the Intersect All set operation used to combine results
from multiple subselects. It returns all applicable values,
including duplicates. These forms are faster, because they do
not compute distinct results. Type IntersectAll is a struct that
inherits from setOp.
*/
type IntersectAll struct {
	setOp
}

/*
The function NewIntersectAll returns a pointer to the
IntersectAll struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewIntersectAll(first, second Subresult) Subresult {
	return &IntersectAll{
		setOp{
			first:  first,
			second: second,
		},
	}
}

/*
It calls the VisitIntersectAll method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *IntersectAll) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitIntersectAll(this)
}

/*
Represents the Except set operation used to combine results
from multiple subselects. EXCEPT returns values from the
first sub-select that are absent from the second sub-select.
Type Except is a struct that inherits from setOp.
*/
type Except struct {
	setOp
}

/*
The function NewExcept returns a pointer to the
Except struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewExcept(first, second Subresult) Subresult {
	return &Except{
		setOp{
			first:  first,
			second: second,
		},
	}
}

/*
It calls the VisitExcept method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *Except) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitExcept(this)
}

/*
Represents the Except All set operation used to combine results
from multiple subselects. It returns all applicable values,
including duplicates. These forms are faster, because they do
not compute distinct results. Type ExceptAll is a struct that
inherits from setOp.
*/
type ExceptAll struct {
	setOp
}

/*
The function NewExceptAll returns a pointer to the
ExceptAll struct with the input argument used to set
the first and second fields in the subresult.
*/
func NewExceptAll(first, second Subresult) Subresult {
	return &ExceptAll{
		setOp{
			first:  first,
			second: second,
		},
	}
}

/*
It calls the VisitExceptAll method by passing
in the receiver and returns the interface. It is a
visitor pattern.
*/
func (this *ExceptAll) Accept(visitor NodeVisitor) (interface{}, error) {
	return visitor.VisitExceptAll(this)
}
