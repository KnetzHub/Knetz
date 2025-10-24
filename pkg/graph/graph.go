package graph

import (
	"fmt"
	
	"github.com/knetz-io/knetz/internal/models"
)

// Graph represents a dependency graph
type Graph struct {
	nodes map[string]*models.GraphNode
	edges []models.GraphEdge
}

// NewGraph creates a new dependency graph
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*models.GraphNode),
		edges: make([]models.GraphEdge, 0),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(node *models.GraphNode) {
	g.nodes[node.ID] = node
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to string, required bool, confidence float64) error {
	if _, ok := g.nodes[from]; !ok {
		return fmt.Errorf("source node %s not found", from)
	}
	if _, ok := g.nodes[to]; !ok {
		return fmt.Errorf("target node %s not found", to)
	}

	g.edges = append(g.edges, models.GraphEdge{
		From:       from,
		To:         to,
		Required:   required,
		Confidence: confidence,
	})

	return nil
}

// GetNode returns a node by ID
func (g *Graph) GetNode(id string) (*models.GraphNode, bool) {
	node, ok := g.nodes[id]
	return node, ok
}

// GetNodes returns all nodes
func (g *Graph) GetNodes() map[string]*models.GraphNode {
	return g.nodes
}

// GetEdges returns all edges
func (g *Graph) GetEdges() []models.GraphEdge {
	return g.edges
}

// GetDependencies returns all dependencies of a node
func (g *Graph) GetDependencies(nodeID string) []string {
	var deps []string
	for _, edge := range g.edges {
		if edge.From == nodeID {
			deps = append(deps, edge.To)
		}
	}
	return deps
}

// GetDependents returns all nodes that depend on the given node
func (g *Graph) GetDependents(nodeID string) []string {
	var dependents []string
	for _, edge := range g.edges {
		if edge.To == nodeID {
			dependents = append(dependents, edge.From)
		}
	}
	return dependents
}

// DetectCycles detects circular dependencies in the graph
func (g *Graph) DetectCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make([]string, 0)

	for nodeID := range g.nodes {
		if !visited[nodeID] {
			if cycleFound := g.dfs(nodeID, visited, recStack, path, &cycles); cycleFound {
				// Cycle detected
			}
		}
	}

	return cycles
}

// dfs performs depth-first search to detect cycles
func (g *Graph) dfs(nodeID string, visited, recStack map[string]bool, path []string, cycles *[][]string) bool {
	visited[nodeID] = true
	recStack[nodeID] = true
	path = append(path, nodeID)

	for _, dep := range g.GetDependencies(nodeID) {
		if !visited[dep] {
			if g.dfs(dep, visited, recStack, path, cycles) {
				return true
			}
		} else if recStack[dep] {
			// Cycle detected - extract the cycle
			cycle := make([]string, 0)
			cycleStart := -1
			for i, n := range path {
				if n == dep {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle = append(cycle, path[cycleStart:]...)
				cycle = append(cycle, dep) // Close the cycle
				*cycles = append(*cycles, cycle)
			}
			return true
		}
	}

	recStack[nodeID] = false
	return false
}

// CalculateDepth calculates the dependency depth for a node
func (g *Graph) CalculateDepth(nodeID string) int {
	visited := make(map[string]bool)
	return g.calculateDepthRecursive(nodeID, visited)
}

// calculateDepthRecursive recursively calculates depth
func (g *Graph) calculateDepthRecursive(nodeID string, visited map[string]bool) int {
	if visited[nodeID] {
		return 0 // Avoid infinite loops
	}

	visited[nodeID] = true
	maxDepth := 0

	for _, dep := range g.GetDependencies(nodeID) {
		depth := g.calculateDepthRecursive(dep, visited)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth + 1
}

// CalculateImpactScore calculates how many services depend on this node
func (g *Graph) CalculateImpactScore(nodeID string) int {
	return len(g.GetDependents(nodeID))
}

// GetTransitiveDependencies returns all transitive dependencies of a node
func (g *Graph) GetTransitiveDependencies(nodeID string) []string {
	visited := make(map[string]bool)
	var deps []string
	g.collectTransitiveDeps(nodeID, visited, &deps)
	return deps
}

// collectTransitiveDeps recursively collects transitive dependencies
func (g *Graph) collectTransitiveDeps(nodeID string, visited map[string]bool, deps *[]string) {
	if visited[nodeID] {
		return
	}
	visited[nodeID] = true

	for _, dep := range g.GetDependencies(nodeID) {
		*deps = append(*deps, dep)
		g.collectTransitiveDeps(dep, visited, deps)
	}
}

// ToModel converts the graph to a model representation
func (g *Graph) ToModel() *models.DependencyGraph {
	return &models.DependencyGraph{
		Nodes: g.nodes,
		Edges: g.edges,
	}
}

// Builder helps build dependency graphs from services
type Builder struct {
	graph *Graph
}

// NewBuilder creates a new graph builder
func NewBuilder() *Builder {
	return &Builder{
		graph: NewGraph(),
	}
}

// AddService adds a service and its dependencies to the graph
func (b *Builder) AddService(service *models.Service) {
	// Create node for the service
	nodeID := fmt.Sprintf("%s@%s:%s", service.Name, service.ClusterName, service.Namespace)
	node := &models.GraphNode{
		ID:          nodeID,
		ServiceName: service.Name,
		Version:     service.Version,
		Cluster:     service.ClusterName,
		Namespace:   service.Namespace,
		Tenant:      service.TenantName,
	}
	b.graph.AddNode(node)

	// Add edges for dependencies
	for _, dep := range service.Dependencies {
		depNodeID := fmt.Sprintf("%s@%s:%s", dep.ServiceName, dep.Cluster, dep.Namespace)
		
		// Create dependency node if it doesn't exist
		if _, exists := b.graph.GetNode(depNodeID); !exists {
			depNode := &models.GraphNode{
				ID:          depNodeID,
				ServiceName: dep.ServiceName,
				Version:     dep.Version,
				Cluster:     dep.Cluster,
				Namespace:   dep.Namespace,
				Tenant:      dep.Tenant,
			}
			b.graph.AddNode(depNode)
		}

		// Add edge
		b.graph.AddEdge(nodeID, depNodeID, dep.Required, dep.Confidence)
	}
}

// Build returns the constructed graph
func (b *Builder) Build() *Graph {
	return b.graph
}

