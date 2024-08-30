package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/o98k-ok/lazy/v2/alfred"
)

var (
	HOST = "bing.com"
)

type BingImg struct {
	Startdate     string `json:"startdate"`
	Fullstartdate string `json:"fullstartdate"`
	Enddate       string `json:"enddate"`
	URL           string `json:"url"`
	Urlbase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Quiz          string `json:"quiz"`
	Wp            bool   `json:"wp"`
	Hsh           string `json:"hsh"`
	Drk           int    `json:"drk"`
	Top           int    `json:"top"`
	Bot           int    `json:"bot"`
	Hs            []any  `json:"hs"`
}

type BingResponse struct {
	Images []BingImg `json:"images"`
}

func main() {
	var dir string
	var count int
	flag.StringVar(&dir, "dir", "", "wallpaper dir")
	flag.IntVar(&count, "count", 9, "wallpaper count")
	flag.Parse()
	if len(dir) == 0 {
		pwd, _ := os.Getwd()
		dir = path.Join(pwd, "wallpaper")
	}

	{
		os.MkdirAll(dir, 0755)
	}

	{
		url := fmt.Sprintf("https://%s/HPImageArchive.aspx?format=js&idx=0&n=%d&mkt=en-US", HOST, count)
		resp, err := netutil.HttpGet(url)
		if err != nil {
			alfred.Log("get %s err %v", url, err)
			return
		}

		var bing BingResponse
		if err = json.NewDecoder(resp.Body).Decode(&bing); err != nil {
			alfred.Log("decode %s response err %v", url, err)
			return
		}

		items := alfred.NewItems()
		items.Items = make([]*alfred.Item, len(bing.Images))
		groups := sync.WaitGroup{}
		for i, v := range bing.Images {
			groups.Add(1)
			go func(idx int, val BingImg) {
				defer groups.Done()

				filename := path.Join(dir, fmt.Sprintf("%s_%s.jpg", val.Startdate, val.Hsh))
				item := alfred.NewItem(val.Title, "", filename)
				item.Icon = &alfred.Icon{
					Path: filename,
				}
				items.Items[idx] = item

				{
					_, err := os.Stat(filename)
					if err == nil {
						return
					}
				}

				err = netutil.DownloadFile(filename, fmt.Sprintf("https://%s/%s", HOST, val.URL))
				if err != nil {
					alfred.Log("download %s err %v", val.URL, err)
				}
			}(i, v)
		}
		groups.Wait()
		items.Show()
	}
}
