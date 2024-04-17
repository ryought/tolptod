package wavelet

// Encode DNA (in byte = 8bits) into 3bit representation
func ToCompactDNA(s []byte) []byte {
	ret := make([]byte, len(s))
	for i, c := range s {
		switch c {
		case 'A', 'a':
			ret[i] = 0b001
		case 'C', 'c':
			ret[i] = 0b011
		case 'G', 'g':
			ret[i] = 0b101
		case 'T', 't':
			ret[i] = 0b111
		case 0:
			ret[i] = 0b000
		case '$':
			ret[i] = 0b010
		case '#':
			ret[i] = 0b110
		case ' ':
			ret[i] = 0b100
		default:
			panic("Failed to encode")
		}
	}
	return ret
}

// Decode DNA 3bit representation into byte (ASCII code)
func FromCompactDNA(s []byte) []byte {
	ret := make([]byte, len(s))
	for i, c := range s {
		switch c {
		case 0b001:
			ret[i] = 'A'
		case 0b011:
			ret[i] = 'C'
		case 0b101:
			ret[i] = 'G'
		case 0b111:
			ret[i] = 'T'
		case 0b000:
			ret[i] = 0
		case 0b010:
			ret[i] = '$'
		case 0b110:
			ret[i] = '#'
		case 0b100:
			ret[i] = ' '
		default:
			panic("Failed to decode")
		}
	}
	return ret
}
