// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"io"
	"bufio"
	"errors"
	"fmt"
	"strings"
)

func ReadSNAP(r io.Reader) (Graph, error) {
	lr := bufio.NewReader(r)

	inHeader := true
	scannedNumVertices := false

	var numVertices, numEdges int

	var g IntAdjacencyListGraph

	var currentNode IntVertex = -1

	for {
		line, isPrefix, err := lr.ReadLine()

		if err == io.EOF {
			break;
		}

		if err != nil {
			return nil, err
		}

		if isPrefix {
			return nil, errors.New("line too big")
		}

		l := string(line)

		if inHeader {
			if strings.HasPrefix(l, "#") {
				if !scannedNumVertices && strings.HasPrefix(l, "# Nodes: ") {
					_, err = fmt.Sscanf(l, "# Nodes: %d Edges: %d", &numVertices, &numEdges)
					if err != nil {
						return nil, err
					}
					if numVertices <= 0 {
						return nil, errors.New("Invalid number of vertices")
					}
					scannedNumVertices = true
					g = make([][]Vertex, numVertices)
				}
			} else {
				inHeader = false
			}
		} else {
			if numVertices == 0 {
				return nil, errors.New("Unknown number of vertices")
			}
			var v1, v2 IntVertex
			_, err = fmt.Sscanf(l, "%d %d", &v1, &v2)
			if err != nil {
				return nil, err
			}

			if v1 < 0 || int(v1) >= numVertices {
				return nil, errors.New("Invalid vertex")
			}

			if v2 < 0 || int(v2) >= numVertices {
				return nil, errors.New("Invalid vertex")
			}

			if v1 != currentNode {
				g[v1] = make([]Vertex, 1, 16)
				g[v1][0] = v2
				currentNode = v1
			} else {
				g[v1] = append(g[v1], v2)
			}
		}
	}
	return g, nil
}
