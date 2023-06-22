package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

// Model Note
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	//load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error ketika membaca .env: %s", err)
	}

	// //membuat router
	r := gin.Default()
	// koneksi database
	db := dbConnection()
	defer db.Close()

	//middleware utk mengizinkan akses dari luar (biar browser bisa akses)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	//controller
	r.GET("/", func(c *gin.Context) {
		//repository / database
		notes, err := getNotes(db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//response
		c.JSON(200, gin.H{
			"data": notes,
		})
	})

	//controller
	r.POST("/", func(c *gin.Context) {
		//binding & validasi request
		var note Note
		if err := c.ShouldBindJSON(&note); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//repository / database
		noteResult, err := createNote(db, note.Title, note.Content)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//response
		c.JSON(200, gin.H{
			"message": "berhasil menambahkan catatan",
			"data":    noteResult,
		})
	})

	r.PUT("/:id", func(c *gin.Context) {
		var err error

		//binding & validasi request
		var note Note
		if err := c.ShouldBindJSON(&note); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		note.ID, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//repository / database
		err = updateNote(db, note.ID, note.Title, note.Content)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		noteResult, err := getNoteByID(db, note.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//response
		c.JSON(200, gin.H{
			"message": "berhasil mengubah catatan",
			"data":    noteResult,
		})
	})

	r.DELETE("/:id", func(c *gin.Context) {
		var err error

		//binding & validasi request
		var note Note
		note.ID, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//repository / database
		err = deleteNotes(db, note.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//response
		c.JSON(200, gin.H{
			"message": "berhasil menghapus catatan",
		})
	})

	r.Run(":8080") // listen and serve on
}

// sambungan ke database
func dbConnection() *sql.DB {
	//koneksi ke database
	db, err := sql.Open("mysql", os.Getenv("DB_CONNECTION_URL"))
	if err != nil {
		panic(err)
	}

	return db
}

// repository
func createNote(db *sql.DB, title, content string) (*Note, error) {
	const createNoteQuery = `
    INSERT INTO notes (title, content) 
    VALUES (?, ?)
    RETURNING id, title, content
  `

	var note Note
	err := db.QueryRow(createNoteQuery, title, content).Scan(&note.ID, &note.Title, &note.Content) //queryRow kalo cuma 1 baris
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func getNoteByID(db *sql.DB, id int) (*Note, error) {
	const getNoteByIDQuery = `
    SELECT id, title, content FROM notes WHERE id = ?
  `

	var note Note
	err := db.QueryRow(getNoteByIDQuery, id).Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func getNotes(db *sql.DB) ([]Note, error) {
	const getNotesQuery = `
    SELECT id, title, content FROM notes
  `

	rows, err := db.Query(getNotesQuery) //qyery kalo bisa lebih dari 1 baris
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func updateNote(db *sql.DB, id int, title, content string) error {
	const updateNoteQuery = `
    UPDATE notes 
    SET title = ?, content = ? 
    WHERE id = ?
  `

	db.QueryRow(updateNoteQuery, title, content, id)

	return nil
}

func deleteNotes(db *sql.DB, id int) error {
	const deleteNoteQuery = `
    DELETE FROM notes WHERE id = ?
  `

	_, err := db.Exec(deleteNoteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// binding & validasi
// func bindNoteRequest(c *gin.Context) (Note, error) {
// 	var note Note
// 	if os.Getenv("REQUEST_TYPE") == "JSON" {
// 		if err := c.ShouldBindWith(&note, binding.JSON); err != nil {
// 			return note, err
// 		}

// 	} else if os.Getenv("REQUEST_TYPE") == "FORM" {
// 		if err := c.ShouldBindWith(&note, binding.Form); err != nil {
// 			return note, err
// 		}

// 	}

// 	return note, nil
// }
