package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	TransResult struct {
		Data []struct {
			Dst string `json:"dst"`
			Src string `json:"src"`
		} `json:"data"`
		From     string `json:"from"`
		To       string `json:"to"`
		Status   int    `json:"status"`
		Type     int    `json:"type"`
		Phonetic []struct {
			SrcStr string `json:"src_str"`
			TrgStr string `json:"trg_str"`
		} `json:"phonetic"`
	} `json:"trans_result"`
	DictResult struct {
		Edict struct {
			Item []struct {
				TrGroup []struct {
					Tr          []string `json:"tr"`
					Example     []string `json:"example"`
					SimilarWord []string `json:"similar_word"`
				} `json:"tr_group"`
				Pos string `json:"pos"`
			} `json:"item"`
			Word string `json:"word"`
		} `json:"edict"`
		From        string `json:"from"`
		SimpleMeans struct {
			WordName  string   `json:"word_name"`
			From      string   `json:"from"`
			WordMeans []string `json:"word_means"`
			Tags      struct {
				Core  []string `json:"core"`
				Other []string `json:"other"`
			} `json:"tags"`
			Exchange struct {
				WordPl []string `json:"word_pl"`
			} `json:"exchange"`
			Symbols []struct {
				PhEn  string `json:"ph_en"`
				PhAm  string `json:"ph_am"`
				Parts []struct {
					Part  string   `json:"part"`
					Means []string `json:"means"`
				} `json:"parts"`
				PhOther string `json:"ph_other"`
			} `json:"symbols"`
		} `json:"simple_means"`
		Common struct {
			Text string `json:"text"`
		} `json:"common"`
		Collins struct {
			Entry []struct {
				Type    string `json:"type"`
				EntryID string `json:"entry_id"`
				Value   []struct {
					MeanType []struct {
						InfoType string `json:"info_type"`
						InfoID   string `json:"info_id"`
						Example  []struct {
							ExampleID string `json:"example_id"`
							TtsSize   string `json:"tts_size"`
							Tran      string `json:"tran"`
							Ex        string `json:"ex"`
							TtsMp3    string `json:"tts_mp3"`
						} `json:"example,omitempty"`
						Posc []struct {
							Tran    string `json:"tran"`
							PoscID  string `json:"posc_id"`
							Example []struct {
								ExampleID string `json:"example_id"`
								Tran      string `json:"tran"`
								Ex        string `json:"ex"`
								TtsMp3    string `json:"tts_mp3"`
							} `json:"example"`
							Def string `json:"def"`
						} `json:"posc,omitempty"`
					} `json:"mean_type"`
					Gramarinfo []struct {
						Tran  string `json:"tran"`
						Type  string `json:"type"`
						Label string `json:"label"`
					} `json:"gramarinfo"`
					Tran   string `json:"tran"`
					Def    string `json:"def"`
					MeanID string `json:"mean_id"`
					Posp   []struct {
						Label string `json:"label"`
					} `json:"posp"`
				} `json:"value"`
			} `json:"entry"`
			WordName      string `json:"word_name"`
			WordID        string `json:"word_id"`
			WordEmphasize string `json:"word_emphasize"`
			Frequence     string `json:"frequence"`
		} `json:"collins"`
		Lang   string `json:"lang"`
		Oxford struct {
			Entry []struct {
				Tag  string `json:"tag"`
				Name string `json:"name"`
				Data []struct {
					Tag  string `json:"tag"`
					Data []struct {
						Tag  string `json:"tag"`
						Data []struct {
							Tag  string `json:"tag"`
							Data []struct {
								Tag  string `json:"tag"`
								Data []struct {
									Tag    string `json:"tag"`
									EnText string `json:"enText,omitempty"`
									ChText string `json:"chText,omitempty"`
									G      string `json:"g,omitempty"`
									Data   []struct {
										Text      string `json:"text"`
										HoverText string `json:"hoverText"`
									} `json:"data,omitempty"`
								} `json:"data"`
							} `json:"data"`
						} `json:"data,omitempty"`
						P     string `json:"p,omitempty"`
						PText string `json:"p_text,omitempty"`
						N     string `json:"n,omitempty"`
						Xt    string `json:"xt,omitempty"`
					} `json:"data"`
				} `json:"data"`
			} `json:"entry"`
			Unbox []struct {
				Tag  string `json:"tag"`
				Type string `json:"type"`
				Name string `json:"name"`
				Data []struct {
					Tag     string `json:"tag"`
					Text    string `json:"text,omitempty"`
					Words   string `json:"words,omitempty"`
					Outdent string `json:"outdent,omitempty"`
					Data    []struct {
						Tag    string `json:"tag"`
						EnText string `json:"enText"`
						ChText string `json:"chText"`
					} `json:"data,omitempty"`
				} `json:"data"`
			} `json:"unbox"`
		} `json:"oxford"`
		BaiduPhrase []struct {
			Tit   []string `json:"tit"`
			Trans []string `json:"trans"`
		} `json:"baidu_phrase"`
		Sanyms []struct {
			Tit  string `json:"tit"`
			Type string `json:"type"`
			Data []struct {
				P string   `json:"p"`
				D []string `json:"d"`
			} `json:"data"`
		} `json:"sanyms"`
		QueryExplainVideo struct {
			ID           int    `json:"id"`
			UserID       string `json:"user_id"`
			UserName     string `json:"user_name"`
			UserPic      string `json:"user_pic"`
			Query        string `json:"query"`
			Direction    string `json:"direction"`
			Type         string `json:"type"`
			Tag          string `json:"tag"`
			Detail       string `json:"detail"`
			Status       string `json:"status"`
			SearchType   string `json:"search_type"`
			FeedURL      string `json:"feed_url"`
			Likes        string `json:"likes"`
			Plays        string `json:"plays"`
			CreatedAt    string `json:"created_at"`
			UpdatedAt    string `json:"updated_at"`
			DuplicateID  string `json:"duplicate_id"`
			RejectReason string `json:"reject_reason"`
			CoverURL     string `json:"coverUrl"`
			VideoURL     string `json:"videoUrl"`
			ThumbURL     string `json:"thumbUrl"`
			VideoTime    string `json:"videoTime"`
			VideoType    string `json:"videoType"`
		} `json:"queryExplainVideo"`
	} `json:"dict_result"`
	LijuResult struct {
		Double string   `json:"double"`
		Tag    []string `json:"tag"`
		Single string   `json:"single"`
	} `json:"liju_result"`
	Logid int `json:"logid"`
}

