package gate

import "fmt"

func (g *Gate) open() {
	fmt.Printf("Gate %s: ⇈\n\r", g.Name)
}

func (g *Gate) close() {
	fmt.Printf("Gate %s: ⇊\n\n", g.Name)
}

func (g *Gate) stop() {
	fmt.Printf("Gate %s: ⇎\n\n", g.Name)
}
