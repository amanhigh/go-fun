//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "KeyType=string,int ValueType=string,int"
package generate

import (
	"github.com/cheekybits/genny/generic"
	"github.com/rs/zerolog/log"
)

/* Use generic.Type to Compile */
type KeyType generic.Type
type ValueType generic.Type

type KeyTypeValueTypeMap struct {
	typedMap map[KeyType]ValueType
}

func (self *KeyTypeValueTypeMap) PrintType() {
	for key, value := range self.typedMap {
		log.Info().Any("Key", key).Any("Value", value).Msg("Map")
	}
}
