package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Token struct {
	Cnt int32
	Val string
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s %d", t.Val, t.Cnt)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Hash(content string, minTokenLen int, quantRate float64) string {
	if len(content) < 1 {
		return GetMD5Hash(content) // fallback to md5
	}

	var maxFreq int32 = 0

	tokens := make(map[string]int)
	var orderedTokens []*Token
	var profiles []*Token

	indexCounter := 0
	curToken := strings.Builder{}
	for _, c := range content {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			curToken.WriteRune(unicode.ToLower(c))
		} else {
			if curToken.Len() > 0 {
				if curToken.Len() > minTokenLen {
					tokenVal := curToken.String()
					index, exists := tokens[tokenVal]
					if !exists {
						index = indexCounter
						indexCounter++

						tok := &Token{
							Cnt: 0,
							Val: tokenVal,
						}
						tokens[tokenVal] = index
						orderedTokens = append(orderedTokens, tok)
					}
					orderedTokens[index].Cnt++
					if orderedTokens[index].Cnt > maxFreq {
						maxFreq = orderedTokens[index].Cnt
					}
				}
				curToken.Reset()
			}
		}
	}

	// Yes! DRY
	if curToken.Len() > 0 {
		if curToken.Len() > minTokenLen {
			tokenVal := curToken.String()
			index, exists := tokens[tokenVal]
			if !exists {
				index = indexCounter
				indexCounter++

				tok := &Token{
					Cnt: 0,
					Val: tokenVal,
				}
				tokens[tokenVal] = index
				orderedTokens = append(orderedTokens, tok)
			}
			orderedTokens[index].Cnt++
			if orderedTokens[index].Cnt > maxFreq {
				maxFreq = orderedTokens[index].Cnt
			}
		}
		curToken.Reset()
	}

	quant := int32(math.Round(float64(maxFreq) * quantRate))
	if quant < 2 {
		if maxFreq > 1 {
			quant = 2
		} else {
			quant = 1
		}
	}

	for _, t := range orderedTokens {
		t.Cnt = (t.Cnt / quant) * quant
		if t.Cnt < quant {
			continue
		}
		profiles = append(profiles, t)
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Cnt < profiles[j].Cnt
	})

	text := strings.Builder{}
	for _, v := range profiles {
		if text.Len() > 0 {
			text.WriteString("\n")
		}
		text.WriteString(v.ToString())
	}

	return GetMD5Hash(text.String())
}

func main() {
	minTokenLen := flag.Int("min_token_len", 2, "minimum token length")
	quantRate := flag.Float64("quant_rate", 0.01, "quant rate")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			break
		}
		fmt.Printf("%s\n", Hash(line, *minTokenLen, *quantRate))
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
