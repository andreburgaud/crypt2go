package padding

import "fmt"

func ExamplePadder_Pad() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D}
	fmt.Printf("%X\n", p)
	padder := NewPkcs5Padding()
	p, err := padder.Pad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D
	// 0A0B0C0D04040404
}

func ExamplePadder_Pad_second() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0A, 0x0B, 0x0C, 0x0D}
	fmt.Printf("%X\n", p)
	padder := NewPkcs5Padding()
	p, err := padder.Pad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D0A0B0C0D
	// 0A0B0C0D0A0B0C0D0808080808080808
}

func ExamplePadder_Pad_third() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D}
	fmt.Printf("%X\n", p)
	padder := NewPkcs7Padding(16) // 16-byte block size
	p, err := padder.Pad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D
	// 0A0B0C0D0C0C0C0C0C0C0C0C0C0C0C0C
}

func ExamplePadder_Unpad() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x04, 0x04, 0x04, 0x04}
	fmt.Printf("%X\n", p)
	padder := NewPkcs5Padding()
	p, err := padder.Unpad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D04040404
	// 0A0B0C0D
}

func ExamplePadder_Unpad_second() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0A, 0x0B, 0x0C, 0x0D,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08}
	fmt.Printf("%X\n", p)
	padder := NewPkcs5Padding()
	p, err := padder.Unpad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D0A0B0C0D0808080808080808
	// 0A0B0C0D0A0B0C0D
}

func ExamplePadder_Unpad_third() {
	p := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0C, 0x0C, 0x0C, 0x0C,
		0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C}
	fmt.Printf("%X\n", p)
	padder := NewPkcs7Padding(16) // 16-byte block size
	p, err := padder.Unpad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%X\n", p)
	// Output:
	// 0A0B0C0D0C0C0C0C0C0C0C0C0C0C0C0C
	// 0A0B0C0D
}

func ExamplePadder_Unpad_empty() {
	p := []byte{0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08}
	fmt.Printf("%X\n", p)
	padder := NewPkcs7Padding(8) // 8-byte block size
	p, err := padder.Unpad(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v\n", p)
	// Output:
	// 0808080808080808
	// []
}
