package building

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostBuilding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b Building
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			t.Errorf("decode error: %v", err)
		}
		if b.Name != "SkyTower" || b.City != "Berlin, Germany" {
			t.Errorf("unexpected building: %+v", b)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	origURL := buildingAPI
	buildingAPI = server.URL
	defer func() { buildingAPI = origURL }()

	b := &Building{Name: "SkyTower", Description: "", BuildingDate: "2025-09-03", Capacity: 1000, City: "Berlin, Germany"}
	if err := PostBuilding(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllBuildingIds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]int{1, 2, 3})
	}))
	defer server.Close()

	origURL := buildingAllIDsAPI
	buildingAllIDsAPI = server.URL
	defer func() { buildingAllIDsAPI = origURL }()

	ids, err := GetAllBuildingIds()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %v", ids)
	}
}

func TestDeleteBuilding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	origURL := buildingAPI
	buildingAPI = server.URL
	defer func() { buildingAPI = origURL }()

	if err := DeleteBuilding(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
