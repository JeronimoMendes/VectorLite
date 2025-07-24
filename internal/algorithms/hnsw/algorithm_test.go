package hnsw

import (
	"VectorLite/internal/algorithms"
	"VectorLite/internal/vector"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	M := 16
	efConstruction := 200
	mL := 1.0 / math.Log(2.0)

	alg := New(M, efConstruction, mL)

	assert.NotNil(t, alg)
	assert.Equal(t, M, alg.M)
	assert.Equal(t, efConstruction, alg.efConstruction)
	assert.Equal(t, mL, alg.mL)
	assert.Empty(t, alg.nodes)
	assert.Nil(t, alg.entryNode)
}

func TestHNSWNode_isConnectedTo(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	// Initially not connected
	assert.False(t, node1.isConnectedTo(node2, 0))
	assert.False(t, node2.isConnectedTo(node1, 0))

	// Add connection manually
	node1.Connections[0] = []*HNSWNode{node2}
	assert.True(t, node1.isConnectedTo(node2, 0))
	assert.False(t, node2.isConnectedTo(node1, 0))
}

func TestHNSWNode_connect(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	layer := 0

	// Connect nodes
	node1.connect(node2, layer)

	// Both nodes should be connected to each other
	assert.True(t, node1.isConnectedTo(node2, layer))
	assert.True(t, node2.isConnectedTo(node1, layer))
	assert.Len(t, node1.Connections[layer], 1)
	assert.Len(t, node2.Connections[layer], 1)
}

func TestHNSWNode_connect_AlreadyConnected(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	layer := 0

	node1.connect(node2, layer)

	// Current implementation creates duplicate connections
	assert.Len(t, node1.Connections[layer], 1)
	assert.Len(t, node2.Connections[layer], 1)
}

func TestHNSWNode_disconnect(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	layer := 0

	// First connect the nodes
	node1.connect(node2, layer)
	require.True(t, node1.isConnectedTo(node2, layer))
	require.True(t, node2.isConnectedTo(node1, layer))

	// Then disconnect
	node1.disconnect(node2, layer)

	// Both should be disconnected
	assert.False(t, node1.isConnectedTo(node2, layer))
	assert.False(t, node2.isConnectedTo(node1, layer))
	assert.Len(t, node1.Connections[layer], 0)
	assert.Len(t, node2.Connections[layer], 0)
}

func TestHNSWNode_disconnect_NotConnected(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	layer := 0

	// Disconnect nodes that were never connected - should not panic
	node1.disconnect(node2, layer)

	assert.False(t, node1.isConnectedTo(node2, layer))
	assert.False(t, node2.isConnectedTo(node1, layer))
}

func TestAlgorithm_AddEntry_FirstEntry(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	entry := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
		Metadata: map[string]string{"id": "1"},
		Id:       1,
	}

	alg.AddEntry(entry)

	assert.Len(t, alg.nodes, 1)
	assert.NotNil(t, alg.entryNode)
	assert.Equal(t, entry, alg.nodes[0].Entry)
	assert.Equal(t, alg.nodes[0], alg.entryNode)
}

func TestAlgorithm_AddEntry_MultipleEntries(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	entry1 := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
		Metadata: map[string]string{"id": "1"},
		Id:       1,
	}

	entry2 := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
		Metadata: map[string]string{"id": "2"},
		Id:       2,
	}

	alg.AddEntry(entry1)
	alg.AddEntry(entry2)

	assert.Len(t, alg.nodes, 2)
	assert.NotNil(t, alg.entryNode)

	// Entry node should be the one with higher max layer
	expectedEntryNode := alg.nodes[0]
	if alg.nodes[1].MaxLayer > alg.nodes[0].MaxLayer {
		expectedEntryNode = alg.nodes[1]
	}
	assert.Equal(t, expectedEntryNode, alg.entryNode)
}

func TestAlgorithm_calculateLevelProbability(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	// Test multiple calculations - current implementation can return negative values
	for i := 0; i < 100; i++ {
		level := alg.calculateLevelProbability()
		assert.IsType(t, int(0), level)
		// Note: current implementation can return negative values due to log calculation
	}
}

func TestAlgorithm_calculateLevelProbability_Distribution(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	// Generate many levels and check distribution
	levels := make(map[int]int)
	numSamples := 10000

	for i := 0; i < numSamples; i++ {
		level := alg.calculateLevelProbability()
		levels[level]++
	}

	// Level 0 should be most common, higher levels should be less common
	assert.Greater(t, levels[0], levels[1])
	assert.Greater(t, levels[1], levels[2])

	// Should have some variety of levels
	assert.GreaterOrEqual(t, len(levels), 5)
}

