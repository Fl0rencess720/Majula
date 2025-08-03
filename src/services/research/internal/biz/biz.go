package biz

type DeepResearchRepo interface {
}

type DeepResearchUsecase struct {
	repo DeepResearchRepo
}

func NewDeepResearchUsecase(repo DeepResearchRepo) *DeepResearchUsecase {
	return &DeepResearchUsecase{repo: repo}
}
