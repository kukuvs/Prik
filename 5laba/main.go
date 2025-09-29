package main

import (
	"fmt"
	"math"
)

// Person
type Person struct {
	Name string
	Age  int
}

func (p Person) Info() {
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}

func (p *Person) Birthday() {
	p.Age++
}

// Shape interface
type Shape interface {
	Area() float64
}

// Circle
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Rectangle
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// PrintAreas
func PrintAreas(shapes []Shape) {
	for i, s := range shapes {
		fmt.Printf("Shape %d area: %.2f\n", i+1, s.Area())
	}
}

// Stringer interface
type Stringer interface {
	String() string
}

// Book
type Book struct {
	Title  string
	Author string
	Pages  int
}

func (b Book) String() string {
	return fmt.Sprintf("%s by %s, %d pages", b.Title, b.Author, b.Pages)
}

func main() {
	// Person
	p := Person{Name: "Ivan", Age: 30}
	p.Info()
	p.Birthday()
	p.Info()

	// Shapes
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 3}
	shapes := []Shape{c, r}
	PrintAreas(shapes)

	// Book
	book := Book{Title: "Go in Action", Author: "William Kennedy", Pages: 300}
	fmt.Println(book.String())
}
