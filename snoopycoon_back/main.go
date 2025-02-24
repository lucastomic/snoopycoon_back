package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lucastomic/snoopycoon_back/database"
	"github.com/lucastomic/snoopycoon_back/domain"
	"github.com/lucastomic/snoopycoon_back/scraper"
	"github.com/lucastomic/snoopycoon_back/usecases"
)

func updateTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔄 Actualizando topic...")

	// Obtener ID desde la URL
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = usecases.UpdateTopic(uint(idInt), topic)
	if err != nil {
		http.Error(w, `{"error":"Error al actualizar topic"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Topic actualizado:", idInt)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topic)
}

func deleteTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🗑️ Recibiendo solicitud DELETE en /api/listeners/{id}")

	// Obtener ID desde la URL
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	err = usecases.DeleteTopic(uint(idInt))
	if err != nil {
		http.Error(w, `{"error":"Error al borrar topic"}`, http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Eliminado correctamente:", idInt)
	w.WriteHeader(http.StatusNoContent)
}

func createTopic(w http.ResponseWriter, r *http.Request) {
	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := usecases.CreateTopic(&topic)
	if err != nil {
		http.Error(w, `{"error":"Error al crear topic"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

func getListeners(w http.ResponseWriter, _ *http.Request) {
	listeners, err := usecases.GetTopics()
	if err != nil {
		http.Error(w, `{"error":"Error al obtener listeners"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(listeners)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := usecases.CreateUser(domain.User{Email: reqBody.Email, Password: reqBody.Password})
	if err != nil {
		http.Error(w, `{"error":"Error al crear usuario"}`, http.StatusInternalServerError)
		return
	}

	setAuthCookie(&w, jwt, 30)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": jwt})
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := usecases.SignIn(reqBody.Email, reqBody.Password)
	if err != nil {
		http.Error(w, `{"error":"Error en autenticación"}`, http.StatusUnauthorized)
		return
	}

	setAuthCookie(&w, jwt, 30)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": jwt})
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, `{"error":"Se requiere un término de búsqueda"}`, http.StatusBadRequest)
		return
	}

	// 1) Llamar a Wallapop
	wallapopItems, wallapopAvg, errWallapop := scraper.ScrapeWallapop(query)
	// 2) Llamar a Vinted
	vintedItems, vintedAvg, errVinted := scraper.ScrapeVintedAPI(query)

	// Si una falla, devuelves error
	if errWallapop != nil || errVinted != nil {
		http.Error(w, `{"error":"Error al scrapear uno de los marketplaces"}`, http.StatusInternalServerError)
		return
	}

	// Construyes una respuesta con datos de ambos
	response := map[string]interface{}{
		"query": query,
		"wallapop": map[string]interface{}{
			"total_items":   wallapopItems,
			"average_price": wallapopAvg,
		},
		"vinted": map[string]interface{}{
			"total_items":   vintedItems,
			"average_price": vintedAvg,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func setAuthCookie(w *http.ResponseWriter, token string, authCookieExpirationInDays int) {
	expirationTime := time.Now().Add(time.Duration(authCookieExpirationInDays) * 24 * time.Hour)
	http.SetCookie(*w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.GetGormDB()
	if err != nil {
		panic(err)
	}
	database.DB = db
	database.DB.AutoMigrate(&domain.User{}, &domain.Topic{})

	// ✅ Configurar rutas correctamente con `mux.NewRouter()`
	r := mux.NewRouter()

	// 🔐 Rutas de autenticación
	r.HandleFunc("/api/auth/signup", createUser).Methods("POST")
	r.HandleFunc("/api/auth/signin", signIn).Methods("POST")

	// 📄 Rutas de búsqueda
	r.HandleFunc("/api/listeners", getListeners).Methods("GET")
	r.HandleFunc("/api/listeners", createTopic).Methods("POST")
	r.HandleFunc("/api/listeners/{id}", deleteTopic).Methods("DELETE")
	r.HandleFunc("/api/listeners/{id}", updateTopic).Methods("PUT")

	// 🛒 Ruta de scraping
	r.HandleFunc("/api/scrape", scrapeHandler).Methods("GET")

	// 🌍 Aplica CORS a todas las rutas
	handler := corsMiddleware(r)

	fmt.Println("✅ Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
