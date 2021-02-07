package model

import (
	_ "sync"
)

// TTT int枚举测试1
type TTT int

const (
	// A1 a1的注释
	A1 TTT = 5

	// A2 a2的注释
	A2

	// A3 a3的注释
	A3
)

// XXX int枚举2
type XXX int

const (
	// X1 x1的注释
	X1 XXX = iota + 2
	// X2 x2的注释
	X2
)

// STR 字符串枚举类型
type STR string

const (
	// STR1 str1的注释
	STR1 STR = "aa"
	// STR2 str2的注释
	STR2 STR = "bb"
)

// dfdkhfd111111
type DDD struct {
	C1  string `json:"c1" bson:"c1"`
	C2  bool   `json:"c2" bson:"c2"`
	DDD []DDD  `bson:"dd"`
}

// Model dflndrleg;ldfjg
type Model struct {
	T1    TTT    `bson:"t1"` // xxx状态
	Model []*DDD `bson:"ddds"`
	STRS  STR    `bson:"str"` // yyy取值

	FFF struct {
		MM []int    `json:"mm" bson:"mm"`
		NN []string `json:"nn" bson:"nn"`
	} `json:"fff" bson:"fff"`

	// QQQ map[string]DDD `json:"qqq" bson:"qqq"`
}

func aa() {
	// var dd pkg.Account
	// fmt.Println(dd)
}
