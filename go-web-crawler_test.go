// go-web-crawler_test
package main

import (
	"fmt"
	"testing"
)

//func UniqLinks(list1, list2 []string) []string

func TestUniqLinks(t *testing.T) {
	l1 := []string{"1", "2", "5"}
	l2 := []string{"10", "5", "4", "10"}
	res := []string{"10", "4", "010"}

	l := UniqLinks(l1, l2)
	fmt.Println("результат", l)
	if (l[0] != res[0]) && (l[1] != res[1]) && (l[2] != res[2]) {
		t.Error(l)
	}
}

//func main() {
//	fmt.Println("Hello World!")
//}

//if v != 1.5 {
//        t.Error("Expected 1.5, got ", v)
//    }
