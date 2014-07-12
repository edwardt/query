//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

/*

Package file provides a file-based implementation of the catalog
package.

*/
package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/couchbaselabs/query/catalog"
	"github.com/couchbaselabs/query/errors"
	"github.com/couchbaselabs/query/expression"
	"github.com/couchbaselabs/query/value"
)

// site is the root for the file-based Site.
type site struct {
	path           string
	namespaces     map[string]*namespace
	namespaceNames []string
}

func (s *site) Id() string {
	return s.path
}

func (s *site) URL() string {
	return "file://" + s.path
}

func (s *site) NamespaceIds() ([]string, errors.Error) {
	return s.NamespaceNames()
}

func (s *site) NamespaceNames() ([]string, errors.Error) {
	return s.namespaceNames, nil
}

func (s *site) NamespaceById(id string) (p catalog.Namespace, e errors.Error) {
	return s.NamespaceByName(id)
}

func (s *site) NamespaceByName(name string) (p catalog.Namespace, e errors.Error) {
	p, ok := s.namespaces[strings.ToUpper(name)]
	if !ok {
		e = errors.NewError(nil, "Namespace "+name+" not found.")
	}

	return
}

// NewSite creates a new file-based site for the given filepath.
func NewSite(path string) (s catalog.Site, e errors.Error) {
	path, er := filepath.Abs(path)
	if er != nil {
		return nil, errors.NewError(er, "")
	}

	fs := &site{path: path}

	e = fs.loadNamespaces()
	if e != nil {
		return
	}

	s = fs
	return
}

func (s *site) loadNamespaces() (e errors.Error) {
	dirEntries, er := ioutil.ReadDir(s.path)
	if er != nil {
		return errors.NewError(er, "")
	}

	s.namespaces = make(map[string]*namespace, len(dirEntries))
	s.namespaceNames = make([]string, 0, len(dirEntries))

	var p *namespace
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			s.namespaceNames = append(s.namespaceNames, dirEntry.Name())
			diru := strings.ToUpper(dirEntry.Name())
			if _, ok := s.namespaces[diru]; ok {
				return errors.NewError(nil, "Duplicate namespace name "+dirEntry.Name())
			}

			p, e = newNamespace(s, dirEntry.Name())
			if e != nil {
				return
			}

			s.namespaces[diru] = p
		}
	}

	return
}

// namespace represents a file-based Namespace.
type namespace struct {
	site          *site
	name          string
	keyspaces     map[string]*keyspace
	keyspaceNames []string
}

func (p *namespace) SiteId() string {
	return p.site.Id()
}

func (p *namespace) Id() string {
	return p.Name()
}

func (p *namespace) Name() string {
	return p.name
}

func (p *namespace) KeyspaceIds() ([]string, errors.Error) {
	return p.KeyspaceNames()
}

func (p *namespace) KeyspaceNames() ([]string, errors.Error) {
	return p.keyspaceNames, nil
}

func (p *namespace) KeyspaceById(id string) (b catalog.Keyspace, e errors.Error) {
	return p.KeyspaceByName(id)
}

func (p *namespace) KeyspaceByName(name string) (b catalog.Keyspace, e errors.Error) {
	b, ok := p.keyspaces[strings.ToUpper(name)]
	if !ok {
		e = errors.NewError(nil, "Keyspace "+name+" not found.")
	}

	return
}

func (p *namespace) path() string {
	return filepath.Join(p.site.path, p.name)
}

// newNamespace creates a new namespace.
func newNamespace(s *site, dir string) (p *namespace, e errors.Error) {
	p = new(namespace)
	p.site = s
	p.name = dir

	e = p.loadKeyspaces()
	return
}

func (p *namespace) loadKeyspaces() (e errors.Error) {
	dirEntries, er := ioutil.ReadDir(p.path())
	if er != nil {
		return errors.NewError(er, "")
	}

	p.keyspaces = make(map[string]*keyspace, len(dirEntries))
	p.keyspaceNames = make([]string, 0, len(dirEntries))

	var b *keyspace
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			diru := strings.ToUpper(dirEntry.Name())
			if _, ok := p.keyspaces[diru]; ok {
				return errors.NewError(nil, "Duplicate keyspace name "+dirEntry.Name())
			}

			b, e = newKeyspace(p, dirEntry.Name())
			if e != nil {
				return
			}

			p.keyspaces[diru] = b
			p.keyspaceNames = append(p.keyspaceNames, b.Name())
		}
	}

	return
}

