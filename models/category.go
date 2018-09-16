package models

func (db *DB) StoreNewNoteCategoryRelationship(
	noteId NoteId,
	category Category,
) error {
	sqlQuery := `
		INSERT INTO note_to_category_relationship (note_id, category)
		VALUES ($1, $2)`

	rows, err := db.Query(sqlQuery, int64(noteId), category.String())
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}
