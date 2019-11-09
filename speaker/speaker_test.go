package speaker

import (
	"fmt"
)

//nolint:lll
func Example() {
	fmt.Println(BehindTheWall.Apply("Hello") + Chainsaw1 + " world!")
	// Output:
	// <speaker effect="behind_the_wall">Hello<speaker effect="-"><speaker audio="alice-sounds-things-chainsaw-1.opus"> world!
}
