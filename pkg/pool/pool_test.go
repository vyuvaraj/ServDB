package pool

import (
	"context"
	"testing"
)

func TestNewConnectionPoolPostgres(t *testing.T) {
	p := NewConnectionPool(5, "postgres")
	defer p.Shutdown(context.Background())
	if p.dialect != "postgres" {
		t.Errorf("expected dialect postgres, got %s", p.dialect)
	}
	if p.maxConns != 5 {
		t.Errorf("expected maxConns 5, got %d", p.maxConns)
	}
}

func TestNewConnectionPoolSQLite(t *testing.T) {
	p := NewConnectionPool(10, "sqlite")
	defer p.Shutdown(context.Background())
	if p.dialect != "sqlite" {
		t.Errorf("expected dialect sqlite, got %s", p.dialect)
	}
}

func TestNewConnectionPoolInvalidDialect(t *testing.T) {
	p := NewConnectionPool(3, "invalid")
	defer p.Shutdown(context.Background())
	if p.dialect != "invalid" {
		t.Errorf("expected invalid, got %s", p.dialect)
	}
}

func TestConnectionPoolAcquireRelease(t *testing.T) {
	p := NewConnectionPool(2, "sqlite")
	defer p.Shutdown(context.Background())

	conn, err := p.Acquire()
	if err != nil {
		t.Fatalf("failed to acquire connection: %v", err)
	}
	if conn == nil {
		t.Fatal("acquired connection is nil")
	}

	p.Release(conn)
}

func TestConnectionPoolStats(t *testing.T) {
	p := NewConnectionPool(3, "postgres")
	defer p.Shutdown(context.Background())

	stats := p.Stats()
	if stats.Dialect != "postgres" {
		t.Errorf("expected stats dialect postgres, got %s", stats.Dialect)
	}
	if stats.MaxConnections != 3 {
		t.Errorf("expected MaxConnections 3, got %d", stats.MaxConnections)
	}
}

func TestConnectionPoolExhaustion(t *testing.T) {
	p := NewConnectionPool(1, "postgres")
	defer p.Shutdown(context.Background())

	c1, err := p.Acquire()
	if err != nil {
		t.Fatalf("c1 acquire failed: %v", err)
	}

	// Active conns: 1, baseMaxConns: 1. Adaptive size scales maxConns to 2
	c2, err := p.Acquire()
	if err != nil {
		t.Fatalf("c2 acquire failed (should adaptively scale): %v", err)
	}

	// Should exhaust now since we hit scaled maxConns (2)
	_, err = p.Acquire()
	if err == nil {
		t.Error("expected pool exhaustion error, got nil")
	}

	p.Release(c1)
	p.Release(c2)
}

func TestConnectionPoolReset(t *testing.T) {
	p := NewConnectionPool(2, "postgres")
	defer p.Shutdown(context.Background())
	p.IncrementQueries()
	stats := p.Stats()
	if stats.TotalQueries != 1 {
		t.Errorf("expected queries 1, got %d", stats.TotalQueries)
	}
}

func TestConnectionPoolCapacity(t *testing.T) {
	p := NewConnectionPool(4, "sqlite")
	defer p.Shutdown(context.Background())
	if p.Dialect() != "sqlite" {
		t.Errorf("expected sqlite dialect, got %s", p.Dialect())
	}
}
