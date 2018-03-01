package model

import "time"

const (
	Static LangType = "Static"
	Dynamic LangType = "Dynamic"
)

// 言語の型ありかなしかを表す
type LangType string


// プログラム言語を表すモデル
type ProgramLang struct {
	ID string `json:"id"`
	Name string `json:"name"`
	LangType LangType `json:"langType"`
	Description string `json:"description"`
	UpdatedAt time.Time `json:"updatedAt"`
}


