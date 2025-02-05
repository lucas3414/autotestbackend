package run_case_service

import (
	dto "go-gin-demo/dto/run_case_dto"
	"go-gin-demo/service/base"
	"go-gin-demo/utils"
)

var runCasService *RunCasService

type RunCasService struct {
	base.BaseService
}

func NewRunCasService() *RunCasService {
	if runCasService == nil {
		runCasService = &RunCasService{
			//Dao: item_dao.NewItemDao(),
		}
	}
	return runCasService
}

func (m *RunCasService) RunCase(iRunCase dto.RunCaseDTO) []map[string]any {
	return utils.RunCase(iRunCase.CaseList)
}
