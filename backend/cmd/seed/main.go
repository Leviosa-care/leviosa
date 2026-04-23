package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
)

type SeedConfig struct {
	DatabaseURL     string
	DatabaseSchema  string
	VaultAddr       string
	VaultToken      string
	EncxKEKAlias    string
	EncxPepperAlias string
	SeedDataPath    string
}

func loadConfig() (*SeedConfig, error) {
	if err := godotenv.Load("development.env"); err != nil {
		log.Printf("Warning: Could not load development.env: %v", err)
	}

	cfg := &SeedConfig{
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		DatabaseSchema:  getEnv("DATABASE_SCHEMA", "auth"),
		VaultAddr:       getEnv("VAULT_ADDR", "http://localhost:8200"),
		VaultToken:      getEnv("VAULT_TOKEN", ""),
		EncxKEKAlias:    getEnv("ENCX_KEK_ALIAS", "leviosa-kek"),
		EncxPepperAlias: getEnv("ENCX_PEPPER_ALIAS", "leviosa"),
		SeedDataPath:    getEnv("SEED_DATA_PATH", "cmd/seed/seed_data.json"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.VaultToken == "" {
		return nil, fmt.Errorf("VAULT_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupCrypto(cfg *SeedConfig) (encx.CryptoService, error) {
	keyProvider, err := hashicorpkeys.NewTransitService()
	if err != nil {
		return nil, fmt.Errorf("create hashicorp key provider: %w", err)
	}

	secretProvider, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		return nil, fmt.Errorf("create hashicorp secret provider: %w", err)
	}

	crypto, err := encx.NewCrypto(context.Background(), keyProvider, secretProvider, encx.Config{
		KEKAlias:    cfg.EncxKEKAlias,
		PepperAlias: cfg.EncxPepperAlias,
	})
	if err != nil {
		return nil, fmt.Errorf("create crypto service: %w", err)
	}

	return crypto, nil
}

type AdminUser struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Role       string `json:"role"`
	Telephone  string `json:"telephone"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
}

func loadAdmins(path string) ([]AdminUser, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read seed data file %q: %w", path, err)
	}
	var admins []AdminUser
	if err := json.Unmarshal(data, &admins); err != nil {
		return nil, fmt.Errorf("parse seed data file: %w", err)
	}
	return admins, nil
}

func seedAdminUser(ctx context.Context, db *pgxpool.Pool, crypto encx.CryptoService, schema string, admin AdminUser) error {
	now := time.Now()

	user := &domain.User{
		ID:         uuid.New(),
		State:      domain.Active,
		Email:      admin.Email,
		Password:   admin.Password,
		FirstName:  admin.FirstName,
		LastName:   admin.LastName,
		Role:       admin.Role,
		Telephone:  admin.Telephone,
		PostalCode: admin.PostalCode,
		City:       admin.City,
		Address1:   admin.Address1,
		Address2:   admin.Address2,
		CreatedAt:  now,
		LoggedInAt: now,
	}

	userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
	if err != nil {
		return fmt.Errorf("encrypt user data: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.users (
			id, state,
			email_encrypted, email_hash,
			password_hash_secure,
			first_name_encrypted,
			last_name_encrypted,
			role_encrypted,
			telephone_encrypted, telephone_hash,
			postal_code_encrypted,
			city_encrypted,
			address1_encrypted,
			address2_encrypted,
			created_at_encrypted,
			logged_in_at_encrypted,
			picture_encrypted,
			birth_date_encrypted,
			gender_encrypted,
			google_id_encrypted,
			apple_id_encrypted,
			stripe_customer_id_encrypted,
			dek_encrypted, key_version
		) VALUES (
			$1, $2,
			$3, $4,
			$5,
			$6,
			$7,
			$8,
			$9, $10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23, $24
		)
		ON CONFLICT (email_hash) DO UPDATE SET
			email_encrypted = EXCLUDED.email_encrypted,
			password_hash_secure = EXCLUDED.password_hash_secure,
			first_name_encrypted = EXCLUDED.first_name_encrypted,
			last_name_encrypted = EXCLUDED.last_name_encrypted,
			role_encrypted = EXCLUDED.role_encrypted,
			telephone_encrypted = EXCLUDED.telephone_encrypted,
			telephone_hash = EXCLUDED.telephone_hash,
			postal_code_encrypted = EXCLUDED.postal_code_encrypted,
			city_encrypted = EXCLUDED.city_encrypted,
			address1_encrypted = EXCLUDED.address1_encrypted,
			address2_encrypted = EXCLUDED.address2_encrypted,
			created_at_encrypted = EXCLUDED.created_at_encrypted,
			logged_in_at_encrypted = EXCLUDED.logged_in_at_encrypted,
			dek_encrypted = EXCLUDED.dek_encrypted,
			key_version = EXCLUDED.key_version
	`, schema)

	_, err = db.Exec(ctx, query,
		userEncx.ID, userEncx.State,
		userEncx.EmailEncrypted, userEncx.EmailHash,
		userEncx.PasswordHashSecure,
		userEncx.FirstNameEncrypted,
		userEncx.LastNameEncrypted,
		userEncx.RoleEncrypted,
		userEncx.TelephoneEncrypted, userEncx.TelephoneHash,
		userEncx.PostalCodeEncrypted,
		userEncx.CityEncrypted,
		userEncx.Address1Encrypted,
		userEncx.Address2Encrypted,
		userEncx.CreatedAtEncrypted,
		userEncx.LoggedInAtEncrypted,
		userEncx.PictureEncrypted,
		userEncx.BirthDateEncrypted,
		userEncx.GenderEncrypted,
		userEncx.GoogleIDEncrypted,
		userEncx.AppleIDEncrypted,
		userEncx.StripeCustomerIDEncrypted,
		userEncx.DEKEncrypted, userEncx.KeyVersion,
	)
	if err != nil {
		return fmt.Errorf("insert user into database: %w", err)
	}

	log.Printf("Seeded admin user: %s (%s %s)", admin.Email, admin.FirstName, admin.LastName)
	return nil
}

func main() {
	ctx := context.Background()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	crypto, err := setupCrypto(cfg)
	if err != nil {
		log.Fatalf("Failed to setup crypto: %v", err)
	}

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	admins, err := loadAdmins(cfg.SeedDataPath)
	if err != nil {
		log.Fatalf("Failed to load admin users: %v", err)
	}

	log.Println("Starting database seeding...")
	log.Println("====================================")

	for _, admin := range admins {
		if err := seedAdminUser(ctx, db, crypto, cfg.DatabaseSchema, admin); err != nil {
			log.Fatalf("Failed to seed admin user %s: %v", admin.Email, err)
		}
	}

	log.Println("====================================")
	log.Println("Database seeding completed successfully!")
}