// keyspace is a file-based keyspace.
type keyspace struct {
	namespace *namespace
	name      string
	indexes   map[string]catalog.Index
	primary   catalog.PrimaryIndex
}

func (b *keyspace) NamespaceId() string {
	return b.namespace.Id()
}

func (b *keyspace) Id() string {
	return b.Name()
}

func (b *keyspace) Name() string {
	return b.name
}

func (b *keyspace) Count() (int64, errors.Error) {
	dirEntries, er := ioutil.ReadDir(b.path())
	if er != nil {
		return 0, errors.NewError(er, "")
	}
	return int64(len(dirEntries)), nil
}

func (b *keyspace) IndexIds() ([]string, errors.Error) {
	rv := make([]string, 0, len(b.indexes))
	for name, _ := range b.indexes {
		rv = append(rv, name)
	}
	return rv, nil
}

func (b *keyspace) IndexNames() ([]string, errors.Error) {
	rv := make([]string, 0, len(b.indexes))
	for name, _ := range b.indexes {
		rv = append(rv, name)
	}
	return rv, nil
}

func (b *keyspace) IndexById(id string) (catalog.Index, errors.Error) {
	return b.IndexByName(id)
}

func (b *keyspace) IndexByName(name string) (catalog.Index, errors.Error) {
	index, ok := b.indexes[name]
	if !ok {
		return nil, errors.NewError(nil, fmt.Sprintf("Index %v not found.", name))
	}
	return index, nil
}

func (b *keyspace) IndexByPrimary() (catalog.PrimaryIndex, errors.Error) {
	return b.primary, nil
}

func (b *keyspace) Indexes() ([]catalog.Index, errors.Error) {
	rv := make([]catalog.Index, 0, len(b.indexes))
	for _, index := range b.indexes {
		rv = append(rv, index)
	}
	return rv, nil
}

func (b *keyspace) CreatePrimaryIndex() (catalog.PrimaryIndex, errors.Error) {
	if b.primary != nil {
		return b.primary, nil
	}

	return nil, errors.NewError(nil, "Not supported.")
}

func (b *keyspace) CreateIndex(name string, equalKey, rangeKey expression.Expressions, using catalog.IndexType) (catalog.Index, errors.Error) {
	return nil, errors.NewError(nil, "Not supported.")
}

func (b *keyspace) Fetch(keys []string) ([]catalog.Pair, errors.Error) {
	rv := make([]catalog.Pair, len(keys))
	for i, k := range keys {
		item, e := b.FetchOne(k)
		if e != nil {
			return nil, e
		}

		rv[i].Key = k
		rv[i].Value = item
	}

	return rv, nil
}

func (b *keyspace) FetchOne(key string) (value.Value, errors.Error) {
	path := filepath.Join(b.path(), key+".json")
	item, e := fetch(path)
	if e != nil {
		item = nil
	}

	return item, e
}

func (b *keyspace) Insert(inserts []catalog.Pair) ([]catalog.Pair, errors.Error) {
	// FIXME
	return nil, errors.NewError(nil, "Not yet implemented.")
}

func (b *keyspace) Update(updates []catalog.Pair) ([]catalog.Pair, errors.Error) {
	// FIXME
	return nil, errors.NewError(nil, "Not yet implemented.")
}

func (b *keyspace) Upsert(upserts []catalog.Pair) ([]catalog.Pair, errors.Error) {
	// FIXME
	return nil, errors.NewError(nil, "Not yet implemented.")
}

func (b *keyspace) Delete(deletes []string) errors.Error {
	// FIXME
	return errors.NewError(nil, "Not yet implemented.")
}

func (b *keyspace) Release() {
}

func (b *keyspace) path() string {
	return filepath.Join(b.namespace.path(), b.name)
}

// newKeyspace creates a new keyspace.
func newKeyspace(p *namespace, dir string) (b *keyspace, e errors.Error) {
	b = new(keyspace)
	b.namespace = p
	b.name = dir

	fi, er := os.Stat(b.path())
	if er != nil {
		return nil, errors.NewError(er, "")
	}

	if !fi.IsDir() {
		return nil, errors.NewError(nil, "Keyspace path must be a directory.")
	}

	b.indexes = make(map[string]catalog.Index, 1)
	pi := new(primaryIndex)
	b.primary = pi
	pi.keyspace = b
	pi.name = "#primary"
	b.indexes[pi.name] = pi

	return
}

