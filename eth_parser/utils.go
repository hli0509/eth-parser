package ethparser

import "fmt"

func HexToDec(hex string) int {
	var dec int
	fmt.Sscanf(hex, "0x%x", &dec)
	return dec
}

func DecToHex(dec int) string {
	return fmt.Sprintf("0x%x", dec)
}