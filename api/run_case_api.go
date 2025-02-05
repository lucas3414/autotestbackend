package api

import (
	"github.com/gin-gonic/gin"
	"go-gin-demo/dto/run_case_dto"
	"go-gin-demo/service/run_case_service"
)

type RunCaseApi struct {
	BaseApi
	RunCaseService *run_case_service.RunCasService
}

func NewRunCaseApi() RunCaseApi {
	return RunCaseApi{
		BaseApi:        NewBaseApi(),
		RunCaseService: run_case_service.NewRunCasService(),
	}
}

func (m RunCaseApi) RunCase(c *gin.Context) {

	var iRunCaseDTO dto.RunCaseDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iRunCaseDTO}).GetError(); err != nil {
		return
	}
	res := m.RunCaseService.RunCase(iRunCaseDTO)

	//if err != nil {
	//	m.Fail(ResponseJson{
	//		Msg: err.Error(),
	//	})
	//	return
	//}

	m.OK(ResponseJson{
		Data: gin.H{
			"run": res,
		},
	})
}