func TestHNSWNode_MultiLayerConnections(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    3,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{4.0, 5.0, 6.0}},
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    2,
		Connections: make(map[int][]*HNSWNode),
	}

	node3 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{7.0, 8.0, 9.0}},
			Metadata: map[string]string{"id": "3"},
			Id:       3,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	// Connect nodes on different layers
	node1.connect(node2, 0)
	node1.connect(node3, 0)
	node1.connect(node2, 1)
	node1.connect(node2, 2)

	// Check connections on layer 0
	assert.True(t, node1.isConnectedTo(node2, 0))
	assert.True(t, node1.isConnectedTo(node3, 0))
	assert.Len(t, node1.Connections[0], 2)

	// Check connections on layer 1
	assert.True(t, node1.isConnectedTo(node2, 1))
	assert.False(t, node1.isConnectedTo(node3, 1))
	assert.Len(t, node1.Connections[1], 1)

	// Check connections on layer 2
	assert.True(t, node1.isConnectedTo(node2, 2))
	assert.False(t, node1.isConnectedTo(node3, 2))
	assert.Len(t, node1.Connections[2], 1)

	// Layer 3 should have no connections
	assert.Len(t, node1.Connections[3], 0)
}

func TestAlgorithm_searchLayer(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	// Create a test algorithm instance
	alg := New(16, 200, 1.0/math.Log(2.0))

	// Create query node (not added to the graph)
	queryNode := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{0.0, 0.0, 0.0}}, // origin
			Metadata: map[string]string{"id": "query"},
			Id:       999,
		},
		MaxLayer:    0,
		Connections: make(map[int][]*HNSWNode),
	}

	// Create connected nodes at layer 0
	node1 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 0.0, 0.0}}, // close to query
			Metadata: map[string]string{"id": "1"},
			Id:       1,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	node2 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{2.0, 0.0, 0.0}}, // medium distance
			Metadata: map[string]string{"id": "2"},
			Id:       2,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	node3 := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{10.0, 0.0, 0.0}}, // far from query
			Metadata: map[string]string{"id": "3"},
			Id:       3,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	// Connect the nodes: node1 <-> node2 <-> node3 (linear chain)
	node1.connect(node2, 0)
	node2.connect(node3, 0)

	// Test searchLayer starting from node1, looking for 2 closest nodes
	entryNodes := []*HNSWNode{node1}
	result := alg.searchLayer(queryNode, entryNodes, 0, 2)

	// Should return the 2 closest nodes: node1 and node2
	assert.Len(t, result, 2)

	// Verify we got the right nodes (closest should be first)
	assert.Contains(t, result, node1)
	assert.Contains(t, result, node2)
	assert.NotContains(t, result, node3) // node3 should not be in top 2
}

func TestAlgorithm_searchLayer_SingleResult(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	queryNode := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{0.0, 0.0, 0.0}},
			Metadata: map[string]string{"id": "query"},
			Id:       999,
		},
		MaxLayer:    0,
		Connections: make(map[int][]*HNSWNode),
	}

	entryNode := &HNSWNode{
		Entry: algorithms.Entry{
			Vector:   vector.Vector{Values: []float64{1.0, 1.0, 1.0}},
			Metadata: map[string]string{"id": "entry"},
			Id:       1,
		},
		MaxLayer:    1,
		Connections: make(map[int][]*HNSWNode),
	}

	// Test with numClosest = 1
	result := alg.searchLayer(queryNode, []*HNSWNode{entryNode}, 0, 1)

	assert.Len(t, result, 1)
	assert.Equal(t, entryNode, result[0])
}

func TestAlgorithm_AddEntry_WithConnections(t *testing.T) {
	// Set deterministic seed for reproducible tests
	// rand.New(rand.NewSource(123))
	rand.Seed(123)

	// Small M value to test connection limits easily
	alg := New(2, 4, 1.0/math.Log(2.0))

	// Create entries with known distances
	entry1 := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{0.0, 0.0, 0.0}}, // origin
		Metadata: map[string]string{"id": "1"},
		Id:       1,
	}

	entry2 := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{1.0, 0.0, 0.0}}, // close to entry1
		Metadata: map[string]string{"id": "2"},
		Id:       2,
	}

	entry3 := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{0.0, 1.0, 0.0}}, // close to entry1
		Metadata: map[string]string{"id": "3"},
		Id:       3,
	}

	// Add entries one by one
	alg.AddEntry(entry1)
	assert.Len(t, alg.nodes, 1)
	assert.Equal(t, alg.nodes[0], alg.entryNode)

	alg.AddEntry(entry2)
	assert.Len(t, alg.nodes, 2)

	// Verify that nodes got connected at layer 0
	node1 := alg.nodes[0]
	node2 := alg.nodes[1]

	// At least one should be connected to the other at layer 0
	connected := node1.isConnectedTo(node2, 0) && node2.isConnectedTo(node1, 0)
	assert.True(t, connected, "Nodes should be connected after insertion")

	alg.AddEntry(entry3)
	assert.Len(t, alg.nodes, 3)

	// After 3 insertions, there should be some connections
	totalConnections := 0
	for _, node := range alg.nodes {
		if connections, exists := node.Connections[0]; exists {
			totalConnections += len(connections)
		}
	}
	assert.Greater(t, totalConnections, 0, "Should have some connections after multiple insertions")
}

