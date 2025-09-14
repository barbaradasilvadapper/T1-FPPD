// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// programa da internet

// EXERCICIOS:

//    1) rode o programa abaixo.
//       o que você conclui sobre a execução observada?
//		 só aparecer "hello" 5 vezes e "world" não apareceu nenhuma, isso
//	     indica que a goroutine "world" não teve tempo de ser executada antes
//       do término da função main, que encerra todo o programa, não há 
//       sincronização ou espera de término da goroutine.

package main

import (
	"fmt"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}
