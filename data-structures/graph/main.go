package main

import (
	"fmt"
)

type Graph struct {
	NumberOfVerticies int
	AdjacentList      map[int][]int
}

func (g *Graph) AddVertex(value int) {
	if g.AdjacentList == nil {
		g.AdjacentList = map[int][]int{}
	}

	if g.AdjacentList[value] == nil {
		g.AdjacentList[value] = []int{}
	}
	g.NumberOfVerticies += 1
}

func (g *Graph) AddEdge(vertex1 int, vertex2 int) {
	if !includes(g.AdjacentList[vertex1], vertex2) {
		g.AdjacentList[vertex1] = append(g.AdjacentList[vertex1], vertex2)
	}

	if !includes(g.AdjacentList[vertex2], vertex1) {
		g.AdjacentList[vertex2] = append(g.AdjacentList[vertex2], vertex1)
	}
}

func (g *Graph) ShowConnection() (output string) {
	for vert, adjList := range g.AdjacentList {
		output += fmt.Sprintf("%d-->%s\n", vert, join(adjList, ","))
	}
	return output
}

func join(arr []int, delim string) (output string) {
	for i, elm := range arr {
		if i == len(arr)-1 {
			output += fmt.Sprintf("%d", elm)
		} else {
			output += fmt.Sprintf("%d%s ", elm, delim)
		}
	}
	return output
}

func includes(arr []int, val int) bool {
	for _, elm := range arr {
		if elm == val {
			return true
		}
	}
	return false
}

func main() {
	var myGraph = Graph{}
	myGraph.AddVertex(0)
	myGraph.AddVertex(1)
	myGraph.AddVertex(2)
	myGraph.AddVertex(3)
	myGraph.AddVertex(4)
	myGraph.AddVertex(5)
	myGraph.AddVertex(6)
	myGraph.AddEdge(3, 1)
	myGraph.AddEdge(3, 4)
	myGraph.AddEdge(4, 2)
	myGraph.AddEdge(4, 5)
	myGraph.AddEdge(1, 2)
	myGraph.AddEdge(1, 0)
	myGraph.AddEdge(0, 2)
	myGraph.AddEdge(6, 5)
	fmt.Print(myGraph.ShowConnection())
}
