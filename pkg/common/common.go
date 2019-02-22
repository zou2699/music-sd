package common

import (
	"encoding/hex"
	"fmt"
	"github.com/zou2699/music-sd/models"
	"github.com/zou2699/music-sd/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

func AddHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Charset", "UTF-8,*;q=0.5")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0")
	req.Header.Add("referer", "http://google.com/")
}

// 加密postForm
func EncryptForm(requestBytes []byte) (encryptedString string) {
	//获取key的二进制字符串作为key
	key, _ := hex.DecodeString("7246674226682325323F5E6544673A51")
	encryptedBytes := utils.AESEncrypt(requestBytes, key)
	//获取加密后的字符串
	encryptedString = hex.EncodeToString(encryptedBytes)
	//解密
	//https://tools.lami.la/jiami/aes
	//decodeBytes, _ := hex.DecodeString(encryptedString)
	//decrypted := AESDecrypt(decodeBytes, key)
	//fmt.Printf("%s", decrypted)
	//fmt.Println(encryptedString)
	return encryptedString

}

// 下载
func MusicDownload(music models.Music) {

	filename := music.Name
	//fmt.Println(music.Name)
	if runtime.GOOS == "windows" {
		compile, err := regexp.Compile("[\\/:*?\"<>|]")
		if err != nil {
			log.Panic(err)
		}
		filename = compile.ReplaceAllString(filename, ",")
	}

	//fmt.Println("开始下载", filename)
	t1 := time.Now()

	response, err := http.Get(music.Url)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Panic(err)
	}

	elapsed := time.Since(t1)
	fmt.Printf("%s 下载完成, 音质%skbps, 耗时: %v\n", filename, music.Rate, elapsed)
}

// Head the url
func GetContentLen(url string) int {
	response, err := http.Head(url)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(response.ContentLength, response.StatusCode)
	if response.StatusCode == http.StatusOK {
		len := response.Header.Get("Content-Length")
		if len != "" {
			i, err := strconv.Atoi(len)
			if err != nil {
				log.Panic(err)
			}
			return i
		}
		return 0
	}
	return 0
}

// 返回随机数
func Random(min, max int) int {
	return rand.Intn(max-min) + min
}
