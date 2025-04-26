package main

import (
	"errors"
	"fmt"
)

//////////////////////////////////////////////////////
// 🧱 Core Entities
//////////////////////////////////////////////////////

type Column struct {
	Name string
	Type string
}

type Row map[string]interface{}

//////////////////////////////////////////////////////
// 🧠 Strategy Pattern for Filtering
//////////////////////////////////////////////////////

type FilterStrategy interface {
	Match(row Row) bool
}

type EqualFilter struct {
	Column string
	Value  interface{}
}

func (f EqualFilter) Match(row Row) bool {
	val, exists := row[f.Column]
	return exists && val == f.Value
}

//////////////////////////////////////////////////////
// 🗃️ Repository Interface (Interface Segregation)
//////////////////////////////////////////////////////

type TableRepository interface {
	Insert(row Row) error
	Select(filters []FilterStrategy) ([]Row, error)
	Update(filters []FilterStrategy, updates Row) error
	Delete(filters []FilterStrategy) error
}

//////////////////////////////////////////////////////
// 💾 In-Memory Table Implementation (SRP)
//////////////////////////////////////////////////////

type InMemoryTable struct {
	Columns []Column
	Rows    []Row
}

func (t *InMemoryTable) validate(row Row) error {
	for _, col := range t.Columns {
		if _, ok := row[col.Name]; !ok {
			return errors.New("missing column: " + col.Name)
		}
	}
	return nil
}

func (t *InMemoryTable) Insert(row Row) error {
	if err := t.validate(row); err != nil {
		return err
	}
	t.Rows = append(t.Rows, row)
	return nil
}

func (t *InMemoryTable) Select(filters []FilterStrategy) ([]Row, error) {
	var result []Row
	for _, row := range t.Rows {
		match := true
		for _, f := range filters {
			if !f.Match(row) {
				match = false
				break
			}
		}
		if match {
			result = append(result, row)
		}
	}
	return result, nil
}

func (t *InMemoryTable) Update(filters []FilterStrategy, updates Row) error {
	for idx, row := range t.Rows {
		match := true
		for _, f := range filters {
			if !f.Match(row) {
				match = false
				break
			}
		}
		if match {
			for k, v := range updates {
				row[k] = v
			}
			t.Rows[idx] = row
		}
	}
	return nil
}

func (t *InMemoryTable) Delete(filters []FilterStrategy) error {
	var newRows []Row
	for _, row := range t.Rows {
		match := true
		for _, f := range filters {
			if !f.Match(row) {
				match = false
				break
			}
		}
		if !match {
			newRows = append(newRows, row)
		}
	}
	t.Rows = newRows
	return nil
}

//////////////////////////////////////////////////////
// 🏭 Factory Pattern for Table Creation
//////////////////////////////////////////////////////

func NewInMemoryTable(columns []Column) *InMemoryTable {
	return &InMemoryTable{
		Columns: columns,
		Rows:    []Row{},
	}
}

//////////////////////////////////////////////////////
// 📦 Database Manager (Open-Closed, Dependency Inversion)
//////////////////////////////////////////////////////

type Database struct {
	Tables map[string]TableRepository
}

func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]TableRepository),
	}
}

func (db *Database) RegisterTable(name string, repo TableRepository) {
	db.Tables[name] = repo
}

func (db *Database) GetTable(name string) TableRepository {
	return db.Tables[name]
}

//////////////////////////////////////////////////////
// 🚀 MAIN: Usage Example
//////////////////////////////////////////////////////

func main() {
	db := NewDatabase()

	// Create a "users" table
	usersTable := NewInMemoryTable([]Column{
		{"id", "int"},
		{"name", "string"},
		{"email", "string"},
	})

	db.RegisterTable("users", usersTable)

	users := db.GetTable("users")
	users.Insert(Row{"id": 1, "name": "Alice", "email": "alice@example.com"})
	users.Insert(Row{"id": 2, "name": "Bob", "email": "bob@example.com"})
	users.Insert(Row{"id": 3, "name": "Charlie", "email": "charlie@example.com"})

	// SELECT: name = "Alice"
	fmt.Println("\n🔍 SELECT where name = 'Alice'")
	results, _ := users.Select([]FilterStrategy{
		EqualFilter{"name", "Alice"},
	})
	fmt.Println(results)

	// UPDATE: Update email where id = 2
	fmt.Println("\n✏️ UPDATE email where id = 2")
	users.Update([]FilterStrategy{
		EqualFilter{"id", 2},
	}, Row{"email": "newbob@example.com"})

	// SELECT ALL
	fmt.Println("\n📋 SELECT ALL users after update")
	allResults, _ := users.Select([]FilterStrategy{})
	fmt.Println(allResults)

	// DELETE: Delete user where name = "Charlie"
	fmt.Println("\n🗑️ DELETE where name = 'Charlie'")
	users.Delete([]FilterStrategy{
		EqualFilter{"name", "Charlie"},
	})

	// SELECT ALL after delete
	fmt.Println("\n📋 SELECT ALL users after delete")
	finalResults, _ := users.Select([]FilterStrategy{})
	fmt.Println(finalResults)
}
