package controller

import (
	"crud/data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	//Pega a Requisição do Body ( Corpo )
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha na Requisicao"))
		return
	}
	//Converte o Usuario
	var usuario usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter"))
		return
	}

	//Abre a conexão com o banco
	db, erro := data.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no Banco"))
		return
	}
	defer db.Close()

	//Evita SqlInjection
	statement, erro := db.Prepare("INSERT INTO usuarios( nome, email ) values (? ,?)")
	if erro != nil {
		w.Write([]byte("Erro ao Criar o Statement"))
	}
	defer statement.Close()

	//Insere o Usuario no Banco de  Dados
	insercao, erro := statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao Executar a inserção"))
		return
	}

	//Retorna o Id do Usuario Inserido
	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter o ID"))
		return
	}
	//Retorna o Status Code
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! Id: %d", idInserido)))

}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {

	db, erro := data.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao Conectar no banco"))
	}

	defer db.Close()

	linhas, erro := db.Query("SELECT * FROM usuarios")
	if erro != nil {
		w.Write([]byte("Erro ao Buscar Dados no Banco"))
	}
	defer linhas.Close()

	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario

		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao Scnaear o usuarios"))
			return
		}
		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter para Json"))
	}

}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao Converter para Int"))
		return
	}

	db, erro := data.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao Conectar no Banco"))
		return
	}
	linha, erro := db.Query("SELECT * FROM usuarios WHERE id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar ID"))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao Scan Usuario"))
			return
		}

	}

	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario a json"))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao Converter para Int"))
		return
	}
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o Corpo"))
		return
	}

	var usuario usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario para Struct"))
		return
	}

	db, erro := data.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao Conectar no banco"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("UPDATE usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao Criar o Statement "))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(&usuario.Nome, &usuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao Atualizar o Usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao Conerter para INT"))
		return
	}

	db, erro := data.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao Conectar no Banco de Dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("DELETE FROM usuarios where id = ? ")
	if erro != nil {
		w.Write([]byte("Erro ao criar o Statement "))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao Deletar Usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
