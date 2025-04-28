package main

import (
	"errors"
	"fmt"
	"sync"
)

type Operator string

const (
	Eq  Operator = "=="
	Gt  Operator = ">"
	Lt  Operator = "<"
	Gte Operator = ">="
	Lte Operator = "<="
)

type LogicalOperator string

const (
	And LogicalOperator = "AND"
	Or  LogicalOperator = "OR"
)

// Query interface (Composite Pattern)
type Query interface {
	Evaluate(row map[string]interface{}) bool
}

// Leaf: Single Condition
type Condition struct {
	Column   string
	Operator Operator
	Value    interface{}
}

func (c *Condition) Evaluate(row map[string]interface{}) bool {
	val, exists := row[c.Column]
	if !exists {
		return false
	}
	return compare(val, c.Value, c.Operator)
}

// Composite: Logical Combination of Queries
type CompositeFilter struct {
	LogicalOp LogicalOperator
	Children  []Query
}

func (cf *CompositeFilter) Evaluate(row map[string]interface{}) bool {
	switch cf.LogicalOp {
	case And:
		for _, child := range cf.Children {
			if !child.Evaluate(row) {
				return false
			}
		}
		return true
	case Or:
		for _, child := range cf.Children {
			if child.Evaluate(row) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Data types
type ColumnDataType interface {
	Validate(val interface{}) error
}

type IntDataType struct {
	MinValue int
	MaxValue int
}

func (i *IntDataType) Validate(val interface{}) error {
	v, ok := val.(int)
	if !ok {
		return errors.New("expected int")
	}
	if v < i.MinValue || v > i.MaxValue {
		return fmt.Errorf("int %d out of bounds (%d-%d)", v, i.MinValue, i.MaxValue)
	}
	return nil
}

type StringDataType struct {
	AllowNull bool
}

func (s *StringDataType) Validate(val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("expected string")
	}
	if !s.AllowNull && v == "" {
		return errors.New("empty string not allowed")
	}
	return nil
}

type SchemaMember struct {
	Name     string
	DataType ColumnDataType
	Required bool
}

type Schema struct {
	Columns map[string]SchemaMember
}

func NewSchema(members []SchemaMember) *Schema {
	cols := make(map[string]SchemaMember)
	for _, m := range members {
		cols[m.Name] = m
	}
	return &Schema{Columns: cols}
}

func (s *Schema) Validate(row map[string]interface{}) error {
	for name, member := range s.Columns {
		val, exists := row[name]
		if !exists {
			if member.Required {
				return fmt.Errorf("missing required field: %s", name)
			}
			continue
		}
		if err := member.DataType.Validate(val); err != nil {
			return fmt.Errorf("validation failed for %s: %v", name, err)
		}
	}
	return nil
}

// Indexing
type Index struct {
	ColumnName string
	IndexMap   map[interface{}]map[int]struct{}
}

func NewIndex(column string) *Index {
	return &Index{
		ColumnName: column,
		IndexMap:   make(map[interface{}]map[int]struct{}),
	}
}

func (idx *Index) Add(value interface{}, id int) {
	if _, exists := idx.IndexMap[value]; !exists {
		idx.IndexMap[value] = make(map[int]struct{})
	}
	idx.IndexMap[value][id] = struct{}{}
}

func (idx *Index) Remove(value interface{}, id int) {
	if rows, exists := idx.IndexMap[value]; exists {
		delete(rows, id)
		if len(rows) == 0 {
			delete(idx.IndexMap, value)
		}
	}
}

// Table
type Table struct {
	Name      string
	Schema    *Schema
	Data      map[int]map[string]interface{}
	AutoID    int
	Indexes   map[string]*Index
	DataLock  sync.RWMutex
	IndexLock sync.RWMutex
}

func NewTable(name string, schema *Schema) *Table {
	return &Table{
		Name:    name,
		Schema:  schema,
		Data:    make(map[int]map[string]interface{}),
		Indexes: make(map[string]*Index),
	}
}

func (t *Table) Insert(row map[string]interface{}) (int, error) {
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	t.AutoID++
	row["id"] = t.AutoID

	if err := t.Schema.Validate(row); err != nil {
		return 0, err
	}
	t.Data[t.AutoID] = row

	for col, idx := range t.Indexes {
		if val, ok := row[col]; ok {
			idx.Add(val, t.AutoID)
		}
	}

	return t.AutoID, nil
}

func (t *Table) Update(id int, updated map[string]interface{}) error {
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	row, exists := t.Data[id]
	if !exists {
		return errors.New("row not found")
	}

	for k, v := range updated {
		row[k] = v
	}

	if err := t.Schema.Validate(row); err != nil {
		return err
	}

	for col, idx := range t.Indexes {
		if val, ok := updated[col]; ok {
			idx.Remove(row[col], id)
			idx.Add(val, id)
		}
	}

	return nil
}

func (t *Table) Delete(id int) error {
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	row, exists := t.Data[id]
	if !exists {
		return errors.New("row not found")
	}

	for col, idx := range t.Indexes {
		if val, ok := row[col]; ok {
			idx.Remove(val, id)
		}
	}

	delete(t.Data, id)
	return nil
}

func (t *Table) CreateIndex(column string) {
	t.IndexLock.Lock()
	defer t.IndexLock.Unlock()

	idx := NewIndex(column)
	for id, row := range t.Data {
		if val, ok := row[column]; ok {
			idx.Add(val, id)
		}
	}
	t.Indexes[column] = idx
}

func compare(v1 interface{}, v2 interface{}, op Operator) bool {
	switch a := v1.(type) {
	case int:
		b, _ := v2.(int)
		switch op {
		case Eq:
			return a == b
		case Gt:
			return a > b
		case Lt:
			return a < b
		case Gte:
			return a >= b
		case Lte:
			return a <= b
		}
	case string:
		b, _ := v2.(string)
		switch op {
		case Eq:
			return a == b
		}
	}
	return false
}

// New Query method using Composite
func (t *Table) Query(q Query) ([]map[string]interface{}, error) {
	t.DataLock.RLock()
	defer t.DataLock.RUnlock()

	var result []map[string]interface{}
	for _, row := range t.Data {
		if row == nil {
			continue
		}
		if q.Evaluate(row) {
			result = append(result, row)
		}
	}
	return result, nil
}

// Database
type Database struct {
	Name   string
	Tables map[string]*Table
}

func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: make(map[string]*Table),
	}
}

