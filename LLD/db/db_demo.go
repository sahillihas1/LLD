package main

import (
	"errors"
	"fmt"
	"sync"
)

// Schema
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

// Table
type Table struct {
	Name     string
	Schema   *Schema
	Data     map[int]map[string]interface{}
	AutoID   int
	DataLock sync.RWMutex
}

func NewTable(name string, schema *Schema) *Table {
	return &Table{
		Name:   name,
		Schema: schema,
		Data:   make(map[int]map[string]interface{}),
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

	return t.AutoID, nil
}

func (t *Table) QuerySimple(filters map[string]interface{}) ([]map[string]interface{}, error) {
	t.DataLock.RLock()
	defer t.DataLock.RUnlock()

	var results []map[string]interface{}

	for _, row := range t.Data {
		if row == nil {
			continue
		}
		match := true
		for col, val := range filters {
			if rowVal, ok := row[col]; !ok || rowVal != val {
				match = false
				break
			}
		}
		if match {
			results = append(results, row)
		}
	}

	return results, nil
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

// Main function
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

	// Simple Query: name = Alice, age = 30
	filters := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	results, _ := users.QuerySimple(filters)
	for _, r := range results {
		fmt.Println(r)
	}
}
