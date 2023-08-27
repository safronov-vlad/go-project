package main


import (
    "os"
    "fmt"
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/django/v3"
    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5/pgxpool"
)

func init() {
    if err := godotenv.Load(); err != nil {
        fmt.Errorf("No .env file found")
    }
}

func main() {
    // шаблонизатор Django
    engine := django.New("./templates", ".django")
    // Не совсем понял для чего TODO прочитать про контексты
    ctx := context.Background()
    // Подключение к БД
    dbUrl, exists := os.LookupEnv("DATABASE_URL")
    if exists == false{
        fmt.Println("No database config string")
	    os.Exit(1)
    }
	dbpool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

    // Запрос в бд
	rows, err := dbpool.Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	// обработка результатов
	var tables []string
    var table_name string
    for rows.Next() {
        err := rows.Scan(&table_name)
        if err != nil {
            panic(err)
        }
        tables = append(tables, table_name)
    }

    // Запуск приложения
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    // Обработчик
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("index", fiber.Map{
            "tables": tables,
        })
    })

    app.Listen(":3000")
}
