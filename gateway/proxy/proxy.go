package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/response"
)

type GatewayProxy struct {
	Targets map[string]string // Memetakan prefix path (misal "/api/auth") ke target URL service internal
}

func NewGatewayProxy() *GatewayProxy {
	return &GatewayProxy{
		Targets: make(map[string]string),
	}
}

// RegisterService mendaftarkan target alamat routing baru
func (gp *GatewayProxy) RegisterService(pathPrefix string, targetURL string) {
	gp.Targets[pathPrefix] = targetURL
	logger.Info("Proxy Terdaftar: Jalur %s -> Meneruskan ke %s", pathPrefix, targetURL)
}

// Handler meluncurkan sistem reverse proxy berdasarkan rute prefix yang cocok
func (gp *GatewayProxy) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var targetURL string

		// Cari kecocokan prefix path terdekat
		for prefix, target := range gp.Targets {
			if strings.HasPrefix(r.URL.Path, prefix) {
				targetURL = target
				break
			}
		}

		if targetURL == "" {
			response.WriteError(w, http.StatusNotFound, "Endpoint API tidak ditemukan di Gateway")
			return
		}

		remote, err := url.Parse(targetURL)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Kesalahan konfigurasi proxy target")
			return
		}

		// Konfigurasi Reverse Proxy bawaan Go
		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Ubah arah request tujuan asli ke target internal microservice
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = remote.Host

		proxy.ServeHTTP(w, r)
	})
}
