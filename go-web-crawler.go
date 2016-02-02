// go-web-crawler
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go-web-crawler/pick"

	"golang.org/x/net/html/charset"

	//	"net/smtp"
	//	"strconv"
	"strings"
)

//---- инициализация глобальных типов и переменных
type ListUrl struct {
	url       string // урл
	fdownload int    // флаг того что данный урл был загружен для анализа: if 1 then загружен ; if 0 не загружен; if -1 ошибка при загрузке
}

var (
	site    string
	koliter string
)

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
func gethtmlpage(url string) ([]byte, bool) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
		//		panic("HTTP error")
		return make([]byte, 0), false
	}
	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("Encoding error:", err)
		return make([]byte, 0), false
		//		panic("Encoding error")
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		fmt.Println("IO error:", err)
		return make([]byte, 0), false
		//		panic("IO error")
	}
	return body, true
}

//сохранить данные в файл
func Savetofile(namef string, str string) error {
	file, err := os.Create(namef) //OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// handle the error here
		return err
	}
	defer file.Close()
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
func internalLinksfromSiteListUrl(ll []ListUrl, dom string) []ListUrl {
	res := make([]ListUrl, 0)
	for _, v := range ll {
		if len(v.url) > 0 {
			if v.url[0] == '/' {
				res = append(res, ListUrl{url: dom + v.url, fdownload: 0})
			} else {
				if strings.HasPrefix(v.url, dom) {
					res = append(res, ListUrl{url: v.url, fdownload: 0})
				}
			}
		}
	}
	return res
}

//возвращает массив строк которые получается при сравнении массивов list1 и list2
//и если нет строки из list2 в массиве list1
func UniqLinks2(list1 []ListUrl, list2 []ListUrl) []ListUrl {
	res := make([]ListUrl, 0)
	for _, v2 := range list2 {
		f := false
		for _, v1 := range list1 {
			if v1.url == v2.url { // && (v1.fdownload == v2.fdownload) { //strings.Compare(v1, v2) == 0 {
				f = true
				break
			}
		}
		if f != true {
			res = append(res, v2)
			f = false
		}
	}
	return res
}

func AddtoEndList2(l1 []ListUrl, l2 []ListUrl) []ListUrl {
	res := l1
	for _, v := range l2 {
		res = append(res, v)
	}
	return res
}

//удаление повторов в массиве
func delPovtor2(l []ListUrl) []ListUrl {
	var f bool
	res := make([]ListUrl, 0)
	for i := 0; i < len(l); i++ {
		f = true
		for j := 0; j < i; j++ {
			if l[i].url == l[j].url { // && (l[i].fdownload == l[j].fdownload) {
				f = false
				break
			}
		}
		if f {
			res = append(res, l[i])
			f = true
		}
	}
	return res
}

//-----------------
// функция парсинга аргументов программы
func parse_args() bool {
	flag.StringVar(&site, "site", "", "Урл который нужно парсить для получения внутренних ссылок .")
	flag.StringVar(&koliter, "koliter", "", "Количества итераций для выкачивания .")
	flag.Parse()
	if site == "" {
		site = "http://echo.msk.ru"
	}
	if koliter == "" {
		koliter = "10"
	}
	return true
}

//-----------------

func main() {
	fmt.Println("Start Programm..")

	if !parse_args() {
		return
	}

	myurl := site
	ckoliter, _ := strconv.Atoi(koliter)
	timestart := time.Now().String()
	//	flagEnd := false // флаг окончания выгрузки
	lurl := make([]ListUrl, 0) // make([]string, 0)
	lurl = append(lurl, ListUrl{url: myurl, fdownload: 0})
	c := 0
	for {
		if (c == ckoliter) || (c > len(lurl)-1) {
			break
		} else {
			fmt.Print("c= ", c)
			body, errs := gethtmlpage(lurl[c].url)
			if errs {
				lurl[c].fdownload = 1 // обработан url
				listlinks := getLnksfromPage(body)
				listlinksurl := make([]ListUrl, 0)
				for _, vv := range listlinks {
					listlinksurl = append(listlinksurl, ListUrl{url: vv, fdownload: 0})
				}
				//			fmt.Println(listlinks)
				listnew := internalLinksfromSiteListUrl(listlinksurl, myurl)
				//				listnew2 := UniqLinks2(lurl, listnew)
				lurl = AddtoEndList2(lurl, listnew)
				lurl = delPovtor2(lurl)
				fmt.Println("   len(lurl)= ", len(lurl))
			} else {
				lurl[c].fdownload = -1 // обработан url была ошибка
			}
			c++
		}
	}
	fmt.Println("c= ", c)
	fmt.Println("len(lurl)= ", len(lurl))

	lurl = delPovtor2(lurl)
	fmt.Println("после удаления дубликатов - len(lurl)= ", len(lurl))

	s := ""
	for _, v := range lurl {
		s += v.url + ";" + strconv.Itoa(v.fdownload) + "\n"
	}
	timeend := time.Now().String()
	fmt.Println("Time Start...", timestart)
	fmt.Println("Time End...", timeend)
	fmt.Println("Start save result...")
	Savetofile("result.csv", s)
	fmt.Println("End save result...")
	fmt.Println("End Programm..")

	//	fmt.Println(lurl)
}
