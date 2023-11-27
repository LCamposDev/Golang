package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

type album struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *gorm.DB

func main() {
	// Carregar arquivo .env
	err1 := godotenv.Load()
	if err1 != nil{
		log.Fatal("Arquivo .env não pôde ser carregado.")
	}

	var err error

	// Conectar ao banco de dados MySQL
	dsn := os.Getenv("DBCONNECT")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate cria automaticamente a tabela "albums" no banco de dados
	err = db.AutoMigrate(&album{})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:id", deleteAlbums)
	router.PATCH("/albums/:id", patchAlbums)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	var albums []album
	db.Find(&albums)
	c.IndentedJSON(200, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

		if err := c.BindJSON(&newAlbum); err != nil {
			return
		}

	// Criar um novo registro de álbum no banco de dados
	db.Create(&newAlbum)
	c.IndentedJSON(201, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	var album album
	id := c.Param("id")

	// Buscar álbum pelo ID no banco de dados
	result := db.First(&album, id)
	if result.Error != nil {
		c.IndentedJSON(404, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(200, album)
}

func deleteAlbums(c *gin.Context) {
	var album album
	id := c.Param("id")

	// Excluir álbum pelo ID no banco de dados.
	result := db.Delete(&album, id)
	if result.Error != nil{
		c.IndentedJSON(404, gin.H{"message": "album not found"})
		return 
	}
	c.IndentedJSON(200, gin.H{"message": "Album deleted!"})
}


func patchAlbums(c *gin.Context) {
    var updatedAlbum album
    id := c.Param("id")

    // Buscar o álbum pelo ID no banco de dados
    result := db.First(&updatedAlbum, id)
    if result.Error != nil {
        c.IndentedJSON(404, gin.H{"message": "album not found"})
        return
    }

    // Decodificar os dados JSON da solicitação para atualizar o álbum existente
    if err := c.BindJSON(&updatedAlbum); err != nil {
        return
    }

    // Atualizar o álbum no banco de dados
    db.Save(&updatedAlbum)
    c.IndentedJSON(200, updatedAlbum)
}