package main

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	original := x
	reversed := 0

	for x > 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}

	return original == reversed
}

// Cach 2
// func isPalindrome(x int) bool {
// 	// số âm hoặc kết thúc bằng 0 (trừ số 0) → false
// 	if x < 0 || (x%10 == 0 && x != 0) {
// 		return false
// 	}

// 	reversed := 0

// 	for x > reversed {
// 		reversed = reversed*10 + x%10
// 		x /= 10
// 	}

// 	// số chẵn: x == reversed
// 	// số lẻ: bỏ chữ số giữa → x == reversed/10
// 	return x == reversed || x == reversed/10
// }
