package dummy

import "fmt"

type Dummy struct {
	name string
	age  int32
}

// THIS is a comment
func (d *Dummy) GetName() string {
	return d.name
}

func (d *Dummy) GetAge() int32 {
	return d.age
}

func (d *Dummy) PrintAge() {
	fmt.Println(d.age)
}
