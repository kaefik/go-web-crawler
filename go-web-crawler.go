// go-web-crawler
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go-web-crawler/pick"

	"golang.org/x/net/html/charset"

	//	"net/smtp"
	//	"strconv"
	"strings"
)

//---- инициализация глобальных типов и переменных
type ListUrl struct {
	url       string // урл
	fdownload bool   // флаг того что данный урл был загружен для анализа
}

//---- END инициализация глобальных типов и переменных

// инициализация файла логов
func InitLogFile(namef string) *log.Logger {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}
	multi := io.MultiWriter(file, os.Stdout)
	LFile := log.New(multi, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LFile
}

//получение страницы из урла url
func gethtmlpage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
		panic("HTTP error")
	}
	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("Encoding error:", err)
		panic("Encoding error")
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		fmt.Println("IO error:", err)
		panic("IO error")
	}
	return body
}

//сохранить данные в файл
func Savetofile(namef string, str string) error {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// handle the error here
		return err
	}
	file.WriteString(str)
	return err
}

// получение всех ссылок на страницы
func getLnksfromPage(body []byte) []string {
	res := make([]string, 0)
	shtml := string(body)
	//	fmt.Println(shtml)

	res, _ = pick.PickAttr(&pick.Option{&shtml, "a", nil}, "href")
	return res

}

//выборка из списка урл ll только те которые являются внутренними страницами указанного домена dom
func internalLinksfromSite(ll []string, dom string) []string {
	res := make([]string, 0)
	for _, v := range ll {
		if len(v) > 0 {
			if v[0] == '/' {
				res = append(res, dom+v)
			} else {
				if strings.HasPrefix(v, dom) {
					res = append(res, v)
				}
			}
		}
	}
	return res
}

//возвращает массив строк которые получается при сравнении массивов list1 и list2
//и если нет строки из list2 в массиве list1
func UniqLinks(list1, list2 []string) []string {
	res := make([]string, 0)
	f := false
	for _, v2 := range list2 {
		for _, v1 := range list1 {
			if v1 == v2 {
				f = true
				//				break
			}
		}
		if f != true {
			res = append(res, v2)
			f = false
		}
	}
	return res
}

func main() {
	fmt.Println("Hello World!")
	myurl := "http://echo.msk.ru"
	//	flagEnd := false // флаг окончания выгрузки
	//	lurl := make([]ListUrl, 0)
	c := 0
	for {
		if c == 1 {
			break
		} else {
			body := gethtmlpage(myurl)
			listlinks := getLnksfromPage(body)
			fmt.Println(listlinks)
			fmt.Println(internalLinksfromSite(listlinks, myurl))
		}
		c++
	}

}
