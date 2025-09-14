// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.


package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5           // numero de jogadores
const M = 4            // numero de cartas na mao

type carta string      // carta é um strirng

var ch [NJ]chan carta  // NJ canais de itens tipo carta  
var novaRodada chan bool //para ver se vai para proxima rodada
var jogadoresQueBaixaram chan int //salvo o id do jogador que vai baixar 

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, jogoComecou chan bool ) {
	mao := cartasIniciais    // estado local - as cartas na mao do jogador
	nroDeCartas := M         // quantas cartas ele tem 
    var cartaRecebida carta   // carta recebida é vazia

	<- jogoComecou //le o canal que o jogo comecou 

	baixaAntes := false // se o jogador já bateu
	
	for {
		if len(jogadoresQueBaixaram) != 0  && !baixaAntes {
			if len(jogadoresQueBaixaram) == NJ-1 {
				fmt.Printf("O Jogador %d foi o último a bater e ele perdeu o jogo.\n", id)
			}
			jogadoresQueBaixaram <- id
			return
		}
			if nroDeCartas == 5 {
			fmt.Println(id, " joga")
			aux := rand.Intn(nroDeCartas)
			cartaParaBaixar := mao[aux]
			if baixou(mao) {
				cartaParaBaixar = cartaDiferente(mao)
			}
			fmt.Printf("Jogador %d passou a carta %s\n adiante ", id, cartaParaBaixar)

			out <- cartaParaBaixar
			mao = append(mao[:aux], mao[aux+1:]...)
			fmt.Printf("Mão do jogador %d: %v\n", id, mao)
			nroDeCartas--

			if baixou(mao) && nroDeCartas == 4 {
				fmt.Printf("Jogador %d bate!\n", id)
				baixaAntes = true
				jogadoresQueBaixaram <- id
				novaRodada <- false
			}

			if nroDeCartas == 0 {
				novaRodada <- true
			}
		} else {
			select {
			case cartaRecebida = <-in:
				fmt.Printf("Jogador %d recebeu a carta %s\n", id, cartaRecebida)
				mao = append(mao, cartaRecebida)
				nroDeCartas++
				fmt.Printf("Nova mão do jogador %d: %v\n", id, mao)
			default: // wait
			}
		}
	}
}

func cartaDiferente (mao []carta) carta{
	//cria um dicionario carta e valor
	count := make(map[carta] int)
	//guarda dele a quantidade de vezes que cada carta aparec 
	for _,n  := range  mao{
		count[n]++
	}
	for _,n  := range  mao{
		//se for igual a 1 quer dizer que a carta so apareceu uma vez
		if count[n]== 1 {
			return n
		}
	}

	return carta(" ")
}

func baixou(mao []carta) bool{
	//conta quantas vezes cada carta aparece 
	//cria um dicionário onde a chave é a carta e o valor é quantas vezes ela apareceu.
	count := make(map[carta]int)
	for _, n := range mao{
		count[n]++
		if count[n] >= 4 {
			return true
		}
	}
	return  false
}




func main() {

	jogadoresQueBaixaram = make(chan int, 5)
	// cria canais de passagem de cartas
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}


	jogoComecou := make(chan bool)
	novaRodada = make(chan bool)

	// cria canais para bater ? 

	// baralho = cria um baralho com NJ*M cartas

	for i := 0; i < NJ; i++ {   // cria os NJ jogadores

		cartasEscolhidas := make([]carta, 4)
		for j := 0; j < 4; j++ {
			randomIndex := rand.Intn(4)
			cartasEscolhidas[j] = carta(rune('D' + randomIndex))
		}


		// cria jogador i conectado com i-1 e i+1, e com as cartas	
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas , jogoComecou) 


		/*
		Cada jogador é criado como uma goroutine.

Cada jogador se conecta ao próximo (ch[(i+1)%NJ] cria o anel circular).

Ele recebe 4 cartas iniciais (mas esse trecho está bugado, porque mistura Card e carta).
		*/

	}
	
	// escolhe um jogador j e escreve uma carta em seu canal de entrada

	for i := 0; i < NJ; i++ {
		jogoComecou <- true
	}


	CartaInicial := carta("Joker")
	ch[0] <- CartaInicial

	time.Sleep(1 * time.Second)

	// espera ate jogadores baterem no(s) canal(is) de batida
	// registra ordem de batida
}


