package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// EdgeTuple contain the HeadVertex and Distance from the Vertex.ID to which they belong
type EdgeTuple struct {
	HeadertexID int
	HeadVertex  *Vertex
	Distance    float64
}

func (n EdgeTuple) String() string {
	return fmt.Sprintf("\nHeadVertex:\t%d\nDistance:\t%f\n", n.HeadertexID, n.Distance)
}

//A Vertex is a point on a graph, which contains []Edges to other Vertexes.
type Vertex struct {
	ID     int
	Edges  []EdgeTuple
	DGS    float64 // Dijkstra Greedy Score
	Index  int     // Index of item in heap
	Length float64 // Calculated length
}

func (v Vertex) String() string {
	return fmt.Sprintf("\nid:\t%d\nedges: %v\nDGS:\t%g\nIndex:\t%d\nLength:\t%g\n\n\n\n", v.ID, v.Edges, v.DGS, v.Index, v.Length)
}

// AddEdge appends an EdgeTuple to a Vertex's []Edges
func (v *Vertex) AddEdge(e EdgeTuple) {
	v.Edges = append(v.Edges, e)
}

// VertexMap tracks Vertexes using [Vertex.ID]*Vertex
var VertexMap = make(map[int]*Vertex)

// A VertexHeap implements heap.Interface and holds Vertexes.
type VertexHeap []*Vertex

func (vh VertexHeap) Len() int { return len(vh) }

func (vh VertexHeap) Less(i, j int) bool {
	return vh[i].DGS < vh[j].DGS
}

func (vh VertexHeap) Swap(i, j int) {
	vh[i], vh[j] = vh[j], vh[i]
	vh[i].Index = i
	vh[j].Index = j
}

// Push adds Vertexes to VertexHeaps
func (vh *VertexHeap) Push(x interface{}) {
	n := len(*vh)
	v := x.(*Vertex)
	v.Index = n
	*vh = append(*vh, v)
}

// Pop returns the Vertex with the lowest DGS and removes it from the heap
func (vh *VertexHeap) Pop() interface{} {
	old := *vh
	n := len(old)
	v := old[n-1]
	v.Index = -1 // for safety, identify it's no longer in heap
	*vh = old[0 : n-1]
	return v
}

// update modifies the DGS of a Vertex in the heap.
func (vh *VertexHeap) update(v *Vertex, DGS float64) {
	v.DGS = DGS
	heap.Fix(vh, v.Index)
}

var vh VertexHeap

// This example creates a VertexHeap with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {

	readFile(os.Args[1])

	makeVertexHeap()

	dijkstra(1)

	for _, v := range VertexMap {
		fmt.Printf("%d:\t%d\n", v.ID, int(v.Length))
	}
}

func readFile(filename string) {

	file, err := os.Open(filename) //should read in file named in CLI
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		thisLine := strings.Fields(scanner.Text())

		thisertexID, err := strconv.Atoi(thisLine[0])

		if err != nil {
			fmt.Printf("couldn't convert number: %v\n", err)
			return
		}

		w, ok := VertexMap[thisertexID]

		if !ok {
			w = &Vertex{thisertexID, []EdgeTuple{}, math.Inf(1), -1, -1}
			VertexMap[thisertexID] = w
		}

		for i := 1; i < len(thisLine); i++ {

			weightedEdge := strings.Split(thisLine[i], ",")

			edgeID, err := strconv.Atoi(weightedEdge[0])
			weightOfEdge, err := strconv.ParseFloat(weightedEdge[1], 64)

			if err != nil {
				fmt.Printf("couldn't convert number: %v\n", err)
				return
			}

			u, ok := VertexMap[edgeID]

			if !ok {
				u = &Vertex{edgeID, []EdgeTuple{}, math.Inf(1), -1, -1}
				VertexMap[edgeID] = u
			}

			w.AddEdge(EdgeTuple{edgeID, u, weightOfEdge})

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func makeVertexHeap() {

	vh = make(VertexHeap, len(VertexMap))

	i := 0

	for _, v := range VertexMap {
		v.Index = i
		vh[i] = v
		i++
	}

	heap.Init(&vh)
}

func dijkstra(id int) {

	workingVertex := VertexMap[id]
	workingVertex.DGS = 0
	workingVertex.Length = 0
	vh.update(workingVertex, workingVertex.DGS)

	for vh.Len() > 0 {
		for _, tuple := range workingVertex.Edges {
			v := tuple.HeadVertex
			TestDGS := workingVertex.Length + tuple.Distance
			if v.DGS > TestDGS {
				vh.update(v, TestDGS)
			}
		}

		workingVertex = heap.Pop(&vh).(*Vertex)

		workingVertex.Length = workingVertex.DGS

	}

}
