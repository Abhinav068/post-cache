package main

import "sync"

type Post struct {
	sync.Mutex  //! @anurag how to stop saving it in DB
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Likes      int    `json:"likes"`
}
