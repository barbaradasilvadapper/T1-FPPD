// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?
// Veja abaixo a resposta ...
//   quantos clientes podem estar sendo tratados concorrentemente ?
//
// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
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

// Esta é a thread de serviço (worker).
// Ela fica em um loop eterno, aguardando por requisições no canal de entrada.
func worker(id int, in chan Request) {
	fmt.Println("                                 Worker", id, "iniciando.")
	for req := range in {
		fmt.Println("                                 Worker", id, "tratando req", req.v)
		req.ch_ret <- req.v * 2 // Calcula e responde ao cliente
	}
}

// Servidor que inicializa um pool de workers para tratar as requisições.
func servidorComPool(in chan Request) {
	// Cria um pool de workers
	for i := 0; i < Pool; i++ {
		go worker(i, in)
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores com Pool de Workers -------")
	serv_chan := make(chan Request) // Canal por onde o servidor recebe os pedidos
	go servidorComPool(serv_chan)   // Lança o processo servidor com o pool
	for i := 0; i < NCL; i++ {      // Lança diversos clientes
		go cliente(i, serv_chan)
	}
	<-make(chan int) // Bloqueia para o programa não terminar
}
