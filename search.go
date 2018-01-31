package qanswer

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/ngaut/log"
	"github.com/silenceper/qanswer/config"
	"github.com/silenceper/qanswer/proto"
	"github.com/silenceper/qanswer/util"
)

//SearchResult 搜索总数跟频率
type SearchResult struct {
	Sum  int32
	Freq int32
}

//GetSearchResult 返回各个搜索引擎返回的结果
func GetSearchResult(question string, answers []string) map[string][]*SearchResult {
	if question == "" {
		return nil
	}
	res := make(map[string][]*SearchResult)
	res["百度"] = baiduSearch(question, answers)
	return res
}

func baiduSearch(question string, answers []string) (result []*SearchResult) {
	resultMap := make(map[string]*SearchResult, len(answers))
	for k, answer := range answers {
		answer = plainAnswer(answer)
		answers[k] = answer

		searchResult := new(SearchResult)
		resultMap[answer] = searchResult
	}

	var wg sync.WaitGroup
	//搜索题目
	wg.Add(1)
	go func() {
		defer wg.Done()
		searchURL := fmt.Sprintf("http://www.baidu.com/s?wd=%s", url.QueryEscape(question))
		//		questionBody, err := util.HTTPGet(searchURL, 5)
		//		if err != nil {
		//			log.Errorf("search question:%s error", question)
		//			return
		//		}

		doc, err := goquery.NewDocument(searchURL)
		questionBody := doc.Text()
		if err == nil {
			span := doc.Find("span.c-gap-right-small")
			color.Green("百度首页推荐信息:\n%s", strings.TrimSpace(span.Text()))
			// 空格表示有多个类
			doc.Find("div.result.c-container").Each(func(i int, s *goquery.Selection) {
				i = i + 1
				if i > 5 {
					return
				}

				resultId := s.Find("div.c-abstract")
				color.Green("搜索栏%d:\n%s\n", i, strings.TrimSpace(resultId.Text()))
			})
		} else {
			log.Errorf("goquery error, %s", err)
		}

		// 保存到本地
		if config.GetConfig().Debug {
			os.MkdirAll(proto.SearchResultPath, os.ModeDir)
			file, createErr := os.Create(proto.SearchResultPath + "QuestionRes.txt")
			if createErr != nil {
				log.Error(createErr)
			} else {
				_, writeErr := file.Write([]byte(questionBody))
				if writeErr != nil {
					log.Error(writeErr)
				}

				file.Close()
			}
		}

		for _, answer := range answers {
			//题目搜索结果中包含的答案的数量
			resultMap[answer].Freq = int32(strings.Count(string(questionBody), answer))
		}
	}()

	for _, answer := range answers {
		wg.Add(1)
		go func(answer string) {
			defer wg.Done()
			//题目+结果搜索的总数
			keyword := fmt.Sprintf("%s %s", question, answer)
			searchURL := fmt.Sprintf("http://www.baidu.com/s?wd=%s", url.QueryEscape(keyword))
			body, err := util.HTTPGet(searchURL, 5)
			if err != nil {
				log.Errorf("search %s error", answer)
			} else {
				countRe, _ := regexp.Compile(`百度为您找到相关结果约([\d\,]+)`)
				result := countRe.FindAllStringSubmatch(string(body), -1)
				if len(result) > 0 {
					sum := result[0][1]
					sum = strings.Replace(sum, ",", "", -1)
					resultMap[answer].Sum = util.MustInt32(sum)
				}
			}
		}(answer)
	}
	wg.Wait()

	//将map转为slice 方便顺序输出
	for _, answer := range answers {
		result = append(result, resultMap[answer])
	}
	return result
}

//plainAnswer 去除答案中的 《》等字符
func plainAnswer(answer string) string {
	answer = strings.TrimPrefix(answer, "《")
	answer = strings.TrimSuffix(answer, "》")
	return answer
}
