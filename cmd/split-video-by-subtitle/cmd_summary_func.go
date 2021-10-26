package main

import (
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
)

func searchOne(db *model.DB, searchKey string, mode SearchMode, limit int) (result SearchResult, err error) {
	var blocks model.EngSubtitleBlockList
	var query string
	var values []interface{}
	if mode == SearchModeNormal {
		mode = "NATURAL LANGUAGE"
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE MATCH(content) AGAINST(? IN %s MODE) LIMIT ?", mode)
		values = append(values, searchKey, limit)
	} else if mode == SearchModeBool {
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE MATCH(content) AGAINST(? IN %s MODE) LIMIT ?", mode)
		values = append(values, searchKey, limit)
	} else if mode == SearchModeLike {
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE content LIKE ? LIMIT ?")
		values = append(values, fmt.Sprintf("%%%s%%", searchKey), limit)
	} else {
		log.Fatal("invalid --mode=" + mode)
	}
	rows, err := db.Raw(query, values...).Rows()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		item := model.EngSubtitleBlock{}
		if err := db.ScanRows(rows, &item); err != nil {
			log.Error(err)
			break
		}
		blocks = append(blocks, item)
	}
	_ = rows.Close()
	log.Info("search: ", searchKey, " result: ", len(blocks))

	movieMap := make(map[int64]model.EngMovie)
	subtitleMap := make(map[int64]model.EngSubtitle)
	subtitleIds := make([]int64, 0)
	for _, item := range blocks {
		subtitleIds = append(subtitleIds, item.SubtitleId)
	}
	var subtitles model.EngSubtitleList
	var movies model.EngMovieList
	var movieIds = make([]int64, 0)
	if len(subtitleIds) > 0 {
		if err := db.Where("id IN (?)", subtitleIds).Find(&subtitles).Error; err != nil {
			log.Fatal(err)
		}
		for _, item := range subtitles {
			subtitleMap[item.Id] = item
			movieIds = append(movieIds, item.MovieId)
		}
	}

	if len(movieIds) > 0 {
		if err := db.Where("id IN (?)", movieIds).Find(&movies).Error; err != nil {
			log.Fatal(err)
		}
		for _, item := range movies {
			movieMap[item.Id] = item
		}
	}
	for _, block := range blocks {
		item := ResultItem{
			Block: block,
		}
		if subtitle, b := subtitleMap[block.SubtitleId]; b {
			item.Subtitle = subtitle
		}
		if movie, b := movieMap[item.Subtitle.MovieId]; b {
			item.Movie = movie
		}
		result.List = append(result.List, item)
	}
	return
}
