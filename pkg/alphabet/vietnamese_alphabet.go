package alphabet

// vietnameseAlphabet represents vietnamese alphabet а-z
type vietnameseAlphabet struct {
	parent Alphabet
}

// NewVietnameseAlphabet returns new instance of VietnameseAlphabet
func NewVietnameseAlphabet() Alphabet {
	return &vietnameseAlphabet{
		NewSequentialAlphabet('a', 'y'),
	}
}

var (
	mSpencialRune = map[rune]rune{
		'à': 'a',
		'á': 'a',
		'ạ': 'a',
		'ả': 'a',
		'ã': 'a',
		'â': 'a',
		'ầ': 'a',
		'ấ': 'a',
		'ậ': 'a',
		'ẩ': 'a',
		'ẫ': 'a',
		'ă': 'a',
		'ằ': 'a',
		'ắ': 'a',
		'ặ': 'a',
		'ẳ': 'a',
		'ẵ': 'a',
		'è': 'e',
		'é': 'e',
		'ẹ': 'e',
		'ẻ': 'e',
		'ẽ': 'e',
		'ê': 'e',
		'ề': 'e',
		'ế': 'e',
		'ệ': 'e',
		'ể': 'e',
		'ễ': 'e',
		'ì': 'i',
		'í': 'i',
		'ị': 'i',
		'ỉ': 'i',
		'ĩ': 'i',
		'ò': 'o',
		'ó': 'o',
		'ọ': 'o',
		'ỏ': 'o',
		'õ': 'o',
		'ô': 'o',
		'ồ': 'o',
		'ố': 'o',
		'ộ': 'o',
		'ổ': 'o',
		'ỗ': 'o',
		'ơ': 'o',
		'ờ': 'o',
		'ớ': 'o',
		'ợ': 'o',
		'ở': 'o',
		'ỡ': 'o',
		'ù': 'u',
		'ú': 'u',
		'ụ': 'u',
		'ủ': 'u',
		'ũ': 'u',
		'ư': 'u',
		'ừ': 'u',
		'ứ': 'u',
		'ự': 'u',
		'ử': 'u',
		'ữ': 'u',
		'ỳ': 'y',
		'ý': 'y',
		'ỵ': 'y',
		'ỷ': 'y',
		'ỹ': 'y',
		'đ': 'd',
	}
)

// Note, that we map ё as e
func (a *vietnameseAlphabet) Has(char rune) bool {
	pureRune, has := mSpencialRune[char]
	if has {
		return a.parent.Has(pureRune)
	}
	return a.parent.Has(char)
}

func (a *vietnameseAlphabet) Size() int {
	return a.parent.Size()
}

func (a *vietnameseAlphabet) Chars() []rune {
	return a.parent.Chars()
}
