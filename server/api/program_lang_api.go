package api

import (
	"github.com/SekiguchiKai/learn_Google_Search_API/server/model"
	"github.com/SekiguchiKai/learn_Google_Search_API/server/search_store"
	"github.com/SekiguchiKai/learn_Google_Search_API/server/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

const (
	_DefaultProgramLangLimit = 20
	_MaxProgramLangLimit     = 50
)

type programLangQueryParams struct {
	Limit  int
	Cursor string
}

func InitProgramLangAPI(g *gin.RouterGroup) {
	g.GET("/langList", getProgramLangList)
	g.GET("/lang/:id", getProgramLang)
	g.POST("/lang/new", createProgramLang)
	g.PUT("/lang/:id", updateProgramLang)
	g.DELETE("/lang/:id", deleteProgramLang)
}

//func searchProgramLangList(c *gin.Context) {
//	util.InfoLog(c, "searchProgramLangList is called")
//
//	params, err := newProgramLangQueryParam(c)
//	if err != nil {
//		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
//	}
//}

func getProgramLang(c *gin.Context) {
	util.InfoLog(c, "getProgramLang is called")

	id := getProgramLangID(c)
	if id == "" {
		util.RespondAndLog(c, http.StatusBadRequest, "id is required")
	}

	var prl model.ProgramLang
	s, err := search_store.NewProgramLangSearch(c.Request)
	if err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	if exists, err := s.GetProgramLang(id, &prl); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	} else if !exists {
		util.RespondAndLog(c, http.StatusBadRequest, "invalid id : %s", id)
		return
	}

	c.JSON(http.StatusOK, prl)
}

func getProgramLangList(c *gin.Context) {
	util.InfoLog(c, "getProgramLangList is called")

	params, err := newProgramLangQueryParam(c)
	if err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	if params.Limit <= 0 {
		params.Limit = _DefaultProgramLangLimit
	} else if params.Limit > _MaxProgramLangLimit {
		params.Limit = _MaxProgramLangLimit
	}

	s, err := search_store.NewProgramLangSearch(c.Request)
	if err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	var list model.ProgramLangOptionList
	if err := s.GetProgramLangList(params.Cursor, params.Limit, &list); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)

}

func createProgramLang(c *gin.Context) {
	util.InfoLog(c, "createProgramLang is called")
	var params model.ProgramLang
	// HTTPリクエストで受け取ったJSONを構造体にロードする
	if err := bindProgramLangFromJson(c, &params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	// 構造体のバリデーションを行う
	if err := validateParamsForProgramLang(params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	// IDを付与する
	prl := model.NewProgramLang(params)
	prl.UpdatedAt = time.Now().UTC()

	s, err := search_store.NewProgramLangSearch(c.Request)
	if err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	if exists, err := s.ExistsProgramLang(prl.ID); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	} else if exists {
		caution := "There is same ProgramLang"
		util.RespondAndLog(c, http.StatusBadRequest, caution)
	}

	if err := s.PutProgramLang(prl); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, nil)
}

func updateProgramLang(c *gin.Context) {
	util.InfoLog(c, "updateProgramLang is called")
	var params model.ProgramLang
	if err := bindProgramLangFromJson(c, &params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateParamsForProgramLang(params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedAt := time.Now().UTC()

	s, err := search_store.NewProgramLangSearch(c.Request)
	if err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	var source model.ProgramLang
	if exists, err := s.GetProgramLang(params.ID, &source); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	} else if !exists {
		util.RespondAndLog(c, http.StatusNotFound, "id = %s is not found", params.ID)
	}

	u := model.UpdatedProgramLang(source, params)
	u.UpdatedAt = updatedAt

	if err := s.PutProgramLang(u); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, nil)

}

func deleteProgramLang(c *gin.Context) {
	util.InfoLog(c, "deleteProgramLang is called")
	id := getProgramLangID(c)

	s, err := search_store.NewProgramLangSearch(c.Request)
	if err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	}

	var prl model.ProgramLang
	if exists, err := s.GetProgramLang(id, &prl); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	} else if !exists {
		util.RespondAndLog(c, http.StatusNotFound, "id = %s is not found", id)
	}

	if err := s.DeleteProgramLang(id); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, nil)

}

func newProgramLangQueryParam(c *gin.Context) (programLangQueryParams, error) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return programLangQueryParams{}, err
	}

	return programLangQueryParams{
		Limit:  limit,
		Cursor: c.Query("cursor"),
	}, nil

}

// HTTPのリクエストボディのjsonデータProgramLangに変換
func bindProgramLangFromJson(c *gin.Context, dst *model.ProgramLang) error {
	if err := c.BindJSON(dst); err != nil {
		return err
	}

	dst.ID = getProgramLangID(c)
	return nil
}

// IDを取得
func getProgramLangID(c *gin.Context) string {
	return c.Param("id")
}

func validateParamsForProgramLang(prl model.ProgramLang) error {
	if prl.Name == "" {
		return errors.New("name is required")
	}

	if prl.LangType == "" {
		return errors.New("langType is required")
	}

	return nil
}
