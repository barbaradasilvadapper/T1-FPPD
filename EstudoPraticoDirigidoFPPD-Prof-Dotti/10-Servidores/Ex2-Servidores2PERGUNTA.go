// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   - Deseja-se tratar os clientes concorrentemente, e nao sequencialmente. Como ficaria a solucao ?
//   - Quantos clientes podem estar sendo tratados concorrentemente ?
//   - Agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//

package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL  = 100
	Pool = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r)
	}
}

// ------------------------------------
// servidor
// thread de servico calcula a resposta e manda direto pelo canal de retorno informado pelo cliente
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq ", id)
	req.ch_ret <- req.v * 2
}

// ------------------------------------
// servidor sequencial
func servidorSeq(in chan Request) {
	for {
		req := <-in
		fmt.Println("                       trataReq ", req)
		req.ch_ret <- req.v * 2 // responde  ao cliente
	}
}

// ------------------------------------
// servidor concorrente sem limite

func servidorConcSemLimite(in chan Request) {
	// servidor fica em loop eterno recebendo pedidos e criando um processo concorrente para tratar cada pedido
	var j int = 0
	for {
		j++
		req := <-in
		go trataReq(j, req)
	}
}

// ------------------------------------
// servidor concorrente que recebe o seu limite limite

func servidorConc(in chan Request, canal chan struct{}) {
	var j int = 0
	for {
		j++
		req := <-in
		canal <- struct{}{} // nunca le o canal entaos seu limite sera fixo em oq e vier nesse caso o pool
		go trataReq(j, req)
	}
}

// ------------------------------------
// main

func main() {
	fmt.Println("------ Servidor Sequencial -------")
	serv_chan_seq := make(chan Request)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan_seq)
	}
	go servidorSeq(serv_chan_seq)

	// --------------------------------
	fmt.Println("------ Servidor Concorrente (sem limite) -------")
	serv_chan_conc := make(chan Request)
	go servidorConcSemLimite(serv_chan_conc)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan_conc)
	}

	// --------------------------------
	fmt.Println("------ Servidor Concorrente (com limite Pool=10) -------")
	serv_chan_pool := make(chan Request)
	canal := make(chan struct{}, Pool)
	go servidorConc(serv_chan_pool, canal)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan_pool)
	}

	// bloqueia main
	select {}
}
