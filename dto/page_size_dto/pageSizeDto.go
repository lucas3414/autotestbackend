package page_size_dto

type PageSizeDTO struct {
	Page int `json:"page,omitempty" uri:"page"`
	Size int `json:"size,omitempty" uri:"size"`
}

func (m *PageSizeDTO) GetPage() int {
	if m.Page <= 1 {
		m.Page = 1
	}
	return m.Page
}

func (m *PageSizeDTO) GetSize() int {
	if m.Size <= 0 {
		m.Size = 10
	}
	return m.Size
}
