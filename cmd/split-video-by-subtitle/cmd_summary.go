package main

import (
	"encoding/csv"
	"fmt"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type summaryCmd struct {
	Flag struct {
		In         string
		OutDir     string
		SearchMode string
		Limit      int
	}
	db *model.DB
}

func (t *summaryCmd) loadSearchKeys() ([]string, error) {
	var ret []string
	data, err := ioutil.ReadFile(t.Flag.In)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "load search key data failed")
	}
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		ret = append(ret, strings.Trim(l, "\n\t\r ,.?"))
	}
	return ret, nil
}

func (t *summaryCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "summary",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			t.db, err = model.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer func(db *model.DB) {
				err := db.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(t.db)
			searchKeys, err := t.loadSearchKeys()
			if err != nil {
				log.Fatal(err)
			}
			if err := lib.TouchDir(t.Flag.OutDir, 0755); err != nil {
				log.Fatal(err)
			}
			out := path.Join(t.Flag.OutDir, time.Now().Format("2006-01-02_150405")+".csv")
			f, err := os.Create(out)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			writer := csv.NewWriter(f)
			if err := writer.Write([]string{
				"index", "search_key", "movie_id", "movie_name", "subtitle_name", "block_id", "content", "start_time", "end_time",
			}); err != nil {
				log.Fatal(err)
			}
			count := 0
			for _, k := range searchKeys {
				result, err := searchOne(t.db, k, SearchMode(t.Flag.SearchMode), t.Flag.Limit)
				if err != nil {
					log.Fatal(err)
				}
				for _, item := range result.List {
					row := []string{
						fmt.Sprintf("%v", count),
						k,
						fmt.Sprintf("%v", item.Movie.Id),
						item.Movie.Name,
						item.Subtitle.SubName,
						fmt.Sprintf("%v", item.Block.BlockId),
						item.Block.Content,
						item.Block.StartTime,
						item.Block.EndTime,
					}
					if err := writer.Write(row); err != nil {
						log.Fatal(err)
					}
					count += 1
				}
			}
		},
	}
	t.initFlag(cmd)
	return cmd
}

func (t *summaryCmd) initFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&t.Flag.In, "in", "tmp/in.list", "search key list")
	cmd.PersistentFlags().StringVar(&t.Flag.OutDir, "out", "tmp/summary-result", "summary result")
	cmd.PersistentFlags().StringVar(&t.Flag.SearchMode, "mode", "like", "search mode")
	cmd.PersistentFlags().IntVar(&t.Flag.Limit, "limit", 10, "limit")
}

func init() {
	rootCmd.AddCommand((&summaryCmd{}).Build())
}
