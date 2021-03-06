package model

import (
	"time"
	"github.com/SekiguchiKai/learn_Google_Search_API/server/util"
)

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


func NewProgramLang(param ProgramLang) ProgramLang {
	param.ID = newUserID(param.Name, param.LangType)
	return param
}

func UpdatedProgramLang(source, param ProgramLang) ProgramLang {
	param.Description = source.Description
	return param
}

func newUserID(name string, langType LangType) string {
	return util.GetHash(name + "@@" + string(langType))
}