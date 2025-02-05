package dto

type RunCaseDTO struct {
	CaseList []map[string]any `json:"case_list" binding:"required" message:"用例不能为空" required_err:"用例不能为空"`
}
