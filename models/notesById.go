package models

import (
	"encoding/json"
	"fmt"
)

type NotesById map[NoteId]*Note

func (notesById NotesById) ToJson() ([]byte, error) {
	// json doesn't support int indexed maps
	notesByIdString := make(map[string]Note, len(notesById))

	for id, note := range notesById {
		notesByIdString[fmt.Sprint(id)] = *note
	}

	return json.Marshal(notesByIdString)
}
