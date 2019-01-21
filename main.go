package main

import (
	"fmt"
	"github.com/zou2699/musicSD/pkg/netease"
	"github.com/zou2699/musicSD/pkg/qq"
)

func main() {
	for {
		var name string
		fmt.Println("请输入要搜索的歌曲，名称和歌手一起输入可以提高匹配(如 Beyond 海阔天空): ")
		fmt.Scanln(&name)
		if name == "" {
			continue
		}
		fmt.Println("开始搜索...")
		musicList := netease.Search(name)
		musicList = append(musicList, qq.Search(name)...)
		fmt.Println("id -- title --- singer --- size")
		for id, music := range musicList {
			fmt.Printf("%v -- ", id)
			fmt.Printf("%v --- ", music.Source)
			fmt.Printf("%s --- ", music.Title)
			fmt.Printf("%s --- ", music.Singer)
			fmt.Printf("%vMB\n", music.Size)
		}

		var id int
		fmt.Println("请输入要下载的歌曲序号, 多个序号用空格隔开: ")
		_, err := fmt.Scanln(&id)
		if err != nil {
			fmt.Printf("输入序号错误\n\n")
			continue
		}

		music := musicList[id]
		if music.Source == "qq" {
			qq.Download(music)
		} else if music.Source == "netease" {
			netease.Download(music)
		}

	}
	//a := fmt.Sprintf("%#x",rune(1))
	//fmt.Println(strings.Repeat(a,2))
}
