// go-web-crawler_test
package main

import (
	//	"fmt"
	"strings"
	"testing"
)

//func UniqLinks(list1, list2 []string) []string
// если массивы равны по длине и по символу то истина
func comparearray(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, _ := range l1 {
		if strings.Compare(l1[i], l2[i]) != 0 {
			return false
		}
	}
	return true
}

func TestUniqLinks(t *testing.T) {
	l1 := []string{"1", "2", "5"}
	l2 := []string{"10", "5", "4", "10"}
	res := []string{"10", "4", "10"}

	l := UniqLinks(l1, l2)
	if comparearray(l, res) == false {
		t.Error(l)
	}

	l1 = make([]string, 0)
	l2 = []string{"10", "5", "4", "10"}
	res = []string{"10", "5", "4", "10"}

	l = UniqLinks(l1, l2)
	if comparearray(l, res) == false {
		t.Error("10 5 4 10", l)
	}

	l1 = []string{"10", "5", "4", "10"}
	l2 = []string{"1", "2", "5"}
	res = []string{"1", "2"}

	l = UniqLinks(l1, l2)
	if comparearray(l, res) == false {
		t.Error(l)
	}

	l1 = []string{"10", "5", "4", "10"}
	l2 = []string{"10", "5", "4", "10"}
	res = make([]string, 0)

	l = UniqLinks(l1, l2)
	if comparearray(l, res) == false {
		t.Error(l)
	}
}

func TestAddtoEndList(t *testing.T) {
	l1 := []string{"1", "2", "5"}
	l2 := []string{"10", "5", "4", "10"}
	res := []string{"1", "2", "5", "10", "5", "4", "10"}
	l := AddtoEndList(l1, l2)
	if comparearray(l, res) == false {
		t.Error(l)
	}

	l1 = make([]string, 0)
	l2 = []string{"10", "5", "4", "10"}
	res = []string{"10", "5", "4", "10"}
	l = AddtoEndList(l1, l2)
	if comparearray(l, res) == false {
		t.Error(l)
	}

}

func TestDelPovtor(t *testing.T) {
	l1 := []string{"1", "2", "5", "2", "99", "", "0", "", "1"}
	res := []string{"1", "2", "5", "99", "", "0"}
	l := delPovtor(l1)
	if comparearray(l, res) == false {
		t.Error(l)
	}

	l1 = []string{"1", "1", "1", "1", "1"}
	res = []string{"1"}
	l = delPovtor(l1)
	if comparearray(l, res) == false {
		t.Error(l)
	}

	l1 = make([]string, 0)
	res = make([]string, 0)
	l = delPovtor(l1)
	if comparearray(l, res) == false {
		t.Error(l)
	}
}

//
