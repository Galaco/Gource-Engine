package loader

import (
	material2 "github.com/galaco/lambda-core/loader/material"
	"github.com/galaco/lambda-core/material"
	"github.com/galaco/lambda-core/mesh/primitive"
	"github.com/galaco/lambda-core/model"
	"github.com/galaco/lambda-core/texture"
	"github.com/golang-source-engine/filesystem"
)

const skyboxRootDir = "skybox/"

// LoadSky loads the skymaterial cubemap.
// The materialname is normally obtained from the worldspawn entity
func LoadSky(materialName string, fs *filesystem.FileSystem) *model.Model {
	sky := model.NewModel(materialName)

	mats := make([]material.IMaterial, 6)

	mats[0] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"up.vmt", fs)
	mats[1] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"dn.vmt", fs)
	mats[2] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"lf.vmt", fs)
	mats[3] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"rt.vmt", fs)
	mats[4] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"ft.vmt", fs)
	mats[5] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"bk.vmt", fs)

	texs := make([]texture.ITexture, 6)
	for i := 0; i < 6; i++ {
		texs[i] = mats[i].(*material.Material).Textures.Albedo
	}

	sky.AddMesh(primitive.NewCube())

	sky.Meshes()[0].SetMaterial(texture.NewCubemap(texs))

	return sky
}
