package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"servpool/pkg/pool"
)

type mockPool struct {
	dialectVal string
}

func (m *mockPool) Acquire() (*pool.DbConn, error)              { return &pool.DbConn{ID: 1}, nil }
func (m *mockPool) Release(conn *pool.DbConn)                   {}
func (m *mockPool) IncrementQueries()                           {}
func (m *mockPool) Stats() pool.PoolStats                       { return pool.PoolStats{Dialect: m.dialectVal} }
func (m *mockPool) Dialect() string                             { return m.dialectVal }
func (m *mockPool) Shutdown(ctx context.Context) error          { return nil }

func TestRoutingServerMetrics(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()
	srv.HandlePrometheusMetrics(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRoutingServerDbHealth(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("GET", "/api/db/health", nil)
	rr := httptest.NewRecorder()
	srv.HandleDbHealth(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRoutingAddRegionPool(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)
	rPool := &mockPool{dialectVal: "mysql"}
	srv.AddRegionPool("us-west", rPool)

	srv.regionPoolsMu.RLock()
	p := srv.regionPools["us-west"]
	srv.regionPoolsMu.RUnlock()
	if p != rPool {
		t.Error("failed to register regional pool")
	}
}

func TestRoutingHasActiveTable(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)
	srv.activeTables["users"] = true
	if !srv.HasActiveTable("users") {
		t.Error("expected users table to be active")
	}
}

func TestRoutingSetPeers(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)
	srv.SetPeers([]string{"peer1"})
	if len(srv.peers) != 1 || srv.peers[0] != "peer1" {
		t.Errorf("expected peers [peer1], got %v", srv.peers)
	}
}

func TestRoutingShutdown(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)
	err := srv.Shutdown(context.Background())
	if err != nil {
		t.Errorf("failed to shutdown gracefully: %v", err)
	}
}

func TestRoutingHandleQueryInvalidJSON(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("POST", "/api/db/query", bytes.NewReader([]byte("{invalid}")))
	rr := httptest.NewRecorder()
	srv.HandleQuery(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestRoutingHandleQueryEmptyQuery(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	body, _ := json.Marshal(QueryRequest{Query: ""})
	req := httptest.NewRequest("POST", "/api/db/query", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	srv.HandleQuery(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}
}

func TestRoutingHandleStats(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("GET", "/api/db/stats", nil)
	rr := httptest.NewRecorder()
	srv.HandleStats(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRoutingHandleAnalytics(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("GET", "/api/db/analytics", nil)
	rr := httptest.NewRecorder()
	srv.HandleAnalytics(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRoutingHandleClearCache(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("POST", "/api/db/cache/clear", nil)
	rr := httptest.NewRecorder()
	srv.HandleClearCache(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRoutingHandleMigrateInvalidJSON(t *testing.T) {
	primary := &mockPool{dialectVal: "postgres"}
	replica := &mockPool{dialectVal: "postgres"}
	srv := NewServer(primary, replica, nil)

	req := httptest.NewRequest("POST", "/api/db/migrate", bytes.NewReader([]byte("{invalid}")))
	rr := httptest.NewRecorder()
	srv.HandleMigrate(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rr.Code)
	}
}
