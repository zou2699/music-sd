package qq

import (
	"encoding/json"
	"fmt"
	"github.com/zou2699/musicSD/models"
	"github.com/zou2699/musicSD/pkg/common"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Search(keyword string) (musicList []models.Music) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://c.y.qq.com/soso/fcgi-bin/search_for_qq_cp", nil)
	common.AddHeader(req)
	req.Header.Set("referer", "http://m.y.qq.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")

	q := req.URL.Query()
	q.Add("w", keyword)
	q.Add("format", "json")
	q.Add("p", "1")
	// count
	q.Add("n", "8")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}
	var respJson models.RespQQ
	err = json.Unmarshal(content, &respJson)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Printf("%s\n", content)
	if respJson.Code != 0 {
		log.Panic("code isn't 0", respJson.Code)
	}

	var music models.Music
	for _, song := range respJson.Data.Song.List {
		// singer
		var singers []string
		for _, singer := range song.Singer {
			singers = append(singers, singer.Name)
		}

		size := song.Size128
		if song.Size320 != 0 {
			size = song.Size320
		}
		mSize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1048576), 64)

		music.Title = song.Songname
		music.ID = song.Songid
		music.MID = song.Songmid
		music.Duration = ""
		music.Singer = strings.Join(singers, ",")
		music.Album = song.Albumname
		music.Size = mSize
		music.Source = "qq"
		musicList = append(musicList, music)
	}
	return musicList
}

func Download(music models.Music) {
	// 根据songmid等信息获得下载链接
	guid := common.Random(100000000, 10000000000)
	req, err := http.NewRequest("GET", "http://base.music.qq.com/fcgi-bin/fcg_musicexpress.fcg", nil)
	common.AddHeader(req)
	req.Header.Set("referer", "http://m.y.qq.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")

	q := req.URL.Query()
	q.Add("guid", strconv.Itoa(guid))
	q.Add("format", "json")
	q.Add("json", "3")

	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}
	//fmt.Printf("%s\n", content)

	respJson := make(map[string]interface{})
	err = json.Unmarshal(content, &respJson)
	if err != nil {
		log.Panic(err)
	}

	vkey := respJson["key"]
	prefixs := []string{"M800", "M500", "C400"}
	for _, prefix := range prefixs {
		url := fmt.Sprintf("http://dl.stream.qqmusic.qq.com/%v%v.mp3?vkey=%v&guid=%v&fromtag=1", prefix, music.MID, vkey, guid)
		fmt.Println(url)
		size := common.GetContentLen(url)
		mSize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1048576), 64)
		if size > 0 {
			music.Url = url
			if prefix == "M800" {
				music.Rate = "320"
			} else {
				music.Rate = "128"
			}
			music.Size = mSize
			break
		}
	}
	music.Name = fmt.Sprintf("%v - %v.mp3", music.Singer, music.Title)
	common.MusicDownload(music)
}
