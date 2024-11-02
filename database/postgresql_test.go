package database

import (
	"fmt"
	"os"
	"testing"
)

func TestCountAllPosts(t *testing.T) {

	Vars = map[string]string{
		"PG_HOST": "localhost",
		"PG_PORT": "15432",
		"PG_USER": "root",
		"PG_DB":   "blog",
		"PG_PASS": "postgres",
	}

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test CountAllPosts", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountAllPosts()
			fmt.Println("CountAllPosts() = ", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountAllPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CountAllPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Setup code before running tests
	setup()
	code := m.Run()
	// Teardown code after running tests
	teardown()
	os.Exit(code)
}

func setup() {
	// Initialize database connection or other setup tasks
	fmt.Println("Setup for tests")
}

func teardown() {
	// Cleanup tasks
	fmt.Println("Teardown after tests")
}
