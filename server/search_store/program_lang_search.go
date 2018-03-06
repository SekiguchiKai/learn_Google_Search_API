package search_store

import (
	"context"
	"github.com/SekiguchiKai/learn_Google_Search_API/server/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
	"net/http"
)

const (
	_ProgramLangIndex = "ProgramLangIndex"
)

// Search APIのIndexとcontext.Contextを保持する構造体
type ProgramLangSearch struct {
	Index *search.Index
	Ctx   context.Context
}

func NewProgramLangSearch(r *http.Request) (ProgramLangSearch, error) {
	ctx := appengine.NewContext(r)
	idx, err := search.Open(_ProgramLangIndex)
	if err != nil {
		return ProgramLangSearch{}, err
	}

	return ProgramLangSearch{Index: idx, Ctx: ctx}, nil
}

func (s ProgramLangSearch) GetProgramLangList(startID string, limit int, dst *model.ProgramLangOptionList) error {
	var opts *search.ListOptions
	newListOptions(startID, limit, opts)
	iterator := s.Index.List(s.Ctx, opts)
	var list []model.ProgramLang
	for {
		var prl model.ProgramLang
		if _, err := iterator.Next(&prl);err == search.Done {
			break
		} else if err != nil {
			return err
		}

		list = append(list, prl)

		if iterator.Cursor() == "" {
			dst.HasNext = false
		}
		dst.StartID = string(iterator.Cursor())
		dst.List = list
	}

	return nil
}

func (s ProgramLangSearch) GetProgramLang(id string, dst *model.ProgramLang) (bool,error) {
	if id == "" {
		return false, nil
	}

	if  err := s.Index.Get(s.Ctx, model.ProgramLang.ID, dst); err != nil {
		if err != search.ErrNoSuchDocument {
			return false, err
		}
		return false, nil
	}
	return true, nil

}

func (s ProgramLangSearch) ExistsProgramLang(id string) (bool, error) {
	var dst model.ProgramLang
	return s.GetProgramLang(id, &dst)
}

// Search APIにProgramLangを元にDocumentを格納する
func (s ProgramLangSearch) PutProgramLang(src model.ProgramLang) error {
	if _, err := s.Index.Put(s.Ctx, model.ProgramLang.ID, src); err != nil {
		return err
	}
	return nil
}

func (s ProgramLangSearch) SearchProgramLang(dst *model.ProgramLangSearchOptions, opts SearchOptionsWrapper, query string) error {
	 iterator := s.Index.Search(s.Ctx, query, &opts.SearchOptions)
	var list []model.ProgramLang
	 for {
	 	var prl model.ProgramLang
	 	if _, err := iterator.Next(&prl); err == search.Done {
			break
		} else if err != nil {
			return err
		}

		 list = append(list, prl)

		 if iterator.Cursor() == "" {
			 dst.HasNext = false
		 }
		 dst.Cursor = string(iterator.Cursor())
		 dst.List = list
	 }

	return nil
}

func (s ProgramLangSearch) DeleteProgramLang(id string) error {
	if  err := s.Index.Delete(s.Ctx, id); err != nil {
		return err
	}
	return nil
}



func newListOptions(startID string, limit int, opts *search.ListOptions) {
	opts.Limit = limit
	opts.StartID = startID
	opts.IDsOnly = false
}
