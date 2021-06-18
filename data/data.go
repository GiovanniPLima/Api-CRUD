package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //Driver de Conexão com MySql
)

//Conecta com o Mysql
func Conectar() (*sql.DB, error) {
	stringConexao := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", stringConexao)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		fmt.Println("Erro na Connection String")
		return nil, erro
	}
	fmt.Println("Conexao com Sucesso")
	return db, nil
}
