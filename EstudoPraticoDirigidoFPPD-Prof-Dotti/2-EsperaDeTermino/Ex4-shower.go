// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Exemplo da Internet

// EXERCICIOS:
//   1) rode o programa abaixo e interprete.
//      todos os valores escritos no canal são lidos?
//	    nao, pois programa termina após o main finalizar o envio dos valores, 
// 		encerrando o programa sem a goroutine ter tempo de ler todos os valores do canal.

//   2) como isto poderia ser resolvido ?

package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	done := make(chan struct{}) // canal para sinalizar o término
	
	// inicia a goroutine que lê do canal
	go shower(ch, done)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch) // fecha o canal após enviar todos os valores
	<-done // espera a goroutine terminar
}

func shower(c chan int, done chan struct{}) {
	//CÓDIGO ORIGINAL (1)
	// for {
	// 	j := <-c
	// 	fmt.Printf("%d\n", j)
	// }

	//CÓDIGO MODIFICADO (2)
	for j := range c {
		fmt.Printf("%d\n", j)
	}
	// sinaliza que terminou
	done <- struct{}{}
}
