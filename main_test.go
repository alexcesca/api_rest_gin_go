package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupRotasTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}
func TestVerificaStatusCodeSaudacao(t *testing.T) {
	r := SetupRotasTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/alex", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	mockDaResposta := `{"API diz:":"E ai alex, tudo beleza?"}`
	bodyRequest, _ := ioutil.ReadAll(resposta.Body)
	assert.Equal(t, mockDaResposta, string(bodyRequest))
}
func TestListaTodosAlunos(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	//fmt.Println(resposta.Body) // lista todos os alunos
}

func TestBuscaPorCpf(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678901", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	//fmt.Println(resposta.Body) // lista todos os alunos

}
func TestBuscaAlunoPorIdHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	pathBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", pathBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	assert.Equal(t, "Nome do aluno teste", alunoMock.Nome, "Os nomes deveriam ser iguais")
	assert.Equal(t, "12345678901", alunoMock.CPF, "Cpf esperado é 12345678901")
	assert.Equal(t, "123456789", alunoMock.RG, "RG do Aluno deveria ser 123456789")
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")

}
func TestDeletaAluno(t *testing.T) {

	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	r := SetupRotasTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	pathBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", pathBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusNoContent, resposta.Code, "Deveriam ser iguais")

}

func TestEditaUmAlunoHandler(t *testing.T) {

	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	aluno := models.Aluno{Nome: "Nome do aluno teste", CPF: "91945678901", RG: "123456700"}
	valorJaon, _ := json.Marshal(aluno)
	pathBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", pathBusca, bytes.NewBuffer(valorJaon))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	assert.Equal(t, aluno.Nome, alunoMock.Nome, "Os nomes deveriam ser iguais")
	assert.Equal(t, aluno.CPF, alunoMock.CPF, "CPF esperado é "+aluno.CPF)
	assert.Equal(t, aluno.RG, alunoMock.RG, "RG do Aluno deveria ser "+aluno.RG)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Nome do aluno teste", CPF: "12345678901", RG: "123456789"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}
func DeletaAlunoMock() {
	aluno := models.Aluno{}
	database.DB.Delete(&aluno, ID)

}
