package entity

import (
	di "github.com/docker/docker/api/types/image"
)

// type Image struct {
// 	ID        string    `json:"id"`
// 	Name      string    `json:"name"`
// 	Tags      []string  `json:"tags"`
// 	Size      int64     `json:"size"`
// 	CreatedAt time.Time `json:"created_at"`
// }

type Image = di.Summary

type ImageHistory = di.HistoryResponseItem
