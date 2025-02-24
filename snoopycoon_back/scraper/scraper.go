package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Price struct {
	Amount json.Number `json:"amount"`
}

type Product struct {
	ID    interface{} `json:"id"`
	Title string      `json:"title"`
	Price Price       `json:"price"`
}

func ScrapeWallapop(query string) (int, float64, error) {
	encodedQuery := url.QueryEscape(query)
	fullURL := fmt.Sprintf("https://api.wallapop.com/api/v3/search?source=recent_searches&keywords=%s&longitude=-3.69196&latitude=40.41956", encodedQuery)

	fmt.Println("🔍 URL generada para Wallapop:", fullURL)
	totalItems := 0
	totalPrice := 0.0
	nextPage := ""

	client := &http.Client{}
	maxPages := 25
	currentPage := 0

	for {
		start := time.Now()
		reqURL := fullURL
		if nextPage != "" {
			reqURL = "https://api.wallapop.com/api/v3/search?next_page=" + nextPage
		}

		fmt.Printf("🌍 Haciendo request a (%d/%d): %s\n", currentPage+1, maxPages, reqURL)

		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			return 0, 0, err
		}

		// --- HEADERS Wallapop ---
		req.Header.Add("Accept", "application/json, text/plain, */*")
		req.Header.Add("Accept-Language", "es,es-419;q=0.9,en;q=0.8")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("DeviceOS", "0")
		req.Header.Add("MPID", "-3440385224273352815")
		req.Header.Add("Origin", "https://es.wallapop.com")
		req.Header.Add("Referer", "https://es.wallapop.com/")
		req.Header.Add("Sec-Fetch-Dest", "empty")
		req.Header.Add("Sec-Fetch-Mode", "cors")
		req.Header.Add("Sec-Fetch-Site", "same-site")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
		req.Header.Add("X-AppVersion", "84130")
		req.Header.Add("X-DeviceID", "e11cb7f9-292f-423b-b94b-c46eb5fb4188")
		req.Header.Add("X-DeviceOS", "0")
		req.Header.Add("sec-ch-ua", `"Not(A:Brand";v="99", "Google Chrome";v="133", "Chromium";v="133"`)
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", `"Windows"`)

		res, err := client.Do(req)
		if err != nil {
			return 0, 0, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("❌ Error HTTP en Wallapop:", res.StatusCode)
			body, _ := io.ReadAll(res.Body)
			fmt.Println("🔍 Respuesta de Wallapop:", string(body))
			return 0, 0, fmt.Errorf("error HTTP %d", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, 0, err
		}

		var data struct {
			Data struct {
				Section struct {
					Payload struct {
						Items    []Product `json:"items"`
						NextPage string    `json:"next_page"`
					} `json:"payload"`
				} `json:"section"`
			} `json:"data"`
			Meta struct {
				NextPage string `json:"next_page"`
			} `json:"meta"`
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("❌ Error al parsear JSON de Wallapop:", err)
			return 0, 0, err
		}

		for _, item := range data.Data.Section.Payload.Items {
			priceFloat, _ := item.Price.Amount.Float64()
			totalPrice += priceFloat
			totalItems++
		}

		fmt.Println("⏱ Tiempo de ejecución de esta página:", time.Since(start))

		if data.Meta.NextPage == "" || currentPage+1 >= maxPages {
			break
		}
		nextPage = data.Meta.NextPage
		currentPage++
	}

	var avgPrice float64
	if totalItems > 0 {
		avgPrice = totalPrice / float64(totalItems)
	}

	return totalItems, avgPrice, nil
}

func getVintedToken() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.vinted.es", nil)
	if err != nil {
		return "", fmt.Errorf("error creando la petición a Vinted: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en la solicitud a Vinted: %v", err)
	}
	defer resp.Body.Close()

	var accessToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token_web" {
			accessToken = cookie.Value
			break
		}
	}

	if accessToken == "" {
		return "", fmt.Errorf("no se encontró access_token_web en las cookies")
	}

	fmt.Println("🔑 Token de Vinted obtenido correctamente")
	return accessToken, nil
}

func ScrapeVintedAPI(query string) (int, float64, error) {
	token, err := getVintedToken()
	if err != nil {
		fmt.Println("❌ Error obteniendo el token de Vinted:", err)
		return 0, 0, err
	}

	encodedQuery := url.QueryEscape(query)
	baseURL := "https://www.vinted.es/api/v2/catalog/items"

	totalItems := 0
	totalPrice := 0.0
	timeParam := ""

	page := 1
	var totalPages int

	for {
		apiURL := fmt.Sprintf("%s?page=%d&per_page=96&search_text=%s%s", baseURL, page, encodedQuery, timeParam)
		fmt.Println("🔍 URL generada para Vinted API:", apiURL)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Println("❌ Error creando la petición HTTP:", err)
			return 0, 0, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Error en la solicitud HTTP a Vinted:", err)
			return 0, 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("❌ Error HTTP en Vinted:", resp.StatusCode)
			fmt.Println("🔍 Respuesta de Vinted:", string(bodyBytes))
			return 0, 0, fmt.Errorf("error HTTP %d", resp.StatusCode)
		}

		var body struct {
			Items      []Product `json:"items"`
			Pagination struct {
				CurrentPage int `json:"current_page"`
				TotalPages  int `json:"total_pages"`
				Time        int `json:"time"`
			} `json:"pagination"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			fmt.Println("❌ Error al parsear JSON de Vinted:", err)
			return 0, 0, err
		}

		if totalPages == 0 {
			totalPages = body.Pagination.TotalPages
			timeParam = fmt.Sprintf("&time=%d", body.Pagination.Time)
			fmt.Println("📌 Total de páginas a recorrer:", totalPages)
		}

		for _, item := range body.Items {
			priceFloat, err := strconv.ParseFloat(item.Price.Amount.String(), 64)
			if err != nil {
				fmt.Println("⚠️ No se pudo convertir el precio de:", item.Title, "Valor:", item.Price.Amount)
				continue
			}
			totalPrice += priceFloat
			totalItems++
		}

		fmt.Printf("✅ Página %d/%d procesada. Productos acumulados: %d\n", page, totalPages, totalItems)

		if page >= totalPages {
			break
		}
		page++
	}

	avgPrice := 0.0
	if totalItems > 0 {
		avgPrice = totalPrice / float64(totalItems)
	}

	fmt.Println("✅ Scraping Vinted API completado. Total productos:", totalItems, "Precio medio:", avgPrice)

	return totalItems, avgPrice, nil
}
