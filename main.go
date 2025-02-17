package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/lucastomic/snoopycoon_back/database"
	"github.com/lucastomic/snoopycoon_back/domain"
	"github.com/lucastomic/snoopycoon_back/usecases"
)

func updateTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Actualizando topic...")
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}

	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = usecases.UpdateTopic(uint(idInt), topic)
	if err != nil {
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topic)
}
func deleteTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Borrando topic...")
	id := r.PathValue("id")
	fmt.Println("Borrando listener con ID " + id)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}
	err = usecases.DeleteTopic(uint(idInt))
	if err != nil {
		fmt.Println("Error al borrar listener: " + err.Error())
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func createTopic(w http.ResponseWriter, r *http.Request) {
	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Creando topic...")
	fmt.Println(topic)
	err := usecases.CreateTopic(&topic)
	if err != nil {
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

func getListeners(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Obteniendo listeners...")
	listeners, err := usecases.GetTopics()
	if err != nil {
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)
		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(listeners)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creando usuario...")
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
		errJSON := map[string]string{"error": err.Error()}
		s, _ := json.Marshal(errJSON)

		http.Error(w, string(s), http.StatusInternalServerError)
		return
	}
	setAuthCookie(&w, jwt, 30)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": jwt})
}

func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Iniciando sesion...")
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Ha iniciado sesiòn %s\n", reqBody.Email)
	jwt, err := usecases.SignIn(reqBody.Email, reqBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Authenticación exitosa")
	setAuthCookie(&w, jwt, 30)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": jwt})
}

func setAuthCookie(
	w *http.ResponseWriter,
	token string,
	authCookieExpirationInDays int,
) http.ResponseWriter {
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
	return *w
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from localhost:3000
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		// If you need to send cookies or other credentials, set the following as well
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
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
	database.DB.AutoMigrate(
		&domain.User{},
		&domain.Topic{},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth/signup", createUser)
	mux.HandleFunc("POST /api/auth/signin", signIn)
	mux.HandleFunc("GET /api/listeners", getListeners)
	mux.HandleFunc("POST /api/listeners", createTopic)
	mux.HandleFunc("DELETE /api/listeners/{id}", deleteTopic)
	mux.HandleFunc("PUT /api/listeners/{id}", updateTopic)

	// Apply CORS middleware to all routes
	handler := corsMiddleware(mux)

	fmt.Println("Levantando el servidor para Brunito en http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
