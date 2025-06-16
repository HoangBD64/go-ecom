package main

import (
	"log"

	"github.com/HoangBD64/go-ecom/pkg/config"
	"github.com/HoangBD64/go-ecom/pkg/di"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error to load the config: ", err)
	}

	server, err := di.InitializeApi(cfg)
	if err != nil {
		log.Fatal("Failed to initialize the api: ", err)
	}

	if server.Start(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}


package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/xuri/excelize/v2"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int
	Name  string
	Email string
}

const (
	dbUser     = "postgres"
	dbPassword = "your_password"
	dbName     = "test"
	dbHost     = "localhost"
	dbPort     = "5432"
)

func main() {
	http.HandleFunc("/export", handler)
	fmt.Println("Server chạy tại http://localhost:8080/export")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// =====================================
// 🧩 1. Goroutine đọc DB (stream)
// =====================================
func fetchUsersFromDB(ch chan<- User, errCh chan<- error) {
	defer close(ch)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		errCh <- err
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		errCh <- err
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			errCh <- err
			return
		}
		ch <- u // gửi từng dòng vào channel
	}
}

// =====================================
// 🧩 2. Goroutine ghi file Excel (stream)
// =====================================
func writeUsersToExcelStream(w http.ResponseWriter, ch <-chan User) error {
	f := excelize.NewFile()
	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}

	// Ghi header
	header := []interface{}{"ID", "Tên", "Email"}
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	sw.SetRow(cell, header)

	rowIndex := 2
	for user := range ch {
		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		row := []interface{}{user.ID, user.Name, user.Email}
		if err := sw.SetRow(cell, row); err != nil {
			return err
		}
		rowIndex++
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	return f.Write(w)
}

// =====================================
// 🧩 3. Handler: nối các goroutine
// =====================================
func handler(w http.ResponseWriter, r *http.Request) {
	// Setup header để tải file
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="users_streamed_channel.xlsx"`)

	// Channel truyền dữ liệu và error
	dataCh := make(chan User, 100)       // Buffered để tăng hiệu suất
	errCh := make(chan error, 1)

	// Goroutine đọc DB
	go fetchUsersFromDB(dataCh, errCh)

	// Ghi ra file Excel trực tiếp từ channel
	if err := writeUsersToExcelStream(w, dataCh); err != nil {
		log.Println("Lỗi ghi Excel:", err)
		http.Error(w, "Lỗi ghi Excel", 500)
		return
	}

	// Check error từ goroutine đọc DB
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("Lỗi khi đọc DB:", err)
			http.Error(w, "Lỗi đọc dữ liệu", 500)
			return
		}
	default:
		// Không có lỗi
	}
}
