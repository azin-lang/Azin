// internal/azin/syntax/green/token.go
package green

import (
	"hash/maphash"
	"strings"
	"sync"

	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type Token struct {
	Base
	text           string
	leadingTrivia  Node
	trailingTrivia Node
}

func (t *Token) IsNil() bool { return t == nil }

type tokenCacheKey struct {
	kind syntax.Kind
	text string
}

const numShards = 32

type cacheShard struct {
	mu    sync.RWMutex
	items map[tokenCacheKey]*Token
}

var (
	hashSeed = maphash.MakeSeed()
	shards   [numShards]cacheShard
)

func init() {
	for i := range numShards {
		shards[i].items = make(map[tokenCacheKey]*Token)
	}
}

func getShard(kind syntax.Kind, text string) *cacheShard {
	h := maphash.String(hashSeed, text)
	idx := (h ^ uint64(kind)) % numShards
	return &shards[idx]
}

// NewToken constructs an immutable green token, safely interning zero-trivia tokens.
func NewToken(kind syntax.Kind, text string, leading, trailing Node) *Token {
	if IsNil(leading) && IsNil(trailing) {
		key := tokenCacheKey{kind: kind, text: text}
		shard := getShard(kind, text)

		// Fast path: Read lock
		shard.mu.RLock()
		if cached, ok := shard.items[key]; ok {
			shard.mu.RUnlock()
			return cached
		}
		shard.mu.RUnlock()

		internedText := strings.Clone(text)
		internedKey := tokenCacheKey{kind: kind, text: internedText}

		token := &Token{
			Base: Base{
				kind:      kind,
				fullWidth: len(internedText),
				flags:     FlagNone,
			},
			text: internedText,
		}

		shard.mu.Lock()
		if cached, ok := shard.items[key]; ok {
			shard.mu.Unlock()
			return cached
		}
		shard.items[internedKey] = token
		shard.mu.Unlock()

		return token
	}

	wL, fL := getProps(leading)
	wT, fT := getProps(trailing)

	return &Token{
		Base: Base{
			kind:      kind,
			fullWidth: wL + wT + len(text),
			flags:     fL | fT,
		},
		text:           text,
		leadingTrivia:  leading,
		trailingTrivia: trailing,
	}
}

func (t *Token) Text() string         { return t.text }
func (t *Token) LeadingTrivia() Node  { return t.leadingTrivia }
func (t *Token) TrailingTrivia() Node { return t.trailingTrivia }

func (t *Token) SlotCount() int      { return 0 }
func (t *Token) Slot(index int) Node { return nil }
