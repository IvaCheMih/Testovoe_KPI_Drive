package facts

import (
	"bytes"
	"fmt"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/dto"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/models"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/services"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type FactsService struct {
	r      *RabbitQueueService
	config services.ConfigurationService
}

func CreateFactsService(rabbitService *RabbitQueueService, configurationService services.ConfigurationService) FactsService {
	return FactsService{
		r:      rabbitService,
		config: configurationService,
	}
}

func (f *FactsService) Get(fact dto.GetFactRequest) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// записываем факт как form/data
	err := putGetFactInBody(writer, fact)
	if err != nil {
		return "", err
	}

	writer.Close()

	// создаём http запрос
	request, err := http.NewRequest(http.MethodPost, f.config.KpiGetApiURL, bytes.NewReader(body.Bytes()))

	if err != nil {
		fmt.Println(err)
	}

	//добавляем авторизацию и тип контента
	request.Header.Add("Authorization", "Bearer 48ab34464a5573519725deb5865cc74c")
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	// делаем запрос
	response, err := client.Do(request)

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("client: could not read response body: %s\n", err)
	}

	// в случае успеха печатаем полученный ответ
	if response.StatusCode == 200 {
		fmt.Println()
		log.Println("[SUCCESS]:")
		fmt.Println(string(resBody))
		fmt.Println()
	}

	return string(resBody), nil
}

func (f *FactsService) SaveMany(facts []models.Fact) error {
	// поштучно отправляем факты в сервис rabbit
	for _, fact := range facts {
		err := f.r.RabbitPublish(fact)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (f *FactsService) SaveOne(fact models.Fact) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// записываем факт как form/data
	err := putFactInBody(writer, fact)
	if err != nil {
		return err
	}

	writer.Close()

	// создаём http запрос

	request, err := http.NewRequest(http.MethodPost, f.config.KpiPostApiURL, bytes.NewReader(body.Bytes()))

	if err != nil {
		fmt.Println(err)
	}

	//добавляем авторизацию и тип контента
	request.Header.Add("Authorization", "Bearer 48ab34464a5573519725deb5865cc74c")
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	// делаем запрос
	response, err := client.Do(request)

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("client: could not read response body: %s\n", err)
	}

	// в случае успеха печатаем полученный ответ
	if response.StatusCode == 200 {
		fmt.Println()
		log.Println("[SUCCESS]:")
		fmt.Println(string(resBody))
		fmt.Println()
	}

	return err
}

func putFactInBody(writer *multipart.Writer, fact models.Fact) error {
	var factForm = make(map[string]string, 10)
	var err error

	factForm["period_start"] = fact.PeriodStart
	factForm["period_end"] = fact.PeriodEnd
	factForm["period_key"] = fact.PeriodKey
	factForm["indicator_to_mo_id"] = fact.IndicatorToMoId
	factForm["indicator_to_mo_fact_id"] = fact.IndicatorToMoFactId
	factForm["value"] = fact.Value
	factForm["fact_time"] = fact.FactTime
	factForm["is_plan"] = fact.IsPlan
	factForm["auth_user_id"] = fact.AuthUserId
	factForm["comment"] = fact.Comment

	for key, value := range factForm {
		err = put(writer, key, value)
	}

	return err
}

func putGetFactInBody(writer *multipart.Writer, fact dto.GetFactRequest) error {
	var factForm = make(map[string]string, 10)
	var err error

	factForm["period_start"] = fact.PeriodStart
	factForm["period_end"] = fact.PeriodEnd
	factForm["period_key"] = fact.PeriodKey
	factForm["indicator_to_mo_id"] = fact.IndicatorToMoId

	for key, value := range factForm {
		err = put(writer, key, value)
	}

	return err
}

func put(writer *multipart.Writer, key string, value string) error {
	fw, err := writer.CreateFormField(key)
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, strings.NewReader(value))

	return err
}
