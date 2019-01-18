package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zou2699/musicSD/models"
	"github.com/zou2699/musicSD/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func NeteaseSearch(keyword string) (musicList []models.MusicNetease) {
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
	//fmt.Printf("%s\n",requestBytes)

	// 获取hex的二进制字符串作为key
	// 获取hex的二进制字符串作为key
	//key := []byte("rFgB&h#%2?^eDg:Q")
	//key, _ := hex.DecodeString("7246674226682325323F5E6544673A51")
	//
	////fmt.Printf("key:%s\n", key)
	//
	//encryptedBytes := utils.AESEncrypt(requestBytes, key)
	//encryptedString := hex.EncodeToString(encryptedBytes)

	encryptedString := EncryptForm(requestBytes)
	//fmt.Printf("encrypted:%v\nhexs:%v\n", encryptedBytes, encryptedString)

	//解密
	//https://tools.lami.la/jiami/aes
	//decodeBytes, _ := hex.DecodeString(encryptedString)
	//decrypted := AESDecrypt(decodeBytes, key)
	//fmt.Printf("%s", decrypted)

	client := http.Client{}
	if err != nil {
		fmt.Println(err)
	}

	form := url.Values{}
	form.Add("eparams", encryptedString)
	//ssk:="18e16cd68fb6c5d75fb0bd7c92e10fc028c2609c253bf23133504e777bd2c9002e9ff9b4b91baca6e61c6d9bee34d9141f3cf8bd67567d4554f8799e56f536ba470862bedc825921b90b3396936ba22ce08a478a1c93ad389a54fba7b8210d03fe7c3ffcfb4fbce6b9428db3c99c6fcef02b9126c60d802e3a1e8406378adc14f1b6efc49f1b3163848ba9cda124b00633d91d3044d8f8b7b2f87842fe1397a8"
	//form.Add("eparams","18e16cd68fb6c5d75fb0bd7c92e10fc028c2609c253bf23133504e777bd2c9002e9ff9b4b91baca6e61c6d9bee34d9141f3cf8bd67567d4554f8799e56f536ba470862bedc825921b90b3396936ba22ce08a478a1c93ad389a54fba7b8210d03fe7c3ffcfb4fbce6b9428db3c99c6fcef02b9126c60d802e3a1e8406378adc14f1b6efc49f1b3163848ba9cda124b00633d91d3044d8f8b7b2f87842fe1397a8")
	//form.Add("name","admin")
	//form.Add("password","admin")
	//var jsonStr = []byte(`{"eparams":"18e16cd68fb6c5d75fb0bd7c92e10fc028c2609c253bf23133504e777bd2c9002e9ff9b4b91baca6e61c6d9bee34d9141f3cf8bd67567d4554f8799e56f536ba470862bedc825921b90b3396936ba22ce08a478a1c93ad389a54fba7b8210d03fe7c3ffcfb4fbce6b9428db3c99c6fcef02b9126c60d802e3a1e8406378adc14f1b6efc49f1b3163848ba9cda124b00633d91d3044d8f8b7b2f87842fe1397a8"}`)

	req, err := http.NewRequest("POST", "http://music.163.com/api/linux/forward", strings.NewReader(form.Encode()))

	//req, err := http.NewRequest("POST", "http://music.163.com/api/linux/forward", strings.NewReader(form.Encode()))

	//fmt.Println(form.Encode())
	// FAKE_HEADERS
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Charset", "UTF-8,*;q=0.5")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0")
	req.Header.Add("referer", "http://music.163.com/")

	//fmt.Println(req.Header)

	fmt.Println("开始搜索")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("%s\n", content)

	//respJson := make(map[string]interface{})

	var respJson models.RespNetease

	err = json.Unmarshal(content, &respJson)
	if err != nil {
		log.Panic(err)
	}

	if respJson.Code != 200 {
		log.Panic("code not 200,it's", respJson.Code)
	}
	//fmt.Printf("%+v\n", respJson)

	var music models.MusicNetease
	//var musicList []models.MusicNetease
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
		msize, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1048576), 64)

		//fmt.Printf("%.2f", float64(size)/1048576)
		music.ID = song.ID
		music.Title = song.Name
		music.Album = song.Al.Name
		music.Size = msize
		music.Duration = ""
		music.Source = "netease"
		music.Singer = strings.Join(singers, "/")

		musicList = append(musicList, music)
		//	fmt.Printf("%+v\n", music)
	}
	return musicList
}

func NeteaseDownload(music models.MusicNetease) {
	// 初始化requestJson
	musicId := "[" + strconv.Itoa(music.ID) + "]"
	//musicId = "[347230]"
	//fmt.Println(music)
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
		log.Fatal(err)
	}
	encryptedString := EncryptForm(requestBytes)

	//解密
	//https://tools.lami.la/jiami/aes
	//key, _ := hex.DecodeString("7246674226682325323F5E6544673A51")
	//decodeBytes, _ := hex.DecodeString(encryptedString)
	//decrypted := utils.AESDecrypt(decodeBytes, key)
	//fmt.Printf("%s\n", decrypted)

	client := http.Client{}
	if err != nil {
		fmt.Println(err)
	}

	form := url.Values{}
	form.Add("eparams", encryptedString)
	//fmt.Println(encryptedString)
	req, err := http.NewRequest("POST", "http://music.163.com/api/linux/forward", strings.NewReader(form.Encode()))

	// FAKE_HEADERS
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Charset", "UTF-8,*;q=0.5")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0")
	req.Header.Add("referer", "http://music.163.com/")

	//req.Host = "http://music.163.com"
	//fmt.Println(req.Header)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(resp.StatusCode)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
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

	mUrl := respJson.Data[0].URL

	fmt.Println("开始下载", music.Singer+" - "+music.Title)
	t1 := time.Now()
	response, err := http.Get(mUrl)
	if err != nil {
		log.Panic(err)
	}
	file, err := os.Create("./" + music.Singer + " - " + music.Title + "." + respJson.Data[0].Type)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Panic(err)
	}

	elapsed := time.Since(t1)
	fmt.Println("下载完成, 共用时", elapsed)

}

func EncryptForm(requestBytes []byte) (encryptedString string) {
	key, _ := hex.DecodeString("7246674226682325323F5E6544673A51")
	encryptedBytes := utils.AESEncrypt(requestBytes, key)
	encryptedString = hex.EncodeToString(encryptedBytes)
	//fmt.Println(encryptedString)
	return encryptedString
}
