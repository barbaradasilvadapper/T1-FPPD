// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// esta solucao utiliza workgroups de Golang.
// um wg é usado para esperar que um grupo de rotinas acabe.
// a solução é equivalente à anterior com canais.
// nesta disciplina, na parte que se refere uso de canais,
// o professor utilizará somente canais como forma de sincronizacao, e nao outras
// abstracoes da biblioteca da linguagem - todas elas podem ser reduzidas ao uso de canais.
// voce pode usar esta funcionalidade se desejar.

// obs: wg é uma variável do tipo sync.WaitGroup em Go, que serve para sincronizar
// a execução de várias goroutines, permitindo que o programa espere até que todas 
// terminem antes de continuar. 
// chama wg.Add(n) para indicar quantas goroutines vai esperar
// cada goroutine chama wg.Done() ao terminar
// main chama wg.Wait() para aguardar todas finalizarem

package main

import (
	"fmt"
	"sync"
)

func say(s string, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
	wg.Done()
}

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go say("world", &waitgroup)
	go say("hello", &waitgroup)
	waitgroup.Wait()
}
