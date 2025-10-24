package graph

import (
	"testing"
	
	"github.com/knetz-io/knetz/internal/models"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	
	if g == nil {
		t.Fatal("expected non-nil graph")
	}
	
	if len(g.GetNodes()) != 0 {
		t.Errorf("expected empty graph, got %d nodes", len(g.GetNodes()))
	}
}

func TestAddNode(t *testing.T) {
	g := NewGraph()
	
	node := &models.GraphNode{
		ID:          "service-a@cluster:namespace",
		ServiceName: "service-a",
		Version:     "1.0.0",
		Cluster:     "cluster",
		Namespace:   "namespace",
	}
	
	g.AddNode(node)
	
	if len(g.GetNodes()) != 1 {
		t.Errorf("expected 1 node, got %d", len(g.GetNodes()))
	}
	
	retrieved, exists := g.GetNode("service-a@cluster:namespace")
	if !exists {
		t.Error("node not found")
	}
	
	if retrieved.ServiceName != "service-a" {
		t.Errorf("expected 'service-a', got '%s'", retrieved.ServiceName)
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	
	nodeA := &models.GraphNode{ID: "a", ServiceName: "service-a"}
	nodeB := &models.GraphNode{ID: "b", ServiceName: "service-b"}
	
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	
	err := g.AddEdge("a", "b", true, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(g.GetEdges()) != 1 {
		t.Errorf("expected 1 edge, got %d", len(g.GetEdges()))
	}
	
	deps := g.GetDependencies("a")
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(deps))
	}
	
	if deps[0] != "b" {
		t.Errorf("expected dependency 'b', got '%s'", deps[0])
	}
}

func TestDetectCycles(t *testing.T) {
	g := NewGraph()
	
	// Create a cycle: a -> b -> c -> a
	nodeA := &models.GraphNode{ID: "a", ServiceName: "service-a"}
	nodeB := &models.GraphNode{ID: "b", ServiceName: "service-b"}
	nodeC := &models.GraphNode{ID: "c", ServiceName: "service-c"}
	
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)
	
	g.AddEdge("a", "b", true, 1.0)
	g.AddEdge("b", "c", true, 1.0)
	g.AddEdge("c", "a", true, 1.0)
	
	cycles := g.DetectCycles()
	
	if len(cycles) == 0 {
		t.Error("expected to detect cycle, but none found")
	}
}

func TestCalculateDepth(t *testing.T) {
	g := NewGraph()
	
	// Create chain: a -> b -> c
	nodeA := &models.GraphNode{ID: "a", ServiceName: "service-a"}
	nodeB := &models.GraphNode{ID: "b", ServiceName: "service-b"}
	nodeC := &models.GraphNode{ID: "c", ServiceName: "service-c"}
	
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)
	
	g.AddEdge("a", "b", true, 1.0)
	g.AddEdge("b", "c", true, 1.0)
	
	depth := g.CalculateDepth("a")
	if depth != 3 {
		t.Errorf("expected depth 3, got %d", depth)
	}
	
	depth = g.CalculateDepth("b")
	if depth != 2 {
		t.Errorf("expected depth 2, got %d", depth)
	}
	
	depth = g.CalculateDepth("c")
	if depth != 1 {
		t.Errorf("expected depth 1, got %d", depth)
	}
}

func TestCalculateImpactScore(t *testing.T) {
	g := NewGraph()
	
	// Create: a -> b, c -> b, d -> b (b has 3 dependents)
	nodeA := &models.GraphNode{ID: "a", ServiceName: "service-a"}
	nodeB := &models.GraphNode{ID: "b", ServiceName: "service-b"}
	nodeC := &models.GraphNode{ID: "c", ServiceName: "service-c"}
	nodeD := &models.GraphNode{ID: "d", ServiceName: "service-d"}
	
	g.AddNode(nodeA)
	g.AddNode(nodeB)
	g.AddNode(nodeC)
	g.AddNode(nodeD)
	
	g.AddEdge("a", "b", true, 1.0)
	g.AddEdge("c", "b", true, 1.0)
	g.AddEdge("d", "b", true, 1.0)
	
	impact := g.CalculateImpactScore("b")
	if impact != 3 {
		t.Errorf("expected impact score 3, got %d", impact)
	}
	
	impact = g.CalculateImpactScore("a")
	if impact != 0 {
		t.Errorf("expected impact score 0, got %d", impact)
	}
}

func TestGetTransitiveDependencies(t *testing.T) {
	g := NewGraph()
	
	// Create chain: a -> b -> c -> d
	nodes := []*models.GraphNode{
		{ID: "a", ServiceName: "service-a"},
		{ID: "b", ServiceName: "service-b"},
		{ID: "c", ServiceName: "service-c"},
		{ID: "d", ServiceName: "service-d"},
	}
	
	for _, node := range nodes {
		g.AddNode(node)
	}
	
	g.AddEdge("a", "b", true, 1.0)
	g.AddEdge("b", "c", true, 1.0)
	g.AddEdge("c", "d", true, 1.0)
	
	transitive := g.GetTransitiveDependencies("a")
	
	// Should include b, c, d
	if len(transitive) != 3 {
		t.Errorf("expected 3 transitive dependencies, got %d", len(transitive))
	}
}

