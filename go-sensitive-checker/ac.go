package main

import (
	"container/list"
	"sort"
	"unicode/utf8"
)

type acNode struct {
	next map[rune]*acNode
	fail *acNode
	out  []string // 命中词（可换索引以省内存）
}

type Match struct {
	Start int    `json:"start"` // [Start, End) 的 rune 下标
	End   int    `json:"end"`
	Word  string `json:"word"`
}

type ACAutomaton struct {
	root     *acNode
	numWords int
}

func NewAC() *ACAutomaton {
	return &ACAutomaton{
		root: &acNode{next: make(map[rune]*acNode)},
	}
}

func (ac *ACAutomaton) Build(words []string) {
	// Trie
	for _, w := range words {
		if w == "" {
			continue
		}
		cur := ac.root
		for _, ch := range w { // rune 级
			if cur.next[ch] == nil {
				cur.next[ch] = &acNode{next: make(map[rune]*acNode)}
			}
			cur = cur.next[ch]
		}
		// 去重追加
		if len(cur.out) == 0 || cur.out[len(cur.out)-1] != w {
			cur.out = append(cur.out, w)
		}
		ac.numWords++
	}

	// fail 构建（BFS）
	q := list.New()
	for _, nxt := range ac.root.next {
		nxt.fail = ac.root
		q.PushBack(nxt)
	}
	for q.Len() > 0 {
		u := q.Remove(q.Front()).(*acNode)
		for ch, v := range u.next {
			f := u.fail
			for f != nil && f.next[ch] == nil {
				f = f.fail
			}
			if f == nil {
				v.fail = ac.root
			} else {
				v.fail = f.next[ch]
				if len(v.fail.out) > 0 {
					v.out = append(v.out, v.fail.out...)
				}
			}
			q.PushBack(v)
		}
	}
}

// Search 返回所有命中（包含重叠）
func (ac *ACAutomaton) Search(s string) []Match {
	if ac.root == nil {
		return nil
	}
	runes := []rune(s)
	var res []Match
	cur := ac.root
	for i, ch := range runes {
		for cur != ac.root && cur.next[ch] == nil {
			cur = cur.fail
		}
		if nxt := cur.next[ch]; nxt != nil {
			cur = nxt
		}
		if len(cur.out) > 0 {
			for _, w := range cur.out {
				wlen := utf8.RuneCountInString(w)
				res = append(res, Match{
					Start: i - wlen + 1,
					End:   i + 1,
					Word:  w,
				})
			}
		}
	}
	return res
}

// Mask: 若以后需要掩码可用（当前项目未开放 HTTP 掩码接口）
func (ac *ACAutomaton) Mask(s string, mask rune) (string, []Match) {
	matches := ac.Search(s)
	if len(matches) == 0 {
		return s, nil
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Start == matches[j].Start {
			return matches[i].End > matches[j].End
		}
		return matches[i].Start < matches[j].Start
	})
	rs := []rune(s)
	cover := make([]bool, len(rs))
	for _, m := range matches {
		if m.Start < 0 {
			continue
		}
		for i := m.Start; i < m.End && i < len(cover); i++ {
			cover[i] = true
		}
	}
	for i, c := range cover {
		if c {
			rs[i] = mask
		}
	}
	return string(rs), matches
}
