package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/go-sql-driver/mysql"
)

type changearticle struct {
	id    string
	title string
	body  string
}

type newarticle struct {
	title string
	body  string
}

func changeart(word []string) changearticle {
	var (
		newid    string
		newtitle string
		newbody  string
	)
	newid = word[2]
	i := 4
	for i = 4; i < len(word); i++ {
		if word[i] == "Body" {
			break
		}
		newtitle = fmt.Sprint(newtitle, word[i])
		newtitle = fmt.Sprint(newtitle, " ")
	}
	for {
		i += 1
		if i > len(word)-1 {
			break
		}
		newbody = fmt.Sprint(newbody, word[i])
		newbody = fmt.Sprint(newbody, " ")
	}
	article := changearticle{title: newtitle, body: newbody, id: newid}
	return article
}

func newart(word []string) newarticle {
	var (
		newtitle string
		newbody  string
	)
	i := 2
	for i = 2; i < len(word); i++ {
		if word[i] == "Body" {
			break
		}
		newtitle = fmt.Sprint(newtitle, word[i])
		newtitle = fmt.Sprint(newtitle, " ")
	}
	for {
		i += 1
		if i > len(word)-1 {
			break
		}
		newbody = fmt.Sprint(newbody, word[i])
		newbody = fmt.Sprint(newbody, " ")
	}
	article := newarticle{title: newtitle, body: newbody}
	return article
}

func main() {
	var count uint32
	words := []string{}
	c, err := client.DialTLS("imap.gmail.com"+":993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")
	defer c.Logout()

	db, err := sql.Open("mysql", "root:Iamspectre16@tcp(127.0.0.1:3306)/luanvan")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Login
	if err := c.Login("letrieusang@gmail.com", "piahmsjzsblanywi"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")
	for {
		mbox, err := c.Select("INBOX", false)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			count = mbox.Messages
		}
		ch := make(chan *imap.Message, 100)
		myseqset := new(imap.SeqSet)
		if count != mbox.Messages {
			count += 1
			myseqset.AddRange(count, mbox.Messages)
			c.Fetch(myseqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, ch)
			count = mbox.Messages
			mess := <-ch
			for _, literal := range mess.Body {
				reader := bufio.NewReader(literal)
				for {
					line, _ := reader.ReadString('\n')
					words = strings.Split(line, " ")
					for k, i := range words {
						if i == "Update" {
							article := changeart(words)
							_, err := db.Exec("UPDATE article SET title=?, body=? WHERE id=?; ", article.title, article.body, article.id)
							if err != nil {
								panic(err.Error())
							}
							log.Println("Updated")
							break
						} else if i == "Insert" {
							article := newart(words)
							_, err := db.Exec("INSERT INTO article(title,body) VALUES(?,?)", article.title, article.body)
							if err != nil {
								panic(err.Error())
							}
							log.Println("Inserted")
							break
						} else if i == "Delete" {
							_, err = db.Exec("DELETE FROM article WHERE id=?;", words[k+1])
							if err != nil {
								log.Fatal(err)
							}
							log.Println("Deleted")
						}
					}
					if line == "" {
						break
					}
				}
			}
		} else {
			count = mbox.Messages
			continue
		} //

	}
}
