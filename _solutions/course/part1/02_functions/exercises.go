// Reference solution for lesson 02.
package lesson02

func RectangleArea(width float64, height float64) float64 {
	return width * height
}

func Perimeter(w, h float64) float64 {
	return 2 * (w + h)
}

func Swap(a, b string) (string, string) {
	return b, a
}

func Divmod(a, b int) (q, r int) {
	q = a / b
	r = a % b
	return
}

func SumAndProduct(a, b, c int) (int, int) {
	return a + b + c, a * b * c
}

func ApplyTwice(f func(int) int, n int) int {
	return f(f(n))
}

func MakeMultiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor
	}
}
