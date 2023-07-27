package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	Status bool   `gorm:"default:false" json:"status" form:"status" `
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
		todo_g.DELETE("/todo/:id", deleteTodoHandler)

	}
	r.Run(":8080")
}

func getTodoHandler(ctx *gin.Context) {
	var todos []Todo
	if err := db.Find(&todos).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务器错误",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "查询成功",
		"data":    todos,
	})
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
	var todo Todo
	if err := ctx.ShouldBind(&todo); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "参数错误",
		})
		return
	}
	if err := db.First(&Todo{}, todo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "未找到对应的更新项",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务器错误",
		})
		return
	}
	if err := db.Omit("created_at").Save(&todo).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务器错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "修改成功",
		"data":    todo,
	})
}

func deleteTodoHandler(ctx *gin.Context) {
	var todo Todo
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "参数错误",
		})
	}
	if err := db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "未找到对应的更新项",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务器错误",
		})
		return
	}
	if err := db.Delete(&Todo{}, todo.ID).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除失败，服务器错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "删除成功",
		"data":    todo,
	})
}
