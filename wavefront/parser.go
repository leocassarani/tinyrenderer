package wavefront

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func ParseModel(r io.Reader) (*Model, error) {
	model := &Model{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "f "):
			face := parseFace(line)
			model.Faces = append(model.Faces, face)
		case strings.HasPrefix(line, "v "):
			vertex := parseVertex(line)
			model.Vertices = append(model.Vertices, vertex)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return model, nil
}

type Model struct {
	Faces    []Face
	Vertices []Vertex
}

func (m Model) VertexAt(idx int) Vertex {
	return m.Vertices[idx-1]
}

func parseFace(line string) Face {
	parts := strings.Split(line, " ")[1:] // Skip the initial "f".
	indices := make([]int, len(parts))

	for i, part := range parts {
		idx := strings.Split(part, "/")[0]
		indices[i] = mustParseInt(idx)
	}

	return Face{indices}
}

type Face struct {
	Indices []int
}

func parseVertex(line string) Vertex {
	parts := strings.Split(line, " ")[1:] // Skip the initial "v".

	return Vertex{
		X: mustParseFloat(parts[0]),
		Y: mustParseFloat(parts[1]),
		Z: mustParseFloat(parts[2]),
	}
}

type Vertex struct {
	X, Y, Z float64
}

func mustParseInt(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

func mustParseFloat(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
