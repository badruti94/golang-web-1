package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type M map[string]interface{}

type student struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade int    `json:"grade"`
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_belajar_golang")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func sqlQuery() ([]student, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select * from tb_student")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []student

	for rows.Next() {
		var each = student{}
		var err = rows.Scan(&each.Id, &each.Name, &each.Age, &each.Grade)

		if err != nil {
			return nil, err
		}
		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func insertData(id string, name string, age string, grade string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("insert into tb_student values (?,?,?,?)", id, name, age, grade)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func updateData(id string, name string, age string, grade string) error {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("update tb_student set name = ?, age = ?, grade = ? where id = ?", name, age, grade, id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func deleteData(id string) error {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("delete from tb_student where id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func main() {

	r := echo.New()

	r.GET("/", func(ctx echo.Context) error {
		result, err := sqlQuery()

		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, result)
	})

	r.GET("/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	r.POST("/", func(c echo.Context) error {
		id := c.FormValue("id")
		name := c.FormValue("name")
		age := c.FormValue("age")
		grade := c.FormValue("grade")

		err := insertData(id, name, age, grade)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "insert success!")
	})

	r.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		name := c.FormValue("name")
		age := c.FormValue("age")
		grade := c.FormValue("grade")

		err := updateData(id, name, age, grade)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "update success!")
	})

	r.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")

		err := deleteData(id)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "delete success!")

	})

	r.Start(":9000")
}
