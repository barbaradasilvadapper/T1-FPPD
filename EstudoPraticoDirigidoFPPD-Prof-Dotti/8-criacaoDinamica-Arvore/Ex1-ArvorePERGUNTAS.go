// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   [OS ACIMA ESTAO RESOLVIDOS]

//   2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func busca(r* Nodo, v int) bool {}...}
//   2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a busca), ou informa
//        que nao encontrou ao final da busca

//   3.a) a operação que escreve todos pares em um canal de saidaPares e
//        todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado

//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
)

type Nodo struct {
	v int
	e *Nodo
	d *Nodo
}

func caminhaERD(r *Nodo) {
	if r != nil {
		caminhaERD(r.e)
		fmt.Print(r.v, ", ")
		caminhaERD(r.d)
	}
}

// -------- SOMA ----------
// soma sequencial recursiva
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

// funcao "wraper" retorna valor
// internamente dispara recursao com somaConcCh
// usando canais
func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}
func somaConcCh(r *Nodo, s chan int) {
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.e, s1)
		go somaConcCh(r.d, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

// -------- BUSCA ----------
// busca sequencial recursiva
// func busca(r *Nodo, val int) bool { }
func busca(r *Nodo, val int) bool {
	if r != nil {
		if r.v == val {
			return true
		}
		return busca(r.e, val) || busca(r.d, val)
	}
	return false
}



// busca concorrente recursiva
// func buscaConc(r *Nodo, val int, ret chan bool) { }
func buscaConc(r *Nodo, val int) bool {
	ret := make(chan bool)
    go buscaConcCh(r, val, ret)
    return <-ret
}
func buscaConcCh(r *Nodo, val int, ret chan bool) {
	if r == nil {
        ret <- false
        return
    }
	if r.v == val {
		ret <- true
		return
	}
	if r != nil {
		if r.v == val {
			ret <- true
			return
		}

		// cria dois canais, um pra cada filho
		s1 := make(chan bool)
		s2 := make(chan bool)

		//chama a busca para esquerda e direita
		go buscaConcCh(r.e, val, s1)
		go buscaConcCh(r.d, val, s2)

		// combina os resultados
		esq := <-s1
		dir := <-s2

		ret <- (esq || dir)

	} else {
		ret <- false
	}
}



// -------- SAIDAS PAR E IMPAR --------
// sequencial 
// func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) { }
func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r != nil {
		retornaParImparSaida(r, saidaP, saidaI)
	}
	// sinaliza que terminou a operação
	fin <- struct{}{}
}

func retornaParImparSaida(r *Nodo, saidaP chan int, saidaI chan int) {
	if r != nil {
		retornaParImparSaida(r.e, saidaP, saidaI)
		if r.v%2 == 0 {
			saidaP <- r.v
		} else {
			saidaI <- r.v
		}
		retornaParImparSaida(r.d, saidaP, saidaI)
	}
	if r == nil {
		return
	}
}


// concorrente recursiva
// func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) { }
func retornaParImparConc(r *Nodo, saidaP, saidaI chan int, fin chan struct{}) {
    retornaParImparConcCh(r, saidaP, saidaI, fin) // chama diretamente, sem go
}

func retornaParImparConcCh(r *Nodo, saidaP, saidaI chan int, fin chan struct{}) {
    if r == nil {
        fin <- struct{}{}
        return
    }

    finE := make(chan struct{})
    finD := make(chan struct{})

    // chama os filhos concorrentemente
    go retornaParImparConcCh(r.e, saidaP, saidaI, finE)
    go retornaParImparConcCh(r.d, saidaP, saidaI, finD)

    // processa o nodo atual
    if r.v%2 == 0 {
        saidaP <- r.v
    } else {
        saidaI <- r.v
    }

    // espera os filhos terminarem
    <-finE
    <-finD

    // sinaliza que terminou
    fin <- struct{}{}
}







// ---------   agora vamos criar a arvore e usar as funcoes acima

func main() {
	root := &Nodo{v: 10,
		e: &Nodo{v: 5,
			e: &Nodo{v: 3,
				e: &Nodo{v: 1, e: nil, d: nil},
				d: &Nodo{v: 4, e: nil, d: nil}},
			d: &Nodo{v: 7,
				e: &Nodo{v: 6, e: nil, d: nil},
				d: &Nodo{v: 8, e: nil, d: nil}}},
		d: &Nodo{v: 15,
			e: &Nodo{v: 13,
				e: &Nodo{v: 12, e: nil, d: nil},
				d: &Nodo{v: 14, e: nil, d: nil}},
			d: &Nodo{v: 18,
				e: &Nodo{v: 17, e: nil, d: nil},
				d: &Nodo{v: 19, e: nil, d: nil}}}}

	//sequencial
	saidaP := make(chan int)
	saidaI := make(chan int)
	fin := make(chan struct{})

	fmt.Println()
	fmt.Println()
	fmt.Println("Valores na árvore: ")
	go retornaParImpar(root, saidaP, saidaI, fin)
	fim := false
	for count := 0; count < 20 && !fim; {
		select {
		case par := <-saidaP:
			fmt.Println("Par:", par)
			count++
		case impar := <-saidaI:
			fmt.Println("Impar:", impar)
			count++
		case <-fin:
			fim = true
		}
	}

	// concorrente
	saidaPConc := make(chan int)
	saidaIConc := make(chan int)
	finConc := make(chan struct{})

	fmt.Println()
	fmt.Println("Valores na árvore (concorrente): ")
	go retornaParImparConc(root, saidaPConc, saidaIConc, finConc)
	fim2 := false
	for !fim2 {
		select {
		case par := <-saidaPConc:
			fmt.Println("Par:", par)
		case impar := <-saidaIConc:
			fmt.Println("Impar:", impar)
		case <-finConc:
			fim2 = true
		}
	}


	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()
	fmt.Println("Busca 17: ", busca(root, 17))
	fmt.Println("Busca 99: ", busca(root, 99))

	fmt.Println()
	fmt.Println("BuscaConc 17: ", buscaConc(root, 17))
	fmt.Println("BuscaConc 99: ", buscaConc(root, 99))
}
