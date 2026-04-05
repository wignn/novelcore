package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *GraphQLServer
}

func (r *accountResolver) ReadingList(c context.Context, o *Account) ([]*ReadingListEntry, error) {
	c, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	entries, err := r.server.readinglistClient.GetReadingList(c, o.ID, "", 0, 100)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result []*ReadingListEntry
	for _, e := range entries {
		ratingInt := int(e.Rating)
		result = append(result, &ReadingListEntry{
			ID: e.ID, NovelID: e.NovelID, Status: e.Status,
			CurrentChapter: e.CurrentChapter, Rating: &ratingInt,
			Notes: &e.Notes, IsFavorite: e.IsFavorite,
			CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
		})
	}
	return result, nil
}
