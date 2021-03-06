//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import "encoding/json"

type Sequence struct {
	children []Operator `json:"children"`
}

func NewSequence(children ...Operator) *Sequence {
	return &Sequence{children}
}

func (this *Sequence) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSequence(this)
}

func (this *Sequence) Readonly() bool {
	for _, child := range this.children {
		if !child.Readonly() {
			return false
		}
	}

	return true
}

func (this *Sequence) New() Operator {
	return &Sequence{}
}

func (this *Sequence) Children() []Operator {
	return this.children
}

func (this *Sequence) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"#operator": "Sequence"}
	r["~children"] = this.children
	return json.Marshal(r)
}

func (this *Sequence) UnmarshalJSON(body []byte) error {
	var _unmarshalled struct {
		_        string            `json:"#operator"`
		Children []json.RawMessage `json:"~children"`
	}
	err := json.Unmarshal(body, &_unmarshalled)
	if err != nil {
		return err
	}

	this.children = []Operator{}

	for _, raw_child := range _unmarshalled.Children {
		var child_type struct {
			Op_name string `json:"#operator"`
		}
		var read_only struct {
			Readonly bool `json:"readonly"`
		}
		err = json.Unmarshal(raw_child, &child_type)
		if err != nil {
			return err
		}

		if child_type.Op_name == "" {
			err = json.Unmarshal(raw_child, &read_only)
			if err != nil {
				return err
			} else {
				// This should be a readonly object
			}
		} else {
			child_op, err := MakeOperator(child_type.Op_name, raw_child)
			if err != nil {
				return err
			}
			this.children = append(this.children, child_op)
		}
	}

	return err
}
