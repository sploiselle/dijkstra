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

type NodeTuple struct {
	HeadVertexId int
	HeadVertex   *Vertex
	Distance     float64
}

func (n NodeTuple) String() string {
	return fmt.Sprintf("\nHeadVertex:\t%d\nDistance:\t%f\n", n.HeadVertexId, n.Distance)
}

type Vertex struct {
	Id     int
	Edges  []NodeTuple
	DGS    float64 // Dijkstra Greedy Score
	Index  int     // Index of item in heap
	Length float64 // Calculated length
}

func (v Vertex) String() string {
	return fmt.Sprintf("\nid:\t%d\nedges: %v\nDGS:\t%g\nIndex:\t%d\nLength:\t%g\n\n\n\n", v.Id, v.Edges, v.DGS, v.Index, v.Length)
}

func (v *Vertex) AddEdge(i NodeTuple) {
	v.Edges = append(v.Edges, i)
}

// Adjacency lists
//[Vertex.Id]*Vertex
var VertexMap = make(map[int]*Vertex)

//[Vertex.Id]distance
var DistanceMap = make(map[int]int)

// A VertexHeap implements heap.Interface and holds Items.
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

func (vh *VertexHeap) Push(x interface{}) {
	n := len(*vh)
	this_vertex := x.(*Vertex)
	this_vertex.Index = n
	*vh = append(*vh, this_vertex)
}

func (vh *VertexHeap) Pop() interface{} {
	old := *vh
	n := len(old)
	this_vertex := old[n-1]
	this_vertex.Index = -1 // for safety, identify it's no longer in heap
	*vh = old[0 : n-1]
	return this_vertex
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

	for _, v := range []int{7, 37, 59, 82, 99, 115, 133, 165, 188, 197} {
		fmt.Printf("%d,", int(VertexMap[v].Length))
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

		thisVertexId, err := strconv.Atoi(thisLine[0])

		if err != nil {
			fmt.Printf("couldn't convert number: %v\n", err)
			return
		}

		w, ok := VertexMap[thisVertexId]

		if !ok {
			w = &Vertex{thisVertexId, []NodeTuple{}, math.Inf(1), -1, -1}
			VertexMap[thisVertexId] = w
		}

		for i := 1; i < len(thisLine); i++ {

			weightedEdge := strings.Split(thisLine[i], ",")

			edgeId, err := strconv.Atoi(weightedEdge[0])
			weightOfEdge, err := strconv.ParseFloat(weightedEdge[1], 64)

			if err != nil {
				fmt.Printf("couldn't convert number: %v\n", err)
				return
			}

			u, ok := VertexMap[edgeId]

			if !ok {
				u = &Vertex{edgeId, []NodeTuple{}, math.Inf(1), -1, -1}
				VertexMap[edgeId] = u
			}

			w.AddEdge(NodeTuple{edgeId, u, weightOfEdge})

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
			test_DGS := workingVertex.Length + tuple.Distance
			if v.DGS > test_DGS {
				vh.update(v, test_DGS)
			}
		}

		workingVertex = heap.Pop(&vh).(*Vertex)

		workingVertex.Length = workingVertex.DGS

	}

}
