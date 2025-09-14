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
//	 infintos, pois uma goroutine é criada para cada requisicao recebida, ou seja,
//   não há limite definido, então todos os clientes podem ser tratados ao mesmo tempo, 
//   dependendo apenas dos recursos do sistema (CPU/memória).

// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//   criando um canal bufferizado com capacidade 10

package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL = 100
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

// servidor que dispara threads de servico
func servidorConc(in chan Request) {
	// servidor fica em loop eterno recebendo pedidos e criando um processo concorrente para tratar cada pedido
	var j int = 0
	for {
		j++
		req := <-in
		go trataReq(j, req)
	}
}


// ------------------------------------
// servidor com limite de 10 clientes concorrentemente
func servidorConcLim(in chan Request, limite int) {
	cont := make(chan struct{}, limite) // canal bufferizado que contada as vagas disponiveis

	var j int = 0
	for {
		req := <-in
		j++

		// "adquire" uma vaga do canal
		// bloqueia se estiver cheio (limite de 10 atingido)
		cont <- struct{}{}

		// trata a requisição em goroutine
		go func(id int, r Request) {
			defer func() { <-cont }() // libera a vaga ao finalizar
			trataReq(id, r)
		}(j, req)
	}
}


// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	serv_chan := make(chan Request) // CANAL POR ONDE SERVIDOR RECEBE PEDIDOS
	//go servidorConc(serv_chan)      // LANÇA PROCESSO SERVIDOR

	go servidorConcLim(serv_chan, Pool) // limite de 10 clientes simultâneos
	for i := 0; i < NCL; i++ {      // LANÇA DIVERSOS CLIENTES
		go cliente(i, serv_chan)
	}
	<-make(chan int)
}
