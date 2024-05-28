package main

import "sort"


type FuzzyObject interface {
	GetStringFuzzy() *string
}

type FuzzyMatch struct {
	Match FuzzyObject
	Score int
}

type FuzzySearcher interface {
	Search(query string, options []FuzzyObject) (r []FuzzyMatch, err error)
}

type ByScore []FuzzyMatch

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool{
	return a[i].Score > a[j].Score
}


type Token struct {
	Val string
}

func newToken(s string) Token {
	return Token{s}
}

func MakeTokens(s string) (arr []Token){
	token := ""
	for _, c := range s {
		if c == ' ' {
			arr = append(arr, newToken(token))
			token = ""
		} else {
			token += string(c)
		}
	}

	if token != "" {
		arr = append(arr, newToken(token))
	}
	return arr
}

func TokenDistance(a, b Token) (d int) {
	d = 100
	for i := 0; i < len(a.Val) && i < len(b.Val); i++ {
		if a.Val[i] != b.Val[i] {
			d--
		}
	}

	dlen := len(a.Val) - len(b.Val)
	if dlen < 0 {
		dlen = -dlen
	}
	d -= dlen
	return
}

type TokenArray []Token

type TokenFuzzy struct{}
func (f *TokenFuzzy) Search(query string, options []FuzzyObject) (r []FuzzyMatch, err error) {
	q := MakeTokens(query);
	for _, option := range options {
		s := MakeTokens(*option.GetStringFuzzy())

		score := 0
		for _, qToken := range q {
			closest := 0
			for _, sToken := range s {
				temp := TokenDistance(qToken, sToken)
				if temp > closest {
					closest = temp
				}
			}

			score += closest
		}


		r = append(r, FuzzyMatch{option, score})
	}

	sort.Sort(ByScore(r))

	return
}