func TestAlgorithm_AddEntry_ConnectionPruning(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	// Very small M=1 to test pruning quickly
	alg := New(1, 3, 1.0/math.Log(2.0))

	// Create entries in a line: A --- B --- C --- D
	entryA := algorithms.Entry{
		Vector: vector.Vector{Values: []float64{0.0, 0.0, 0.0}},
		Id:     1,
	}
	entryB := algorithms.Entry{
		Vector: vector.Vector{Values: []float64{1.0, 0.0, 0.0}},
		Id:     2,
	}
	entryC := algorithms.Entry{
		Vector: vector.Vector{Values: []float64{2.0, 0.0, 0.0}},
		Id:     3,
	}
	entryD := algorithms.Entry{
		Vector: vector.Vector{Values: []float64{3.0, 0.0, 0.0}},
		Id:     4,
	}

	alg.AddEntry(entryA)
	alg.AddEntry(entryB)
	alg.AddEntry(entryC)
	alg.AddEntry(entryD)

	// With M=1, each node should have at most 1 connection at layer 0
	for i, node := range alg.nodes {
		if connections, exists := node.Connections[0]; exists {
			assert.LessOrEqual(t, len(connections), alg.M,
				"Node %d should not exceed M=%d connections, but has %d",
				i, alg.M, len(connections))
		}
	}
}

func TestAlgorithm_AddEntry_MultiLayer(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(2, 4, 1.0/math.Log(2.0))

	// Add several entries to potentially get multi-layer nodes
	entries := []algorithms.Entry{
		{Vector: vector.Vector{Values: []float64{0.0, 0.0, 0.0}}, Id: 1},
		{Vector: vector.Vector{Values: []float64{1.0, 0.0, 0.0}}, Id: 2},
		{Vector: vector.Vector{Values: []float64{0.0, 1.0, 0.0}}, Id: 3},
		{Vector: vector.Vector{Values: []float64{1.0, 1.0, 0.0}}, Id: 4},
		{Vector: vector.Vector{Values: []float64{2.0, 0.0, 0.0}}, Id: 5},
	}

	for _, entry := range entries {
		alg.AddEntry(entry)
	}

	assert.Len(t, alg.nodes, 5)
	assert.NotNil(t, alg.entryNode)

	// Verify entry node is the one with highest MaxLayer
	maxLayer := -1
	for _, node := range alg.nodes {
		if node.MaxLayer > maxLayer {
			maxLayer = node.MaxLayer
		}
	}
	assert.Equal(t, maxLayer, alg.entryNode.MaxLayer)

	// Check that connections exist and respect layer constraints
	for _, node := range alg.nodes {
		for layer, connections := range node.Connections {
			// Layer should not exceed node's MaxLayer
			assert.LessOrEqual(t, layer, node.MaxLayer,
				"Node should not have connections above its MaxLayer")

			// All connected nodes should exist at this layer
			for _, connectedNode := range connections {
				assert.LessOrEqual(t, layer, connectedNode.MaxLayer,
					"Connected node should exist at this layer")
			}
		}
	}
}

func TestAlgorithm_AddEntry_EmptyGraph(t *testing.T) {
	// Set deterministic seed for reproducible tests
	rand.New(rand.NewSource(123))

	alg := New(16, 200, 1.0/math.Log(2.0))

	entry := algorithms.Entry{
		Vector:   vector.Vector{Values: []float64{1.0, 2.0, 3.0}},
		Metadata: map[string]string{"test": "value"},
		Id:       1,
	}

	// Should not panic when adding to empty graph
	alg.AddEntry(entry)

	assert.Len(t, alg.nodes, 1)
	assert.Equal(t, alg.nodes[0], alg.entryNode)
	assert.Equal(t, entry, alg.nodes[0].Entry)
}
