package main

import (
	"fmt"
	"github.com/zou2699/musicSD/pkg/api"
)

func main() {
	for {
		var name string
		fmt.Println("请输入要搜索的歌曲，名称和歌手一起输入可以提高匹配(如 海阔天空): ")
		fmt.Scanln(&name)
		if name == "" {
			continue
		}
		musicList := api.NeteaseSearch(name)
		fmt.Println("id -- title --- singer --- size")
		for id, music := range musicList {
			fmt.Printf("%v -- ", id)
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
		api.NeteaseDownload(musicList[id])
	}
	//a := fmt.Sprintf("%#x",rune(1))
	//fmt.Println(strings.Repeat(a,2))
}
