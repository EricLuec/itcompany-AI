package sector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostSector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var s Sector
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			t.Errorf("decode error: %v", err)
		}
		if s.Name != "IT" || s.SalaryClass != "A" {
			t.Errorf("unexpected sector: %+v", s)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	origURL := sectorAPI
	sectorAPI = server.URL
	defer func() { sectorAPI = origURL }()

	s := &Sector{Name: "IT", Description: "", SalaryClass: "A"}
	if err := PostSector(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllSectorIds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]int{1, 2, 3})
	}))
	defer server.Close()

	origURL := sectorAllIDsAPI
	sectorAllIDsAPI = server.URL
	defer func() { sectorAllIDsAPI = origURL }()

	ids, err := GetAllSectorIds()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %v", ids)
	}
}

func TestDeleteSector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	origURL := sectorAPI
	sectorAPI = server.URL
	defer func() { sectorAPI = origURL }()

	if err := DeleteSector(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetOneSector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Sector{
			{ID: 1, Name: "IT", Description: "Tech", SalaryClass: "A"},
		})
	}))
	defer server.Close()

	origURL := sectorAPI
	sectorAPI = server.URL
	defer func() { sectorAPI = origURL }()

	sec, err := GetOneSector()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sec.Name != "IT" {
		t.Errorf("expected IT, got %s", sec.Name)
	}
}
