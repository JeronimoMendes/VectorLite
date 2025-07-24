package hnsw

import (
	"VectorLite/internal/algorithms"
	"math"
	"math/rand"
	"slices"
	"sort"
)

type Algorithm struct {
	nodes          []*HNSWNode
	entryNode      *HNSWNode
	M              int
	efConstruction int
	mL             float64
}

func New(M int, efConstruction int, mL float64) *Algorithm {
	return &Algorithm{
		nodes:          []*HNSWNode{},
		entryNode:      nil,
		M:              M,
		efConstruction: efConstruction,
		mL:             mL,
	}
}

/*
A HNSW node can be present in multiple layers of the graph structure
This means it has a max layer and it can be connected to other nodes at each of it's levels.

Let's take the example of node nA:

	Layer 1: nodeC  <-->  nodeA  <-->  nodeB
	Layer 2: nodeD  <-->  nodeA  <-->  nodeB
	Layer 3: nodeA  <-->  nodeB

	nodeA max layer is 3, and it's connected to multiple other nodes throughout it's layers
*/
type HNSWNode struct {
	Entry       algorithms.Entry
	MaxLayer    int
	Connections map[int][]*HNSWNode
}

type CandidateNode struct {
	Node  *HNSWNode
	Score float64
}

func (n *HNSWNode) isConnectedTo(otherNode *HNSWNode, layer int) bool {
	return slices.Contains(n.Connections[layer], otherNode)
}

func (n *HNSWNode) connect(otherNode *HNSWNode, layer int) {
	if n == otherNode {
		panic("Can't connect HNSWNode to itself")
	}

	if n.isConnectedTo(otherNode, layer) {
		return
	}

	n.Connections[layer] = append(n.Connections[layer], otherNode)
	if !otherNode.isConnectedTo(n, layer) {
		otherNode.connect(n, layer)
	}
}

func (n *HNSWNode) disconnect(otherNode *HNSWNode, layer int) {
	n.Connections[layer] = slices.DeleteFunc(n.Connections[layer], func(i *HNSWNode) bool {
		return i == otherNode
	})
	if otherNode.isConnectedTo(n, layer) {
		otherNode.disconnect(n, layer)
	}
}

func (n *HNSWNode) getScore(otherNode *HNSWNode) float64 {
	return n.Entry.Vector.Cosine_similarity(&otherNode.Entry.Vector)
}

func (a *Algorithm) AddEntry(entry algorithms.Entry) {
	newNode := &HNSWNode{
		Entry:       entry,
		MaxLayer:    a.calculateLevelProbability(),
		Connections: make(map[int][]*HNSWNode),
	}
	a.nodes = append(a.nodes, newNode)

	// new nodes becomes king of the hill if we had no previous king or
	// new king is better than previous
	if a.entryNode == nil {
		a.entryNode = newNode
		return
	}

	currentLayer := a.entryNode.MaxLayer
	entryPoints := []*HNSWNode{a.entryNode}

	// till we reach newNode's max layer we do a rapid descent
	for currentLayer > newNode.MaxLayer {
		entryPoints = a.searchLayer(newNode, entryPoints, currentLayer, 1)
		currentLayer--
	}

	for currentLayer >= 0 {
		entryPoints = a.searchLayer(newNode, entryPoints, currentLayer, a.efConstruction)

		// now we need to connect to this new entryPoints
		var connectTo []*HNSWNode

		// we can't connect to more than the max connections per node
		if len(entryPoints) > a.M {
			connectTo = entryPoints[:a.M]
		} else {
			connectTo = entryPoints
		}
		for _, connectNode := range connectTo {
			// this may or may not create a connection, depends how strong the score is for our new node
			a.createConnection(newNode, connectNode, currentLayer)
		}

		currentLayer--
	}

	if a.entryNode.MaxLayer < newNode.MaxLayer {
		a.entryNode = newNode
	}
}

func (a *Algorithm) createConnection(newNode *HNSWNode, existentNode *HNSWNode, layer int) {
	if len(existentNode.Connections[layer]) == a.M {
		// we need to check the lowest connection for this node
		// if it's lowest than the newNode, we pop it and connect to the newNode
		sort.Slice(existentNode.Connections[layer], func(i int, j int) bool {
			return existentNode.getScore(existentNode.Connections[layer][i]) > existentNode.getScore(existentNode.Connections[layer][j])
		})

		weakestConnection := existentNode.Connections[layer][a.M-1]
		if existentNode.getScore(newNode) > existentNode.getScore(weakestConnection) {
			existentNode.disconnect(weakestConnection, layer)
			existentNode.connect(newNode, layer)
		}
	} else {
		existentNode.connect(newNode, layer)
	}
}

func (a *Algorithm) calculateLevelProbability() int {
	level := int(math.Floor(-math.Log(rand.Float64() * a.mL)))
	if level < 0 {
		level = 0
	}
	return level
}

/*
*   This method searches for the closest nodes on a layer
*
*   node
*     node we are executing the search for
*   entryNodes
*     starting nodes for the search
*   layer
*     Layer we are alt
*   numClosest
*     number of close nodes to return
 */
func (a *Algorithm) searchLayer(node *HNSWNode, entryNodes []*HNSWNode, layer int, numClosest int) []*HNSWNode {
	candidates := []*CandidateNode{}
	visited := []*HNSWNode{}

	for _, entryNode := range entryNodes {
		candidates = append(candidates, &CandidateNode{
			Node:  entryNode,
			Score: node.getScore(entryNode),
		})
	}

	sort.Slice(candidates, func(i int, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	for len(candidates) > 0 {
		var current *CandidateNode

		for _, candidate := range candidates {
			if !slices.Contains(visited, candidate.Node) {
				current = candidate
				break
			}
		}

		if current == nil {
			break
		}

		// current candidate sucks big time and we just quit finding a better one
		if len(candidates) >= numClosest && current.Score < candidates[numClosest-1].Score {
			break
		}

		visited = append(visited, current.Node)

		for _, conn := range current.Node.Connections[layer] {
			// we don't want to check visited nodes
			if !slices.Contains(visited, conn) {
				connScore := node.getScore(conn)

				connCandidateNode := &CandidateNode{
					Node:  conn,
					Score: connScore,
				}
				updated := false
				if len(candidates) == numClosest {
					worst := candidates[len(candidates)-1]
					if connScore > worst.Score {
						candidates = candidates[:numClosest-1]
						candidates = append(candidates, connCandidateNode)
						updated = true
					}
				} else {
					candidates = append(candidates, connCandidateNode)
					updated = true
				}

				// TODO: keep slice sorted on inserts so we don't need to sorte everytimej
				if updated {
					sort.Slice(candidates, func(i int, j int) bool {
						return candidates[i].Score > candidates[j].Score
					})
				}
			}
		}
	}

	closest := []*HNSWNode{}
	for _, c := range candidates {
		closest = append(closest, c.Node)
	}

	return closest
}
