package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Etkinlik veri modeli
type Event struct {
	ID       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	City     string `json:"city" binding:"required"`
	Capacity int    `json:"capacity"`
}

var (
	// Verileri bellek içinde tutmak için map kullanıyoruz[cite: 1].
	events = make(map[string]*Event)
	// Eşzamanlı (concurrent) okuma/yazma işlemleri için Mutex[cite: 1].
	mu sync.RWMutex
)

func main() {
	r := gin.Default()

	r.GET("/events", getEvents)
	r.POST("/events", createEvent)
	r.POST("/events/:id/register", registerEvent)

	r.Run(":8080") // Sunucuyu başlat
}

// 1. GET /events[cite: 1]
func getEvents(c *gin.Context) {
	city := c.Query("city") // Opsiyonel ?city= parametresi[cite: 1]

	mu.RLock()
	defer mu.RUnlock()

	var result []Event
	for _, e := range events {
		if city == "" || e.City == city {
			result = append(result, *e)
		}
	}

	c.JSON(http.StatusOK, result) // Başarılı Yanıt: 200 OK[cite: 1]
}

// 2. POST /events[cite: 1]
func createEvent(c *gin.Context) {
	var newEvent Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Eksik veya hatalı veri"})
		return
	}

	// capacity değeri 0 veya negatif gelirse reddedilmeli[cite: 1]
	if newEvent.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kapasite 0 veya negatif olamaz"}) // 400 Bad Request[cite: 1]
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Aynı ID ile kayıt var mı kontrolü (Ekstra güvenlik)
	if _, exists := events[newEvent.ID]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bu ID'ye sahip bir etkinlik zaten var"})
		return
	}

	events[newEvent.ID] = &newEvent
	c.JSON(http.StatusCreated, newEvent) // Başarılı Yanıt: 201 Created
}

// 3. POST /events/{id}/register[cite: 2]
func registerEvent(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	defer mu.Unlock()

	event, exists := events[id]
	if !exists {
		// Etkinlik bulunamazsa: 404 Not Found[cite: 2]
		c.JSON(http.StatusNotFound, gin.H{"error": "Etkinlik bulunamadı"})
		return
	}

	if event.Capacity <= 0 {
		// Kontenjan kalmamışsa: 409 Conflict ve açıklayıcı JSON mesajı[cite: 2]
		c.JSON(http.StatusConflict, gin.H{"error": "Etkinlik kontenjanı dolmuştur"}) 
		return
	}

	// İlgili etkinliğe kayıt alır ve capacity değerini 1 azaltır[cite: 2]
	event.Capacity--
	
	// Başarılı Yanıt: 200 OK[cite: 2]
	c.JSON(http.StatusOK, gin.H{
		"message": "Kayıt başarılı",
		"remaining_capacity": event.Capacity,
	})
}