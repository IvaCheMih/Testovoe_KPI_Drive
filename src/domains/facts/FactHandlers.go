package facts

import (
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/dto"
	"github.com/gofiber/fiber/v2"
	"log"
)

type FactHandlers struct {
	factsService *FactsService
}

func CreateFactHandlers(factsService *FactsService) FactHandlers {
	return FactHandlers{
		factsService: factsService,
	}
}

// CreateFacts godoc
// @Summary create fact.
// @Description create fact.
// @Tags createFact
// @Accept json
// @Produce json
// @Param session body dto.CreateFactRequest true "request"
// @Success 200 {object} dto.CreateFactResponse
// @Router /CreateFacts/ [post]
func (f *FactHandlers) CreateFacts(c *fiber.Ctx) error {

	// получаем массив фактов:
	facts, err := dto.CreateRequestFact(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// передаём саммив фактов на слой сервисов
	err = f.factsService.SaveMany(facts.Facts)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(200)
}

// GetFact godoc
// @Summary get fact.
// @Description get fact.
// @Tags getFact
// @Accept json
// @Produce json
// @Param session body dto.GetFactRequest true "request"
// @Success 200 {object} dto.GetFactResponse
// @Router /GetFact/ [post]
func (f *FactHandlers) GetFact(c *fiber.Ctx) error {
	fact, err := dto.GetRequestFact(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	responseGetFact, err := f.factsService.Get(fact)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(responseGetFact)
}
