package dto

// RegisterRequest menampung payload untuk POST /api/auth/register
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// LoginRequest menampung payload untuk POST /api/auth/login
type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
	IPAddress  string `json:"ip_address"`
}

// LoginResponse mengembalikan token akses dan detail pengalihan halaman dashboard
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DashboardURL string `json:"dashboard_url"` // Diambil dari roles.dashboard_url
}

// UpdateUserRequest menampung payload untuk PUT /api/auth/users/{id}
type UpdateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"` // Menggunakan pointer untuk mendeteksi kehadiran data
}
