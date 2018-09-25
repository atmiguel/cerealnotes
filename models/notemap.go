package models

import (
	"encoding/json"
	"fmt"
)

type NoteMap map[NoteId]*Note

func (noteMap NoteMap) ToJson() ([]byte, error) {
	// json doesn't support int indexed maps
	notesByIdString := make(map[string]Note, len(noteMap))

	for id, note := range noteMap {
		notesByIdString[fmt.Sprint(id)] = *note
	}

	return json.Marshal(notesByIdString)
}
