// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"strconv"
)

type IntVertex int

func (i IntVertex) ID() int {
	return int(i)
}

func (i IntVertex) String() string {
	return strconv.Itoa(int(i))
}

type IntAdjacencyListGraph [][]Vertex

func (g IntAdjacencyListGraph) NumVertices() int {
	return len(g)
}

func (g IntAdjacencyListGraph) Neighbors(v Vertex) []Vertex {
	return g[v.ID()]
}

func (g IntAdjacencyListGraph) Vertex(id int) Vertex {
	return IntVertex(id)
}

func (g IntAdjacencyListGraph) Add(u, v int) {
	if g[u] == nil {
		g[u] = make([]Vertex, 1, 16)
		g[u][0] = IntVertex(v)
	} else {
		g[u] = append(g[u], IntVertex(v))
	}
}


