package netease

import (
	"encoding/json"
	"fmt"
	"github.com/zou2699/music-sd/models"
	"github.com/zou2699/music-sd/pkg/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 搜索netease音乐
func Search(keyword string) (musicList []models.Music) {
	// 初始化requestJson
	requestJSON := map[string]interface{}{
		"method": "POST",
		"url":    "http://music.163.com/api/cloudsearch/pc",
		"params": map[string]interface{}{
			"s":      keyword,
			"type":   1,
			"offset": 0,
			"limit":  8,
		},
	}

	//fmt.Println("requestJSON:",requestJSON)
	// json化数据
	requestBytes, err := json.Marshal(requestJSON)
	if err != nil {
		log.Fatal(err)
	}
	// 加密postForm
	encryptedString := common.EncryptForm(requestBytes)

	// post form
	form := url.Values{}
	form.Add("eparams", encryptedString)

	req, err := http.NewRequest("POST", "http://music.163.com/api/linux/forward", strings.NewReader(form.Encode()))

	// FAKE_HEADERS
	common.AddHeader(req)
	req.Header.Set("referer", "http://music.163.com/")

	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}

	var respJson models.RespNetease
	err = json.Unmarshal(content, &respJson)
	if err != nil {
		log.Panic(err)
	}

	if respJson.Code != 200 {
		log.Panic("code not 200,it's", respJson.Code)
	}

	// 结构化music
	var music models.Music
	for _, song := range respJson.Result.Songs {

		if song.Privilege.Fl == 0 {
			// 没有版权
			continue
		}
		// singer
		var singers []string
		for _, singer := range song.Ar {
			singers = append(singers, singer.Name)
		}

		// 大小
		var size int
		if song.Privilege.Fl >= 320000 {
			size = song.H.Size
		} else if song.Privilege.Fl >= 192000 {
			size = song.M.Size
		} else {
			size = song.L.Size
		}
		// 转化位MB
		//mSize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1048576), 64)
		mSize := fmt.Sprintf("%.2f", float64(size)/1048576)
		//fmt.Printf("%.2f", float64(size)/1048576)
		music.ID = song.ID
		music.Title = song.Name
		music.Album = song.Al.Name
		music.Size = mSize
		music.Duration = time.Unix(int64(song.Dt)/1000, 0).Format("04:05")
		music.Source = "NETEASE"
		music.Singer = strings.Join(singers, ",")

		musicList = append(musicList, music)
	}
	return musicList
}

// 下载netease音乐
func Download(music models.Music) {
	musicId := "[" + strconv.Itoa(music.ID) + "]"
	// 初始化requestJson
	requestJSON := map[string]interface{}{
		"method": "POST",
		"url":    "http://music.163.com/api/song/enhance/player/url",
		"params": map[string]interface{}{
			"ids": musicId,
			"br":  320000,
		},
	}

	// json化数据
	requestBytes, err := json.Marshal(requestJSON)
	if err != nil {
		log.Panic(err)
	}
	encryptedString := common.EncryptForm(requestBytes)

	// post form
	form := url.Values{}
	form.Add("eparams", encryptedString)

	req, err := http.NewRequest("POST", "http://music.163.com/api/linux/forward", strings.NewReader(form.Encode()))

	// FAKE_HEADERS
	common.AddHeader(req)
	req.Header.Set("referer", "http://music.163.com/")

	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	//fmt.Println(resp.StatusCode)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}
	//fmt.Printf("content: %s\n", content)

	var respJson models.MusicDownloadNetease

	err = json.Unmarshal(content, &respJson)
	if err != nil {
		log.Panic(err)
	}

	//fmt.Println(respJson.Code)
	if respJson.Code != 200 {
		log.Panic("code not 200 ", respJson.Code)
	}

	music.Url = respJson.Data[0].URL
	music.Name = fmt.Sprintf("%v - %v.%v", music.Singer, music.Title, respJson.Data[0].Type)
	music.Rate = strconv.Itoa(respJson.Data[0].Br / 1000)

	common.MusicDownload(music)
}
