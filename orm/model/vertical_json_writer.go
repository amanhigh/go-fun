//This is Aman's Generated File
//Request you not to mess with it :)

package model

import (
	"encoding/json"
	"io"
)

func (obj Vertical) WriteTo(writer io.Writer) (int64, error) {
	data, err := json.Marshal(&obj)
	if err != nil {
		return 0, err
	}
	length, err := writer.Write(data)
	return int64(length), err
}