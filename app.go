package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDb() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

type Todo struct {
	gorm.Model
	Title  string `json:"title" form:"title" binding:"required"`
	Status int    `json:"status" form:"status" binding:"required"`
}

func main() {
	// 连接数据库
	if err := initDb(); err != nil {
		fmt.Println("连接数据库失败", err)
		panic(err)
	}
	db.AutoMigrate(&Todo{})
	r := gin.Default()
	// 添加路由组
	todo_g := r.Group("api/v1")
	{
		todo_g.GET("/todo", getTodoHandler)
		todo_g.POST("/todo", createTodoHandler)
		todo_g.PATCH("/todo", updateTodoHandler)
		todo_g.DELETE("/todo", deleteTodoHandler)

	}
	r.Run(":8080")
}

func getTodoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "getTodoHandler")
}

func createTodoHandler(ctx *gin.Context) {
	var todo Todo
	if err := ctx.ShouldBind(&todo); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "参数错误",
		})
		return
	}
	if db.Create(&todo).Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "添加成功",
			"data":    todo,
		})
	}
}

func updateTodoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "updateTodoHandler")
}

func deleteTodoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "deleteTodoHandler")
}
