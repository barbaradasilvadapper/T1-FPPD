// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
//
// ASSUNTO - Compreensão de concorrência e canais com buffer
//
// EXERCÍCIO:
//     1) Avalie o comportamento do programa para tamBuff
//        0 e 10.   Voce consegue explicar a diferença ?
// 		  Com buffer 0, as operações são sincronizadas e lentas, pois envio e recebimento
//		  precisam ocorrer simultaneamente, o que pode causar bloqueios e esperas.
//		  Com buffer 10, há mais concorrência e desempenho, pois produtor e consumidor 
// 	      trabalham de forma mais independente, já que a função de envio pode continuar
//		  mesmo que o consumidor esteja ocupado, até o limite do buffer.

//     2) Qual versao tem maior nivel de concorrencia ?
//        A versão com buffer maior (tamBuff = 10) tem maior nível de concorrência, pois
//        fonte e destino podem trabalhar mais independentemente.

//     3) Faça uma versão que tem vários processos destino
//        que podem consumir os dados de forma não determinística.
//        Ou seja, processos diferentes podem consumir quantidades
//        diferentes de itens,  conforme sua velocidade.
//        Como você coordenaria o término dos processos depois do
//        consumo dos N valores ?


// CÓDIGO ORIGINAL (1 e 2)
// package main

// const N = 100
// const tamBuff = 0

// func fonteDeDados(saida chan int) {
// 	for i := 1; i < N; i++ {
// 		println(i, " -> ")
// 		saida <- i
// 	}
// }

// func destinoDosDados(entrada chan int) {
// 	for i := 1; i < N; i++ {
// 		v := <-entrada
// 		println("                  -> ", v)
// 	}
// }

// func main() {
// 	c := make(chan int, tamBuff)
// 	go fonteDeDados(c)
// 	destinoDosDados(c)
// }


// CÓDIGO MODIFICADO (3)
// O produtor gera os números de 1..N e fecha o canal.
// Os consumidores rodam em paralelo e pegam valores na ordem em que conseguirem.
// Como o consumo é concorrente, alguns consumidores podem pegar mais valores que outros.
// O WaitGroup coordena o fim dos consumidores (cada um termina quando o canal fecha)
package main

import (
	"fmt"
	"sync"
)

const N = 100
const tamBuff = 0

func fonteDeDados(saida chan int) {
	for i := 1; i <= N; i++ {
		fmt.Println(i, " -> ")
		saida <- i
	}
	close(saida) // fecha o canal para sinalizar fim dos dados
}

func destinoDosDados(id int, entrada chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range entrada { // lê até o canal fechar
		fmt.Println("                  Consumidor", id, " -> ", v)
	}
}

func main() {
	c := make(chan int, tamBuff)
	var wg sync.WaitGroup

	// lança vários consumidores (destinos)
	numConsumidores := 5
	wg.Add(numConsumidores)
	for i := 1; i <= numConsumidores; i++ {
		go destinoDosDados(i, c, &wg)
	}

	// produtor (fonte)
	go fonteDeDados(c)

	// espera todos consumidores terminarem
	wg.Wait()
	fmt.Println("Fim")
}
