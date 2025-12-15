package config

type Config struct {
	Server     ServerConfig     `toml:"SERVER"`
	Database   DatabaseConfig   `toml:"DATABASE"`
	Security   SecurityConfig   `toml:"SECURITY"`
	CORS       CORSConfig       `toml:"CORS"`
	FileUpload FileUploadConfig `toml:"FILEUPLOAD"`
	Pagination PaginationConfig `toml:"PAGINATION"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type DatabaseConfig struct {
	Driver string `toml:"driver"` // "sqlite", "postgres", or "mysql"

	// SQLITE
	Path string `toml:"path"`

	// MYSQL / POSTGRES
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`

	// POSTGRES ONLY
	SSLMode string `toml:"sslmode"`
}

type SecurityConfig struct {
	JWTSecret      string `toml:"jwt_secret"`
	JWTExpiryHours int    `toml:"jwt_expiry_hours"`
}

type CORSConfig struct {
	AllowedOrigins []string `toml:"allowed_origins"`
}

type FileUploadConfig struct {
	UploadDir        string   `toml:"upload_dir"`
	MaxUploadSizeMB  int      `toml:"max_upload_size_mb"`
	AllowedFileTypes []string `toml:"allowed_file_types"`
}

type PaginationConfig struct {
	DefaultPageSize int `toml:"default_page_size"`
	MaxPageSize     int `toml:"max_page_size"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Driver:   "sqlite",
			Path:     "./sustainwear.db",
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			DBName:   "sustainwear",
			SSLMode:  "disable",
		},
		Security: SecurityConfig{
			JWTSecret:      "secret-key",
			JWTExpiryHours: 24,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		},
		FileUpload: FileUploadConfig{
			UploadDir:        "./uploads",
			MaxUploadSizeMB:  10,
			AllowedFileTypes: []string{".jpg", ".png", ".jpeg", ".webp"},
		},
		Pagination: PaginationConfig{
			DefaultPageSize: 20,
			MaxPageSize:     100,
		},
	}
}
