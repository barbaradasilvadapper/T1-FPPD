// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// >>> Veja o Ex0 desta série
// ABRE E FECHA CONCORRENCIA
// Há várias formas de esperar o término de processos concorrentes.

// EXERCICIOS:

//   1)  isto seria uma solução para sincronizar o final do programa ?
//       sim, pois o canal fin é usado para sinalizar quando cada goroutine 
// 		 termina. O main espera (<-fin) duas vezes, garantindo que ambas as goroutines
//       tenham concluído sua execução antes de o programa terminar.
//       obs: o canal fin é do tipo chan struct{}, que é um canal vazio usado apenas para sinalização.

//   2)  aumente para criar 10 processos concorrentes say(...).
//       como voce faz a espera de todos ?
// OBS:  tente um comando de repeticao.
//       um loop que inicia 10 goroutines e outro loop que espera 10 sinais no canal fin,
//       garantindo que todas as goroutines tenham concluído antes do programa terminar.

package main

import (
	"fmt"
)

func say(s string, c chan struct{}) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
	c <- struct{}{}
}

func main() {

	//CÓDIGO ORIGINAL (1)
	// fin := make(chan struct{})
	// go say("world", fin)
	// go say("hello", fin)
	// <-fin
	// <-fin

	//CÓDIGO MODIFICADO (2)
	fin := make(chan struct{})
    for i := 0; i < 10; i++ {
        go say(fmt.Sprintf("goroutine %d", i), fin)
    }
    for i := 0; i < 10; i++ {
        <-fin
    }
}
