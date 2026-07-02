package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"bisnis-rinzi/gateway/middleware"
	"bisnis-rinzi/packages/backend/database/redis"
	"bisnis-rinzi/packages/backend/jwt"
	"bisnis-rinzi/packages/backend/logger"
)

// RegisterRoutes menggabungkan seluruh jalur endpoint dan mengarahkannya via Reverse Proxy.
func RegisterRoutes(redisClient *redis.RedisClient) http.Handler {
	mux := http.NewServeMux()

	// 1. Definisikan Alamat Target URL Layanan Internal (bisa dibaca dari .env)
	authURL := getEnvTarget("AUTH_SERVICE_URL", "http://localhost:8081")
	inventoryURL := getEnvTarget("INVENTORY_SERVICE_URL", "http://localhost:8082")
	cashURL := getEnvTarget("CASH_SERVICE_URL", "http://localhost:8083")
	posURL := getEnvTarget("POS_SERVICE_URL", "http://localhost:8084")
	rentalURL := getEnvTarget("RENTAL_SERVICE_URL", "http://localhost:8085")
	financeURL := getEnvTarget("FINANCE_SERVICE_URL", "http://localhost:8086")

	// 2. Buat Handler Reverse Proxy untuk Masing-masing Layanan
	proxyAuth := createReverseProxy(authURL)
	proxyInventory := createReverseProxy(inventoryURL)
	proxyCash := createReverseProxy(cashURL)
	proxyPOS := createReverseProxy(posURL)
	proxyRental := createReverseProxy(rentalURL)
	proxyFinance := createReverseProxy(financeURL)

	// 3. Mapping Rute Publik (Tanpa Validasi Token)
	mux.Handle("/api/auth/login", proxyAuth)
	mux.Handle("/api/auth/register", proxyAuth)
	mux.Handle("/api/auth/refresh-token", proxyAuth)
	
	// Endpoint Publik untuk akses media secara aman melalui proxy
	mux.Handle("/api/public/inventory/media/", proxyInventory)

	// 4. Mapping Rute Terproteksi (Wajib Melewati AuthMiddleware)
	// Memanfaatkan ServeMux Sub-Routing Pattern
	protectedMux := http.NewServeMux()

	// Satukan sisa rute auth, inventory, cash, pos, rental, dan finance ke sub-router terproteksi
	protectedMux.Handle("/api/auth/", proxyAuth)
	protectedMux.Handle("/api/inventory/", proxyInventory)
	protectedMux.Handle("/api/cash/", proxyCash)
	protectedMux.Handle("/api/pos/", proxyPOS)
	protectedMux.Handle("/api/rental/", proxyRental)
	protectedMux.Handle("/api/finance/", proxyFinance)

	// Terapkan AuthMiddleware murni hanya pada rute terproteksi
	mux.Handle("/api/", middleware.AuthMiddleware(protectedMux, redisClient))

	// 5. Bungkus Router Utama dengan Global Middleware (CORS & Logging)
	var finalHandler http.Handler = mux
	finalHandler = middleware.LoggerMiddleware(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)

	return finalHandler
}

// Fungsi pembantu untuk membuat objek Reverse Proxy dari string URL target
func createReverseProxy(target string) http.Handler {
	urlTarget, err := url.Parse(target)
	if err != nil {
		logger.Error("Gagal inisialisasi proxy target %s: %v", target, err)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(urlTarget)

	// Kustomisasi director untuk menyesuaikan header request asli dari client
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = urlTarget.Host

		// Opsional: Kirimkan hasil claims JWT ke header layanan internal agar sub-service tahu siapa yang melakukan request
		if claims, ok := req.Context().Value(middleware.UserContextKey).(*jwt.Claims); ok {
			req.Header.Set("X-User-Id", claims.UserID)
			req.Header.Set("X-User-Role", claims.Role)
			req.Header.Set("X-User-Name", claims.Username)
		}
	}

	return proxy
}

func getEnvTarget(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
