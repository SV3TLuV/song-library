package converter

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"song-library-api/src/cmd/api/internal/model"
)

func ToRecordFromGroup(group model.Group) goqu.Record {
	record := goqu.Record{}

	if group.ID != uuid.Nil {
		record["id"] = group.ID
	}
	if group.Name != "" {
		record["name"] = group.Name
	}
	if !group.CreatedAt.IsZero() {
		record["created_at"] = group.CreatedAt
	}
	if !group.UpdatedAt.IsZero() {
		record["updated_at"] = group.UpdatedAt
	}

	return record
}
