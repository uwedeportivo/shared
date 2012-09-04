// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

type Graph interface {
	NumVertices() int
	Neighbors(v Vertex) []Vertex
	Vertex(id int) Vertex
}

type Vertex interface {
    String() string
    ID() int
}

