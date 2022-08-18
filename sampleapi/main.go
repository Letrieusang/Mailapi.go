package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

type Post struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"body"`
	CreateDate string `json:"createdate"`
	ChangeDate string `json:"changedate"`
}

func main() {

	r := gin.Default()
	client := r.Group("/api")
	{
		client.GET("/story/", Readall)
		client.GET("/story/:id", Readone)
		client.POST("/story/create", Create)
		client.PATCH("/story/update/:id", Update)
		client.DELETE("/story/:id", Delete)
	}
	r.Run(":8080")
}

func DBConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:Iamspectre16@tcp(127.0.0.1:3306)/luanvan")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Readall(c *gin.Context) {
	db := DBConn()
	rows, err := db.Query("SELECT * FROM article ORDER BY createdate")
	if err != nil {
		c.JSON(500, gin.H{
			"messages": "Story not found",
		})
	}
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Id")
	f.SetCellValue("Sheet1", "B1", "Title")
	f.SetCellValue("Sheet1", "C1", "Body")
	f.SetCellValue("Sheet1", "D1", "CreateDate")
	f.SetCellValue("Sheet1", "E1", "ChangeDate")
	poster := []Post{}
	post := Post{}
	count := 2
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.ChangeDate)
		if err != nil {
			panic(err.Error())
		}
		poster = append(poster, post)
		A := fmt.Sprint("A", count)
		B := fmt.Sprint("B", count)
		C := fmt.Sprint("C", count)
		D := fmt.Sprint("D", count)
		E := fmt.Sprint("E", count)
		f.SetCellValue("Sheet1", A, post.Id)
		f.SetCellValue("Sheet1", B, post.Title)
		f.SetCellValue("Sheet1", C, post.Content)
		f.SetCellValue("Sheet1", D, post.CreateDate)
		f.SetCellValue("Sheet1", E, post.ChangeDate)
		count += 1
	}
	c.JSON(200, poster)
	defer db.Close()
	if err := f.SaveAs("simple.xlsx"); err != nil {
		log.Fatal(err)
	}
}

func Readone(c *gin.Context) {
	db := DBConn()
	rows, err := db.Query("SELECT id, title, body FROM article WHERE id = " + c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"messages": "Story not found",
		})
	}
	poster := []Post{}
	post := Post{}
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.ChangeDate)
		if err != nil {
			panic(err.Error())
		}
		poster = append(poster, post)
	}
	c.JSON(200, poster)
	defer db.Close()
}

func Create(c *gin.Context) {
	db := DBConn()
	type CreatePost struct {
		Title string `form:"title" json:"title" binding:"required"`
		Body  string `form:"body" json:"body" binding:"required"`
	}
	var json CreatePost
	if err := c.ShouldBindJSON(&json); err == nil {
		insPost, err := db.Prepare("INSERT INTO article(title,body) VALUES(?,?)")
		if err != nil {
			c.JSON(500, gin.H{
				"messages": err,
			})
		}

		insPost.Exec(json.Title, json.Body)
		c.JSON(200, gin.H{
			"messages": "inserted",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	defer db.Close()
}
func Update(c *gin.Context) {
	db := DBConn()
	type UpdatePost struct {
		Title string `form:"title" json:"title" binding:"required"`
		Body  string `form:"body" json:"body" binding:"required"`
	}
	var json UpdatePost
	if err := c.ShouldBindJSON(&json); err == nil {
		edit, err := db.Prepare("UPDATE article SET title=?, body=? WHERE id= " + c.Param("id"))
		if err != nil {
			panic(err.Error())
		}
		edit.Exec(json.Title, json.Body)
		c.JSON(200, gin.H{
			"messages": "edited",
		})
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer db.Close()
}
func Delete(c *gin.Context) {
	db := DBConn()

	delete, err := db.Prepare("DELETE FROM article WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delete.Exec(c.Param("id"))
	c.JSON(200, gin.H{
		"messages": "deleted",
	})

	defer db.Close()
}
