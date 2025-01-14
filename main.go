package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Car представляет сущность автомобиля
type Car struct {
	ID         int64  `json:"id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Mileage    int64  `json:"mileage"`
	OwnerCount int    `json:"owner_count"`
}

// Furniture представляет сущность мебели
type Furniture struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Producer string `json:"producer"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Length   int    `json:"length"`
}

// Flower представляет сущность цветочной базы
type Flower struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Arrival  string  `json:"arrival"`
}

// Server хранит данные об объектах
type Server struct {
	carsID      int64
	furnitureID int64
	flowersID   int64
	cars        []Car
	furniture   []Furniture
	flowers     []Flower
	mutex       sync.Mutex
}

// HomePage отображает главную страницу
func HomePage(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Welcome</title>
	</head>
	<body>
		<h1>REST API</h1>
		<p>Choose an option:</p>
		<ul>
			<li><a href="/cars">Cars</a></li>
			<li><a href="/furniture">Furniture</a></li>
			<li><a href="/flowers">Flowers</a></li>
		</ul>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

// CreateCar создает новый автомобиль
func (s *Server) CreateCar(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	if err := json.NewDecoder(r.Body).Decode(&newCar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.carsID++
	newCar.ID = s.carsID
	s.cars = append(s.cars, newCar)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCar)
}

// GetCars возвращает список автомобилей
func (s *Server) GetCars(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.cars)
}

// CreateFurniture создает новый объект мебели
func (s *Server) CreateFurniture(w http.ResponseWriter, r *http.Request) {
	var newFurniture Furniture
	if err := json.NewDecoder(r.Body).Decode(&newFurniture); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.furnitureID++
	newFurniture.ID = s.furnitureID
	s.furniture = append(s.furniture, newFurniture)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newFurniture)
}

// GetFurniture возвращает список мебели
func (s *Server) GetFurniture(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.furniture)
}

// CreateFlower создает новую запись о цветах
func (s *Server) CreateFlower(w http.ResponseWriter, r *http.Request) {
	var newFlower Flower
	if err := json.NewDecoder(r.Body).Decode(&newFlower); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.flowersID++
	newFlower.ID = s.flowersID
	s.flowers = append(s.flowers, newFlower)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newFlower)
}

// GetFlowers возвращает список цветов
func (s *Server) GetFlowers(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.flowers)
}

func main() {
	server := &Server{}

	// Главная страница
	http.HandleFunc("/", HomePage)

	// Обработчики для автомобилей
	http.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			server.CreateCar(w, r)
		} else if r.Method == http.MethodGet {
			server.GetCars(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработчики для мебели
	http.HandleFunc("/furniture", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			server.CreateFurniture(w, r)
		} else if r.Method == http.MethodGet {
			server.GetFurniture(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработчики для цветов
	http.HandleFunc("/flowers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			server.CreateFlower(w, r)
		} else if r.Method == http.MethodGet {
			server.GetFlowers(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