func (db *Database) CreateTable(name string, schema *Schema) {
	db.Tables[name] = NewTable(name, schema)
}

// Server
type Server struct {
	Databases map[string]*Database
}

func NewServer() *Server {
	return &Server{
		Databases: make(map[string]*Database),
	}
}

func (s *Server) CreateDatabase(name string) {
	s.Databases[name] = NewDatabase(name)
}

func main() {
	server := NewServer()

	server.CreateDatabase("testdb")
	db := server.Databases["testdb"]

	schema := NewSchema([]SchemaMember{
		{Name: "id", DataType: &IntDataType{MinValue: 0, MaxValue: 10000}, Required: true},
		{Name: "name", DataType: &StringDataType{AllowNull: false}, Required: true},
		{Name: "age", DataType: &IntDataType{MinValue: 0, MaxValue: 150}, Required: true},
		{Name: "city", DataType: &StringDataType{AllowNull: true}, Required: false},
	})

	db.CreateTable("users", schema)
	users := db.Tables["users"]

	users.Insert(map[string]interface{}{"name": "Alice", "age": 30, "city": "Paris"})
	users.Insert(map[string]interface{}{"name": "Bob", "age": 25, "city": "London"})
	users.Insert(map[string]interface{}{"name": "Alice", "age": 35, "city": "Berlin"})
	users.Insert(map[string]interface{}{"name": "Charlie", "age": 28, "city": "Paris"})

	// Create Index
	users.CreateIndex("name")

	// Build a complex query:
	// (name == "Alice" AND age > 30) OR (city == "Paris")
	query := &CompositeFilter{
		LogicalOp: Or,
		Children: []Query{
			&CompositeFilter{
				LogicalOp: And,
				Children: []Query{
					&Condition{Column: "name", Operator: Eq, Value: "Alice"},
					&Condition{Column: "age", Operator: Gt, Value: 30},
				},
			},
			&Condition{Column: "city", Operator: Eq, Value: "Paris"},
		},
	}

	results, _ := users.Query(query)
	for _, r := range results {
		fmt.Println(r)
	}
}
