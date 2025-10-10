package merchant

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/infrastructure"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/handler"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/repository"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Routes struct {
	handler *handler.MerchantHandler
}

func NewRoutes(pool *pgxpool.Pool, minioClient *infrastructure.MinioClient, jwtSecret string) *Routes {
	repo := repository.NewMerchantRepositoryImpl(pool)
	svc := service.NewMerchantService(repo)
	h := handler.NewMerchantHandler(svc)

	return &Routes{
		handler: h,
	}
}

func (r *Routes) Register(router *gin.Engine, jwtSecret string) {
	SetupRoutes(router, r.handler, jwtSecret)
}
