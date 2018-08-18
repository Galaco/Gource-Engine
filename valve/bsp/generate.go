package bsp

import (
	"math"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/texinfo"
)

func GenerateFacesFromBSP(file *bsp.Bsp) ([]float32, [][]uint16, []texinfo.TexInfo) {
	var verts []float32
	var expFaces [][]uint16
	var expTexInfos []texinfo.TexInfo

	fl := *file.GetLump(bsp.LUMP_FACES).GetContents()
	faces := (fl).(lumps.Face).GetData().(*[]face.Face)

	vs := *file.GetLump(bsp.LUMP_VERTEXES).GetContents()
	vertexes := (vs).(lumps.Vertex).GetData().(*[]mgl32.Vec3)

	sf := *file.GetLump(bsp.LUMP_SURFEDGES).GetContents()
	surfEdges := (sf).(lumps.Surfedge).GetData().(*[]int32)

	ed := *file.GetLump(bsp.LUMP_EDGES).GetContents()
	edges := (ed).(lumps.Edge).GetData().(*[][2]uint16)

	ti := *file.GetLump(bsp.LUMP_TEXINFO).GetContents()
	texInfos := ti.(lumps.TexInfo).GetData().(*[]texinfo.TexInfo)


	for _,v := range *vertexes {
		verts = append(verts, v.X(), v.Y(), v.Z())
	}

	// NOTE: We are converting from face to triangles here.
	for _,f := range *faces {
		expF := make([]uint16, 0)
		//// Just so we render triangles

		// All surfedges associated with this face
		// surfEdges are basically indices into the edges lump
		faceSurfEdges := (*surfEdges)[f.FirstEdge:(f.FirstEdge+int32(f.NumEdges))]
		rootIndex := uint16(0)
		for idx,surfEdge := range faceSurfEdges {
			edge := (*edges)[int(math.Abs(float64(surfEdge)))]
			e1 := 0
			e2 := 1
			if surfEdge < 0 {
				e1 = 1
				e2 = 0
			}
			//Capture root indice
			if idx == 0 {
				rootIndex = edge[e1]
			} else {
				// Just create a triangle for every edge now
				expF = append(expF, rootIndex, edge[e1], edge[e2])
			}
		}

		expFaces = append(expFaces, expF)
		expTexInfos = append(expTexInfos, (*texInfos)[f.TexInfo])
	}

	return verts, expFaces, expTexInfos
}

func TexCoordsForFaceFromTexInfo(vertexes []float32, tx *texinfo.TexInfo) []float32{
	uvs := make([]float32, (len(vertexes) / 3) * 2)

	for idx := 0; idx < len(vertexes); idx += 3 {
		//u = tv0,0 * x + tv0,1 * y + tv0,2 * z + tv0,3
		u := (tx.TextureVecsTexelsPerWorldUnits[0][0] * vertexes[idx]) +
			(tx.TextureVecsTexelsPerWorldUnits[0][1] * vertexes[idx+1]) +
			(tx.TextureVecsTexelsPerWorldUnits[0][2] * vertexes[idx+2]) +
			tx.TextureVecsTexelsPerWorldUnits[0][3]

		//v = tv1,0 * x + tv1,1 * y + tv1,2 * z + tv1,3
		v := (tx.TextureVecsTexelsPerWorldUnits[1][0] * vertexes[idx]) +
			(tx.TextureVecsTexelsPerWorldUnits[1][1] * vertexes[idx+1]) +
			(tx.TextureVecsTexelsPerWorldUnits[1][2] * vertexes[idx+2]) +
			tx.TextureVecsTexelsPerWorldUnits[1][3]

		uvs = append(uvs, u, v)
	}

	return uvs
}