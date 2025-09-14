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

const NJ = 5       // número de jogadores
const M = 4        // número de cartas na mão

type carta string  // carta é um strirng

// canais
var ch [NJ]chan carta    // NJ canais de itens tipo carta    
var novaRodada chan bool     
var canalBater chan int      


func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, comecoJogo chan bool) {
	mao := cartasIniciais         // estado local - cartas na mão do jogador
	nroDeCartas := M              // quantas cartas ele tem
	cartaRecebida := carta(" ")   // carta recebida é vazia
	jaBateu := false              // indica se o jogador já bateu

	fmt.Printf("jogador %d tem essas cartas na mao: %v\n", id, mao)

	// sinal de início do jogo
	<-comecoJogo  

	for {
		
		// verifica se alguém já bateu antes
		if len(canalBater) != 0 && !jaBateu {
			fmt.Printf("jogador %d bateu!\n", id)
			// se for o último a bater, perde o jogo
			if len(canalBater) == NJ-1 {
				fmt.Printf("o jogador %d foi o ultimo a bater e perdeu o jogo.\n", id)
			}
			// bate
			canalBater <- id
			return
		}

		// se o jogador tiver cartas suficientes para jogar
		if nroDeCartas == 5 {
			fmt.Println(id, "joga")
			aux := rand.Intn(nroDeCartas)   // escolhe uma carta aleatória da mão
			cartaParaPassar := mao[aux]     // carta que será passada

			// se tiver 4 cartas iguais, tenta passar uma diferente
			if temQuatroIguais(mao) {
				cartaParaPassar = diferente(mao)
			}

			fmt.Printf("jogador %d passou a carta %s adiante\n", id, cartaParaPassar)
			out <- cartaParaPassar  // envia carta para o próximo jogador

			// remove a carta passada da mão
			mao = append(mao[:aux], mao[aux+1:]...)
			fmt.Printf("mao do jogador %d: %v\n", id, mao)
			nroDeCartas-- 

			// 4 cartas iguais (bater)
			if temQuatroIguais(mao) && nroDeCartas == 4 {
				fmt.Printf("jogador %d bate!\n", id)
				jaBateu = true
				canalBater <- id
				novaRodada <- false
			}

			// se acabou as cartas, sinaliza nova rodada
			if nroDeCartas == 0 {
				novaRodada <- true
			}
		} else { 
			select {
			case cartaRecebida = <-in:  // recebe carta do canal
				fmt.Printf("jogador %d recebeu a carta %s\n", id, cartaRecebida)
				mao = append(mao, cartaRecebida)  // adiciona carta à mão
				nroDeCartas++
				fmt.Printf("nova mao do jogador %d: %v\n", id, mao)
			default:
				// espera passivamente se não houver carta disponível
			}
		}
	}
}

// função que verifica se há 4 cartas iguais na mão
func temQuatroIguais(mao []carta) bool {
	contagem := make(map[carta]int)
	for _, c := range mao {
		contagem[c]++
		if contagem[c] >= 4 {
			return true
		}
	}
	return false
}

// função que retorna uma carta diferente das repetidas
func diferente(mao []carta) carta {
	contagem := make(map[carta]int)
	for _, c := range mao {
		contagem[c]++
		if contagem[c] == 1 {
			return c
		}
	}
	return carta(" ")
}

func main() {
	// cria canais de passagem de cartas
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	canalBater = make(chan int, NJ)   // canal para registrar quem bateu
	novaRodada = make(chan bool)      // canal para sinalizar nova rodada
	comecoJogo := make(chan bool)     // canal para sinalizar início do jogo

	// cria jogadores
	for i := 0; i < NJ; i++ {
		// baralho = cria um baralho com NJ*M cartas
		// cartasEscolhidas = escolhe aleatoriamente (e tira) M cartas do baralho para o jogador i
		cartasEscolhidas := make([]carta, M)
		for j := 0; j < M; j++ {
			randomIndex := rand.Intn(4)
			cartasEscolhidas[j] = carta(rune('a' + randomIndex))
		}
		// cria goroutine do jogador com conexão circular (anel)
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas, comecoJogo)
	}

	// inicia o jogo 
	for i := 0; i < NJ; i++ {
		comecoJogo <- true
	}

	// escolhe um jogador j e escreve uma carta em seu canal de entrada
	cartaInicial := carta("joker")
	ch[0] <- cartaInicial

	// espera um pouco para simular jogadas antes do programa terminar
	time.Sleep(1 * time.Second)
}
