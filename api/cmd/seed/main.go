package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"dac/project-tracker/internal/config"
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/infrastructure/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌱 Seed — Criar usuário no ProjectTracker")
	fmt.Println()

	name := prompt(reader, "Nome: ")
	email := prompt(reader, "Email: ")
	password := promptPassword("Senha: ")
	role := prompt(reader, "Role (admin/manager/member) [member]: ")
	if role == "" {
		role = "member"
	}

	cfg := &config.Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "project_tracker"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar no banco: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Falha ao gerar hash: %v", err)
	}

	user := model.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     role,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Falha ao criar usuário: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Usuário criado com sucesso!")
	fmt.Printf("   Nome:  %s\n", user.Name)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Role:  %s\n", user.Role)
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(label)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func promptPassword(label string) string {
	fmt.Print(label)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatalf("Falha ao ler senha: %v", err)
	}
	return string(bytePassword)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
