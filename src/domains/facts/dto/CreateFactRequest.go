package dto

import "github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/models"

type CreateFactRequest struct {
	Facts []models.Fact
}

type CreateFactResponse struct {
}
