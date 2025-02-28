package upc 

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

type Product struct {
	Success bool `json:"success"`
	Barcode string `json:"barcode"`
	Title string `json:"title"`
	Alias string `json:"alias"`
	Description string `json:"description"`
	Brand string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
	MPN string `json:"mpn"`
	MSRP string `json:"msrp"`
	ASIN string `json:"ASIN"`
	Category string `json:"category"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (env *env) GetProductByCode(code string) (Product, error) {

	env.logger.Info("Getting product by UPC code")

	product := Product{}

	//WARNING: using the api key like this is not OAuth compliant but so far I cannot get it to work the receommended way
	requestURL := fmt.Sprintf("%sproduct/%s?apikey=%s", env.client.Endpoint, code, env.client.Key)
	//requestURL := fmt.Sprintf("%sproduct/%s", env.client.Endpoint, code)
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return product, err
	}
	
	//bearerToken := fmt.Sprintf("Bearer %s", env.client.Key)
	//request.Header.Set("authorization", bearerToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return product, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return product, err
	}

	sb := string(body)
	fmt.Println(sb)

	err = json.Unmarshal(body, &product)
	if err != nil {
		return product, err
	}
	
	return product, nil
	
}