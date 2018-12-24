package common

import (
	entity2 "github.com/galaco/Gource-Engine/core/entity"
	"github.com/galaco/Gource-Engine/game/entity"
)

type PropDynamic struct {
	entity2.Base
	entity.PropBase
}

func (entity *PropDynamic) New() entity2.IEntity {
	return &PropDynamic{}
}

func (entity PropDynamic) Classname() string {
	return "prop_dynamic"
}