package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func PerIPRateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.ClientIP())
		if err != nil {
			ip = c.ClientIP()
		}

		mu.Lock()
		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{limiter: rate.NewLimiter(r, b)}
			log.Printf("Creating new limiter for IP: %s", ip)
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		mu.Unlock()
		c.Next()
	}
}
