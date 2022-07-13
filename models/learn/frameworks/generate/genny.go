//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "KeyType=string,int ValueType=string,int"
package frameworks

import (
	"fmt"

	"github.com/cheekybits/genny/generic"
)

/* Use generic.Type to Compile */
type KeyType generic.Type
type ValueType generic.Type

type KeyTypeValueTypeMap struct {
	typedMap map[KeyType]ValueType
}

func (self *KeyTypeValueTypeMap) PrintType() {
	for key, value := range self.typedMap {
		fmt.Printf("Key:%v Value: %v\n", key, value)
	}
}
