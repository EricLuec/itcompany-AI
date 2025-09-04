package employees

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateEmployee(t *testing.T) {
	origSector := getSector
	getSector = func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"name":        "Engineering",
			"description": "Building cool stuff",
			"salaryClass": "A",
		}, nil
	}
	defer func() { getSector = origSector }()

	origRandom := getRandomUserFunc
	getRandomUserFunc = func() (string, string, error) { return "John", "Doe", nil }
	defer func() { getRandomUserFunc = origRandom }()

	emp, err := GenerateEmployee()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if emp.FirstName != "John" || emp.LastName != "Doe" {
		t.Errorf("unexpected name: %s %s", emp.FirstName, emp.LastName)
	}
}

func TestPostEmployee(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var emp Employee
		if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
			t.Errorf("decode error: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	orig := employeeAPI
	employeeAPI = server.URL
	defer func() { employeeAPI = orig }()

	err := PostEmployee(&Employee{FirstName: "Alice", LastName: "Smith"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllEmployeeIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]int{1, 2, 3})
	}))
	defer server.Close()

	orig := employeeAllIDsAPI
	employeeAllIDsAPI = server.URL
	defer func() { employeeAllIDsAPI = orig }()

	ids, err := GetAllEmployeeIDs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %v", ids)
	}
}

func TestDeleteEmployee(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	orig := employeeAPI
	employeeAPI = server.URL
	defer func() { employeeAPI = orig }()

	if err := DeleteEmployee(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