// primaryIndex performs full keyspace scans.
type primaryIndex struct {
	name     string
	keyspace *keyspace
}

func (pi *primaryIndex) KeyspaceId() string {
	return pi.keyspace.Id()
}

func (pi *primaryIndex) Id() string {
	return pi.Name()
}

func (pi *primaryIndex) Name() string {
	return pi.name
}

func (pi *primaryIndex) Type() catalog.IndexType {
	return catalog.UNSPECIFIED
}

func (pi *primaryIndex) Drop() errors.Error {
	return errors.NewError(nil, "This primary index cannot be dropped.")
}

func (pi *primaryIndex) EqualKey() expression.Expressions {
	return nil
}

func (pi *primaryIndex) RangeKey() expression.Expressions {
	// FIXME
	return nil
}

func (pi *primaryIndex) Condition() expression.Expression {
	return nil
}

func (pi *primaryIndex) Statistics(span *catalog.Span) (catalog.Statistics, errors.Error) {
	return nil, nil
}

func (pi *primaryIndex) Scan(span *catalog.Span, distinct bool, limit int64, conn *catalog.IndexConnection) {
	defer close(conn.EntryChannel())

	// For primary indexes, bounds must always be strings, so we
	// can just enforce that directly
	low, high := "", ""

	// Ensure that lower bound is a string, if any
	if len(span.Range.Low) > 0 {
		a := span.Range.Low[0].Actual()
		switch a := a.(type) {
		case string:
			low = a
		default:
			conn.SendError(errors.NewError(nil, fmt.Sprintf("Invalid lower bound %v of type %T.", a, a)))
			return
		}
	}

	// Ensure that upper bound is a string, if any
	if len(span.Range.High) > 0 {
		a := span.Range.High[0].Actual()
		switch a := a.(type) {
		case string:
			high = a
		default:
			conn.SendError(errors.NewError(nil, fmt.Sprintf("Invalid upper bound %v of type %T.", a, a)))
			return
		}
	}

	dirEntries, er := ioutil.ReadDir(pi.keyspace.path())
	if er != nil {
		conn.SendError(errors.NewError(er, ""))
		return
	}

	var n int64 = 0
	for _, dirEntry := range dirEntries {
		if limit > 0 && n > limit {
			break
		}

		id := documentPathToId(dirEntry.Name())

		if low != "" &&
			(id < low ||
				(id == low && (span.Range.Inclusion&catalog.LOW == 0))) {
			continue
		}

		low = ""

		if high != "" &&
			(id > high ||
				(id == high && (span.Range.Inclusion&catalog.HIGH == 0))) {
			break
		}

		if !dirEntry.IsDir() {
			entry := catalog.IndexEntry{PrimaryKey: id}
			conn.EntryChannel() <- &entry
			n++
		}
	}
}

func (pi *primaryIndex) ScanEntries(limit int64, conn *catalog.IndexConnection) {
	defer close(conn.EntryChannel())

	dirEntries, er := ioutil.ReadDir(pi.keyspace.path())
	if er != nil {
		conn.SendError(errors.NewError(er, ""))
		return
	}

	for i, dirEntry := range dirEntries {
		if limit > 0 && int64(i) > limit {
			break
		}
		if !dirEntry.IsDir() {
			entry := catalog.IndexEntry{PrimaryKey: documentPathToId(dirEntry.Name())}
			conn.EntryChannel() <- &entry
		}
	}
}

func fetch(path string) (item value.Value, e errors.Error) {
	bytes, er := ioutil.ReadFile(path)
	if er != nil {
		if os.IsNotExist(er) {
			// file doesn't exist should simply return nil, nil
			return
		}
		return nil, errors.NewError(er, "")
	}

	doc := value.NewAnnotatedValue(value.NewValueFromBytes(bytes))
	doc.SetAttachment("meta", map[string]interface{}{"id": documentPathToId(path)})
	item = doc

	return
}

func documentPathToId(p string) string {
	_, file := filepath.Split(p)
	ext := filepath.Ext(file)
	return file[0 : len(file)-len(ext)]
}