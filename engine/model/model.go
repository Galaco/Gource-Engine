package model

import (
	"github.com/galaco/Gource-Engine/engine/mesh"
)

// A collection of renderable primitives/submeshes
type Model struct {
	meshes   []mesh.IMesh
}

// Add a new primitive
func (model *Model) AddMesh(meshes ...mesh.IMesh) {
	model.meshes = append(model.meshes, meshes...)
}

// Get all primitives/submeshes
func (model *Model) GetMeshes() []mesh.IMesh {
	return model.meshes
}

func (model *Model) Reset() {
	model.meshes = []mesh.IMesh{}
}

func NewModel(meshes ...mesh.IMesh) *Model {
	return &Model{
		meshes:   meshes,
	}
}