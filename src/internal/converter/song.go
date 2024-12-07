package converter

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"song-library-api/src/internal/model"
)

func ToRecordFromSong(song model.Song) goqu.Record {
	record := goqu.Record{}

	if song.ID != uuid.Nil {
		record["id"] = song.ID
	}
	if song.GroupID != uuid.Nil {
		record["group_id"] = song.GroupID
	}
	if song.Song != "" {
		record["song"] = song.Song
	}
	if song.Text != "" {
		record["text"] = song.Text
	}
	if song.Link != "" {
		record["link"] = song.Link
	}
	if !song.ReleaseDate.IsZero() {
		record["release_date"] = song.ReleaseDate
	}
	if !song.CreatedAt.IsZero() {
		record["created_at"] = song.CreatedAt
	}
	if !song.UpdatedAt.IsZero() {
		record["updated_at"] = song.UpdatedAt
	}

	return record
}

func ToViewFromSong(song model.Song) model.SongView {
	return model.SongView{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		Text:        song.Text,
		Link:        song.Link,
		ReleaseDate: song.ReleaseDate,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
	}
}

func ToViewsFromSong(songs []model.Song) []model.SongView {
	views := make([]model.SongView, 0, len(songs))
	for _, song := range songs {
		views = append(views, ToViewFromSong(song))
	}
	return views
}
