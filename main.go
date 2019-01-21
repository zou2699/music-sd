package main

import (
	"bufio"
	"fmt"
	"github.com/zou2699/musicSD/pkg/netease"
	"github.com/zou2699/musicSD/pkg/qq"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	for {
		var name string
		fmt.Println("请输入要搜索的歌曲，名称和歌手一起输入可以提高匹配(如 海阔天空 Beyond): ")
		fmt.Scanln(&name)
		//name = "海阔天空"
		if name == "" {
			continue
		}
		fmt.Println("开始搜索...")
		musicList := netease.Search(name)
		musicList = append(musicList, qq.Search(name)...)
		for id, music := range musicList {
			fmt.Printf("[%2d] %7s | %s %5sMB - %s - %s - %s\n", id, music.Source, music.Duration, music.Size, music.Title, music.Singer, music.Album)
		}

		fmt.Println("请输入要下载的歌曲序号, 多个序号用空格隔开: ")
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			fmt.Printf("输入序号错误\n\n")
			continue
		}
		fmt.Println("src ", input)
		ids := strings.Fields(input)

		//_, err = fmt.Scanln(&ids)
		//if err != nil {
		//	fmt.Println(err)
		//	fmt.Printf("输入序号错误\n\n")
		//	continue
		//}
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			i, err := strconv.Atoi(id)
			if err != nil {
				log.Panic(err)
			}

			music := musicList[i]
			if music.Source == "QQ" {
				go func() {
					qq.Download(music)
					wg.Done()
				}()
			} else if music.Source == "NETEASE" {
				go func() {
					netease.Download(music)
					wg.Done()
				}()
			}
		}
		wg.Wait()
		fmt.Println("###################")
		fmt.Println()
	}
}
