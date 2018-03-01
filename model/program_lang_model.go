package model

import "time"

const (
	Static LangType = "Static"
	Dynamic LangType = "Dynamic"
)

// 言語の型ありかなしかを表す
type LangType string

type ProgramLangList struct {
	ProgramLangOptionList ProgramLangOptionList `json:"programLangOptionList"`
	ProgramLangSearchOptions ProgramLangSearchOptions `json:"programLangSearchOptions"`
}

type ProgramLangOptionList struct {
	List []ProgramLang `json:"list"`
	HasNext bool `json:"hasNext"`
	StartID string `json:"startID"`
}

type ProgramLangSearchOptions struct {
	List []ProgramLang `json:"list"`
	HasNext bool `json:"hasNext"`
	Cursor string `json:"cursor"`
}

// プログラム言語を表すモデル
type ProgramLang struct {
	ID string `json:"id"`
	Name string `json:"name"`
	LangType LangType `json:"langType"`
	Description string `json:"description"`
	UpdatedAt time.Time `json:"updatedAt"`
}


