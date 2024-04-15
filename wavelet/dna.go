package wavelet

func ToDNA(s []byte) {
	// code -> byte
	// to := []byte{'$', '#', 'A', 'C', 'G', 'T', 'N', ' '}
	to := []byte{'A', 'C', 'G', 'T'}
	// byte -> code
	var from [256]byte
	for i := 0; i < 4; i++ {
		from[to[i]] = byte(i)
	}
}
