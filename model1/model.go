package model

import (
	_ "sync"
)

type TTT int

// dddddd
const (
	// a1a1
	A1 TTT = 5

	// a2
	A2

	// dfdfd

	// a3
	A3
)

// STR 字符串枚举类型
type STR string

const (
	// STR1 str1
	STR1 STR = "aa"
	// STR2 str2
	STR2 STR = "bb"
)

// dfdkhfd111111
// type DDD struct {
// 	C1  string `json:"c1" bson:"c1"`
// 	C2  bool   `json:"c2" bson:"c2"`
// 	DDD []DDD  `bson:"dd"`
// }

// Model dflndrleg;ldfjg
type Model struct {
	T1 TTT `bson:"t1"`
	// Model []*DDD         `bson:"ddds"`
	// KKK   *modelfuck.RRR `json:"kkk" bson:"kkk"`

	// FFF struct {
	// 	MM []int    `json:"mm" bson:"mm"`
	// 	NN []string `json:"nn" bson:"nn"`
	// } `json:"fff" bson:"fff"`

	// QQQ map[string]DDD `json:"qqq" bson:"qqq"`
}

func aa() {
	// var dd pkg.Account
	// fmt.Println(dd)
}
