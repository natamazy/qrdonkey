package cmd

import "fmt"

func PrintDonkey() {
	fmt.Printf("\033[92m" + QRdonkey)
	fmt.Printf("\033[0m")
}

const QRdonkey = `
                          /\          /\
                         ( \\        // )
                          \ \\      // /
                           \_\\||||//_/
                            \/ _  _ \
                           \/|(O)(O)|
                          \/ |      |
      ___________________\/  \      /
     //                //     |____|
    //                ||     /      \
   //|                \|     \ 0  0 /
  // \       )         V    / \____/   /--------\
 //   \     /        (     /          | QR saved |
""     \   /_________|  |_/            \--------/
       /  /\   /     |  ||
      /  / /  /      \  ||
      | |  | |        | ||
      | |  | |        | ||
      |_|  |_|        |_||
       \_\  \_\        \_\\
`