type DictRequest2 struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse2 struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func query(wg *sync.WaitGroup, word string) {

	defer wg.Done()
	client := &http.Client{}
	var data = strings.NewReader(`from=en&to=zh&query=hello&transtype=translang&simple_means_flag=3&sign=54706.276099&token=901cdbb2ddcd3168341657ae2a0ac7c5&domain=common&ts=1690521569118`)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/v2transapi?from=en&to=zh", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Acs-Token", "1690521554994_1690521569201_1cx2YWCDMMM9QeQ9nOHLd6RwvvbuSShY/VjejoIQkVtaWhjEyaMzEfQswuGQUOTXBqw1rtVZpj2FYgBXLDdHLgGjMXJ9+PuYeF5jisG7dlYuB9pfG69w/UcHg336c19xpgy6M9EbkIkdv4Ko+tAAQRkWpGtWbmDCdiImBhkwqITgUcaCvvveAlj6uPueaz09Gfa+yDgUFm8h+UOAv6XSvyV8YkHKy+LaAvTuhWC/AZbaTPDqguZOg04q8iGaqs4Kfv68TI69Pwg4c+mxMwVtupiYhyxnOka3nECxLr9Mnl4EokeJIzY0jDopqtoRI0Yx42HW3QMOoCKr8ow5JFPC6D2IEwsZnsqhu6aQaPOly7bvQIinWXvaQsEvlJObbBUnzk6sKWHFPlqPr+gwM9Dx9X5jd/xst6/rDp+6aajcXZ/kpI8cSmPwnE73Nlbm0EOs/fqwHg5Fg/Nf+uQHbqxIQtfwOQf/9DETrRrt0MutY7QHRk1pwF/Pn8vQGvcadCQZ")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "BAIDUID=FEAC8DCFE5A285ED37B2AEE52E44C5D9:FG=1; BAIDUID_BFESS=FEAC8DCFE5A285ED37B2AEE52E44C5D9:FG=1; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1690470395; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1690521556; ab_sr=1.0.1_MWZlMzRlOTA3OWI1NjIzZDU4NzA2YzY5YTI4NGZhODBjYWM1Mjk0ZWQxNjE5M2NjNDA5NTY4NDdjMjAwMzRhMGI2N2QzYWRjYjhmNzM3ZDkwODY5MmUyMmI1NWY3ZTM5ZjVmZWJiN2Y5ODg0M2RjMjhmYTIyZjhiYzVhN2I5ZmEwYTJjOGFjMTIxMWIxNDk1N2Y1ODM1ZWJmMTIzMWJmNA==")
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.DictResult.SimpleMeans.Symbols[0].PhEn, "US:", dictResponse.DictResult.SimpleMeans.Symbols[0].PhAm)
	for _, item := range dictResponse.DictResult.SimpleMeans.Symbols[0].Parts[0].Means {
		fmt.Println(item)
	}
}

func query2(wg *sync.WaitGroup, word string) {
	defer wg.Done()
	client := &http.Client{} //创建一个http client，这个函数还可以指定最大的返回时间
	request := DictRequest2{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request) //设置一个request结构体，将json序列化
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	//生成一个http请求
	if err != nil {
		log.Fatal(err)
	}

	//设置请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")

	resp, err := client.Do(req) //真正发起请求
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close() //按照go的习惯，需要手动关闭resp流，defer的意思是在函数结束后关闭，释放内存
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 { //得到的response不一定正确，需要检测一下返回的状态码
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse2
	err = json.Unmarshal(bodyText, &dictResponse) //构造一个与json相同的结构体，并将json序列化到这个结构体变量里面
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs) //得到音标
	for _, item := range dictResponse.Dictionary.Explanations {                                           //得到解释
		fmt.Println(item)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	var wg sync.WaitGroup
	wg.Add(2)
	// 并发执行 function1 和 function2
	go query(&wg, word)
	go query2(&wg, word)
	// 等待所有函数执行完成
	wg.Wait()
	fmt.Println("All functions completed")
}
