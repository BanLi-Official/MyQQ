package MyFmt

import(
	"bufio"
	"os"
)


func Scanf(a *string) {
    reader := bufio.NewReader(os.Stdin)
    data, _, _ := reader.ReadLine()
    *a = string(data)
}