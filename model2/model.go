package model

type RRR struct {
	DDD int    `json:"dddd" bson:"ddd"`
	EEE string `json:"eee" bson:"eee"`
	FFF struct {
		AA int `json:"aa"` // 最里层注释
		BB string
		CC string `json:"cc" bson:"cc"` // 最里层注释2
	} `json:"fff" bson:"fff"`
}
