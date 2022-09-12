package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Payment struct {
	gorm.Model
	Title string
	Price int
	Day   string
}

type Result struct {
	Total int
}

// DB migration
func dbInit() {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbInit())")
	}
	defer db.Close()
	db.AutoMigrate(&Payment{})
}

// DB Create
func dbInsert(title string, price int, day string) {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbInsert())")
	}
	defer db.Close()
	db.Create(&Payment{Title: title, Price: price, Day: day})
}

// DB Update
func dbUpdate(id int, title string, price int, day string) {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbUpdate())")
	}
	defer db.Close()
	var payment Payment
	db.First(&payment, id)
	payment.Title = title
	payment.Price = price
	payment.Day = day
	db.Save(&payment)
}

// DB Delete
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbDelete())")
	}
	defer db.Close()
	var payment Payment
	db.First(&payment, id)
	db.Unscoped().Delete(&payment)
}

// DB All Get
func dbGetAll() []Payment {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetAll())")
	}
	defer db.Close()
	var payment []Payment
	db.Order("created_at desc").Find(&payment)
	return payment
}

// DB One Get
func dbGetOne(id int) Payment {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetOne())")
	}
	defer db.Close()
	var payment Payment
	db.First(&payment, id)
	return payment
}

func dbGetNum() int {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetNum())")
	}
	defer db.Close()
	var num int
	db.Table("payments").Count(&num)
	return num
}

func dbGetPrice() Result {
	db, err := gorm.Open("sqlite3", "payment.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetPrice())")
	}
	defer db.Close()
	var result Result
	db.Table("payments").Select("sum(price) as total").Scan(&result)
	return result
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(c *gin.Context) {
		books := dbGetAll()
		num := dbGetNum()
		sumPrice := dbGetPrice()
		c.HTML(200, "index.html", gin.H{"books": books, "num": num, "sumPrice": sumPrice.Total})
	})

	//Create
	router.POST("/new", func(c *gin.Context) {
		title := c.PostForm("title")
		p := c.PostForm("price")
		day := c.PostForm("day")
		price, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		dbInsert(title, price, day)
		c.Redirect(302, "/")
	})

	//Edit
	router.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		payment := dbGetOne(id)
		c.HTML(200, "edit.html", gin.H{"payment": payment})
	})

	//Update
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		title := c.PostForm("title")
		p := c.PostForm("price")
		day := c.PostForm("day")
		price, errPrice := strconv.Atoi(p)
		if errPrice != nil {
			panic(errPrice)
		}
		dbUpdate(id, title, price, day)
		c.Redirect(302, "/")
	})

	//delete
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		c.Redirect(302, "/")
	})

	//delete_confirm
	router.GET("/delete_confirm/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		payment := dbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"payment": payment})
	})

	router.Run()
}
