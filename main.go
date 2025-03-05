package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// Общая функция для обновления объекта
func updateEntity[T any](w http.ResponseWriter, r *http.Request, entities *[]T, getID func(T) int64, setID func(*T, int64)) {
	var updatedEntity T
	if err := json.NewDecoder(r.Body).Decode(&updatedEntity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := getID(updatedEntity)

	for i, entity := range *entities {
		if getID(entity) == id {
			(*entities)[i] = updatedEntity
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedEntity)
			return
		}
	}

	http.Error(w, "Entity not found", http.StatusNotFound)
}

// Функция удаления сущности
func deleteEntity[T any](w http.ResponseWriter, r *http.Request, entities *[]T, getID func(T) int64) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, entity := range *entities {
		if getID(entity) == id {
			*entities = append((*entities)[:i], (*entities)[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Entity not found", http.StatusNotFound)
}

// Удаление автомобиля
func (s *Server) DeleteCar(w http.ResponseWriter, r *http.Request) {
	deleteEntity(w, r, &s.cars, func(c Car) int64 { return c.ID })
}

// Удаление мебели
func (s *Server) DeleteFurniture(w http.ResponseWriter, r *http.Request) {
	deleteEntity(w, r, &s.furniture, func(f Furniture) int64 { return f.ID })
}

// Удаление цветов
func (s *Server) DeleteFlower(w http.ResponseWriter, r *http.Request) {
	deleteEntity(w, r, &s.flowers, func(f Flower) int64 { return f.ID })
}

// Методы для работы с автомобилями
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

func (s *Server) GetCars(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.cars)
}

func (s *Server) UpdateCar(w http.ResponseWriter, r *http.Request) {
	updateEntity(w, r, &s.cars, func(c Car) int64 { return c.ID }, func(c *Car, id int64) { c.ID = id })
}

// Методы для работы с мебелью
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

func (s *Server) GetFurniture(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.furniture)
}

func (s *Server) UpdateFurniture(w http.ResponseWriter, r *http.Request) {
	updateEntity(w, r, &s.furniture, func(f Furniture) int64 { return f.ID }, func(f *Furniture, id int64) { f.ID = id })
}

// Методы для работы с цветами
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

func (s *Server) GetFlowers(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.flowers)
}

func (s *Server) UpdateFlower(w http.ResponseWriter, r *http.Request) {
	updateEntity(w, r, &s.flowers, func(f Flower) int64 { return f.ID }, func(f *Flower, id int64) { f.ID = id })
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
		} else if r.Method == http.MethodPut {
			server.UpdateCar(w, r)
		} else if r.Method == http.MethodDelete {
			server.DeleteCar(w, r)
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
		} else if r.Method == http.MethodPut {
			server.UpdateFurniture(w, r)
		} else if r.Method == http.MethodDelete {
			server.DeleteFurniture(w, r)
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
		} else if r.Method == http.MethodPut {
			server.UpdateFlower(w, r)
		} else if r.Method == http.MethodDelete {
			server.DeleteFlower(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
