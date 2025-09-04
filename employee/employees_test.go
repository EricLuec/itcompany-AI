package employees

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Hilfsfunktionen für Mocking
func mockGetRandomUser() (string, string, error) {
	return "John", "Doe", nil
}

func mockGetOneSector() (*Sector, error) {
	return &Sector{
		ID:          1,
		Name:        "Engineering",
		Description: "Building cool stuff",
		SalaryClass: "A",
	}, nil
}

func TestGenerateEmployee(t *testing.T) {
	// Temporär die echten Funktionen überschreiben
	origRandomUser := getRandomUser
	origGetSector := GetOneSector

	getRandomUser = mockGetRandomUser
	GetOneSector = func() (*Sector, error) { return mockGetOneSector() }

	defer func() {
		getRandomUser = origRandomUser
		GetOneSector = origGetSector
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

	// Temporär die URL anpassen
	origURL := employeeAPI
	employeeAPI = server.URL
	defer func() { employeeAPI = origURL }()

	emp := &Employee{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		HireDate:  "2025-09-03",
		Salary:    1000,
		Sector: Sector{
			ID:          1,
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
