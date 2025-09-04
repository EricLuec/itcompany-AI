package employees

import (
	"encoding/json"
	sector2 "itCompany-AI/sector"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateEmployee(t *testing.T) {
	origRandom := getRandomUserFunc
	origSector := getSectorFunc
	getRandomUserFunc = func() (string, string, error) { return "John", "Doe", nil }
	getSectorFunc = func() (*sector2.Sector, error) {
		return &sector2.Sector{
			Name:        "Engineering",
			Description: "Building stuff",
			SalaryClass: "A",
		}, nil
	}
	defer func() {
		getRandomUserFunc = origRandom
		getSectorFunc = origSector
	}()

	emp, err := GenerateEmployee()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if emp.FirstName != "John" || emp.LastName != "Doe" {
		t.Errorf("unexpected name: %s %s", emp.FirstName, emp.LastName)
	}
	if emp.Sector.Name != "Engineering" {
		t.Errorf("unexpected sector name: %s", emp.Sector.Name)
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

	origURL := employeeAPI
	employeeAPI = server.URL
	defer func() { employeeAPI = origURL }()

	emp := &Employee{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		HireDate:  "2025-09-03",
		Salary:    1000,
		Sector: sector2.Sector{
			Name:        "IT",
			Description: "",
			SalaryClass: "A",
		},
	}

	if err := PostEmployee(emp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllEmployeeIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]int{1, 2, 3})
	}))
	defer server.Close()

	origURL := employeeAllIDsAPI
	employeeAllIDsAPI = server.URL
	defer func() { employeeAllIDsAPI = origURL }()

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

	origURL := employeeAPI
	employeeAPI = server.URL
	defer func() { employeeAPI = origURL }()

	if err := DeleteEmployee(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
