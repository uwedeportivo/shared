// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

type BFSVisitor interface {
	Visit(u, v Vertex, g Graph)
}

type bfsColor int

const (
	white bfsColor = iota
	gray
	black
)

func BFS(g Graph, s Vertex, visitors []BFSVisitor) {
	color := make([]bfsColor, g.NumVertices())
	queue := newQueue(1024)

	color[s.ID()] = gray
	queue.enqueue(s.ID())

	for ;queue.size() > 0; {
		uid, _ := queue.dequeue()
		u := g.Vertex(uid)
		for _, v := range g.Neighbors(u) {
			if color[v.ID()] == white {
				for _, visitor := range visitors {
					visitor.Visit(u, v, g)
				}
				color[v.ID()] = gray
				queue.enqueue(v.ID())
			}
		}
		color[u.ID()] = black
	}
}

type BFSTreeRecorder []int

func (r BFSTreeRecorder) Visit(u, v Vertex, g Graph) {
	r[v.ID()] = u.ID()
}

type BFSLevelRecorder []int

func (r BFSLevelRecorder) Visit(u, v Vertex, g Graph) {
	r[v.ID()] = r[u.ID()] + 1
}