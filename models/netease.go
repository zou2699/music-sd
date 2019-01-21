package models

type RespNetease struct {
	Result struct {
		Songs []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			Pst  int    `json:"pst"`
			T    int    `json:"t"`
			Ar   []struct {
				ID    int           `json:"id"`
				Name  string        `json:"name"`
				Tns   []interface{} `json:"tns"`
				Alias []interface{} `json:"alias"`
			} `json:"ar"`
			Alia []interface{} `json:"alia"`
			Pop  float64       `json:"pop"`
			St   int           `json:"st"`
			Rt   interface{}   `json:"rt"`
			Fee  int           `json:"fee"`
			V    int           `json:"v"`
			Crbt interface{}   `json:"crbt"`
			Cf   string        `json:"cf"`
			Al   struct {
				ID     int           `json:"id"`
				Name   string        `json:"name"`
				PicURL string        `json:"picUrl"`
				Tns    []interface{} `json:"tns"`
				Pic    int64         `json:"pic"`
			} `json:"al"`
			Dt int `json:"dt"`
			H  struct {
				Br   int     `json:"br"`
				Fid  int     `json:"fid"`
				Size int     `json:"size"`
				Vd   float64 `json:"vd"`
			} `json:"h"`
			M struct {
				Br   int     `json:"br"`
				Fid  int     `json:"fid"`
				Size int     `json:"size"`
				Vd   float64 `json:"vd"`
			} `json:"m"`
			L struct {
				Br   int     `json:"br"`
				Fid  int     `json:"fid"`
				Size int     `json:"size"`
				Vd   float64 `json:"vd"`
			} `json:"l"`
			A           interface{}   `json:"a"`
			Cd          string        `json:"cd"`
			No          int           `json:"no"`
			RtURL       interface{}   `json:"rtUrl"`
			Ftype       int           `json:"ftype"`
			RtUrls      []interface{} `json:"rtUrls"`
			DjID        int           `json:"djId"`
			Copyright   int           `json:"copyright"`
			SID         int           `json:"s_id"`
			Rtype       int           `json:"rtype"`
			Rurl        interface{}   `json:"rurl"`
			Mst         int           `json:"mst"`
			Cp          int           `json:"cp"`
			Mv          int           `json:"mv"`
			PublishTime int64         `json:"publishTime"`
			Privilege   struct {
				ID    int  `json:"id"`
				Fee   int  `json:"fee"`
				Payed int  `json:"payed"`
				St    int  `json:"st"`
				Pl    int  `json:"pl"`
				Dl    int  `json:"dl"`
				Sp    int  `json:"sp"`
				Cp    int  `json:"cp"`
				Subp  int  `json:"subp"`
				Cs    bool `json:"cs"`
				Maxbr int  `json:"maxbr"`
				Fl    int  `json:"fl"`
				Toast bool `json:"toast"`
				Flag  int  `json:"flag"`
			} `json:"privilege"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

type Music struct {
	Album    string  `json:"album"`
	Singer   string  `json:"singer"`
	Source   string  `json:"source"`
	Duration string  `json:"duration"`
	Title    string  `json:"title"`
	ID       int     `json:"id"`
	Size     float64 `json:"size"`
	MID      string  `json:"mid"`
	Url      string  `json:"url"`
	Rate     string  `json:"rate"`
	Name     string  `json:"name"`
}

type MusicDownloadNetease struct {
	Data []struct {
		ID            int         `json:"id"`
		URL           string      `json:"url"`
		Br            int         `json:"br"`
		Size          int         `json:"size"`
		Md5           string      `json:"md5"`
		Code          int         `json:"code"`
		Expi          int         `json:"expi"`
		Type          string      `json:"type"`
		Gain          float64     `json:"gain"`
		Fee           int         `json:"fee"`
		Uf            interface{} `json:"uf"`
		Payed         int         `json:"payed"`
		Flag          int         `json:"flag"`
		CanExtend     bool        `json:"canExtend"`
		FreeTrialInfo interface{} `json:"freeTrialInfo"`
	} `json:"data"`
	Code int `json:"code"`
}
