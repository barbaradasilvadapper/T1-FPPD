// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// >>> Veja antes Ex4 desta série.
// EXERCICIOS:
//     1) esta é uma solução para a questão anterior?
//     Funcionou neste caso, mas depende do agendamento e da velocidade do sistema. 
//     O correto seria sincronizar explicitamente, por exemplo, usando um canal de 
//     confirmação ou sync.WaitGroup para garantir que todos os valores foram processados 
//     antes do programa terminar.

//     2) o que garante que todos valores serão lidos antes do programa acabar?
// 	   Não há nada que garanta que todos os 1000 valores enviados pelo main serão lidos 
// 	   antes do main terminar. O main apenas faz os ch <- i e depois manda um valor no quit. 
// 	   Assim que isso acontece, a goroutine pode sair antes de consumir todos os valores 
//     ainda presentes no canal ch.


package main

import "fmt"

func main() {
	ch := make(chan int)
	quit := make(chan bool)
	go shower(ch, quit)
	for i := 0; i < 1000; i++ {
		ch <- i
	}
	quit <- false // or true, does not matter
}

func shower(c chan int, quit chan bool) {
	for {
		select {
		case j := <-c:
			fmt.Printf("%d\n", j)
		case <-quit:
			break
		}
	}
}
