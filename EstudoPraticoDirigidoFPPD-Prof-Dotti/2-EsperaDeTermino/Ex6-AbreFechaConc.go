// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// ABRE E FECHA CONCORRENCIA
// Há várias formas de esperar o término de processos concorrentes.
// Aqui temos um programa que lança processos concorrentes e espera o término dos mesmos.
//
// EXERCICIOS:
//    1) modifique o número de processos e o numero de iteracoes de cada processo

//    2) avalie o resultado obtido do ponto de vista da velocidade relativa entre os processos e da justiça.
//    a velocidade relativa entre os processos é independente, ou seja, cada processo
//    pode avançar em sua execução de forma independente dos outros. A justiça
//    não é garantida, pois alguns processos podem ser agendados mais frequentemente
//    do que outros, resultando em uma distribuição desigual de tempo de CPU entre eles
// 	  (alguns processos aparecem mais vezes seguidas do que outros)

//    3) observe que os itens comunicados pelo canal fin são vazios.
//       isto significa que o importante neste caso é somente a sincronização.
//	  Quando uma goroutine escreve fin <- struct{}{}, significa “terminei”.
//    O valor em si não interessa, só interessa o fato de o envio ter acontecido.
//    Por isso é vazio (struct{} não ocupa espaço na memória).

//    4) 'fin' é um canal sincrono.
//       faria diferença se 'fin' fosse assíncrono, ou seja, se tivesse um buffer para armazenar itens?
//    Funcionalmente não muda nada, mas muda o tempo de bloqueio das goroutines.
//	  Com canal síncrono, a goroutine que termina fica bloqueada até que o main leia do canal.
//    Com canal assíncrono, as goroutines não bloqueiam ao enviar, pois o buffer armazena os sinais 
//    até o main ler.

package main

import (
	"fmt"
)

func algoConcorrente(id int, par int, fin chan struct{}) {
	for i := 0; i < par; i++ {
		fmt.Println(id, " fazendo algo ", i)
	}
	fin <- struct{}{} // sinaliza final
}

func main() {
	fin := make(chan struct{})

	// cria 5 rotinas concorrentes
	for i := 0; i < 10; i++ {
		go algoConcorrente(i, 3, fin) // passa canal fin para avisar o termino
	}

	// espera o termino das rotinas
	for i := 0; i < 10; i++ {
		<-fin // wait for 5 processes to write in ch
	}
	fmt.Println("fim")
}
