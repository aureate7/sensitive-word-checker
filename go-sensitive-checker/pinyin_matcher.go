package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// PinyinMatcher builds pinyin alias index for Chinese sensitive words.
// It supports built-in aliases and an optional external alias file.
type PinyinMatcher struct {
	normalizer     *Normalizer
	aliases        map[string][]string // word -> manual aliases
	enableAuto     bool
	enableInitials bool
	args           pinyin.Args
}

func NewPinyinMatcher(normalizer *Normalizer, aliasFile string, enableAuto bool, enableInitials bool) *PinyinMatcher {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	args.Heteronym = false
	args.Fallback = func(r rune, _ pinyin.Args) []string {
		r = fullToHalf(r)
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return []string{strings.ToLower(string(r))}
		}
		return nil
	}
	m := &PinyinMatcher{
		normalizer:     normalizer,
		aliases:        defaultPinyinAliases(),
		enableAuto:     enableAuto,
		enableInitials: enableInitials,
		args:           args,
	}
	m.loadAliasFile(aliasFile)
	return m
}

func (m *PinyinMatcher) BuildAliasIndex(words map[string]struct{}) *AliasIndex {
	aliasToKW := make(map[string][]string)
	for w := range words {
		aliasSet := map[string]struct{}{}
		for _, rawAlias := range m.aliases[w] {
			alias := m.normalizer.NormalizeTextAggressive(rawAlias).NormalizedText
			if alias == "" {
				continue
			}
			aliasSet[alias] = struct{}{}
		}
		if m.enableAuto {
			for _, rawAlias := range m.autoPinyinAliases(w) {
				alias := m.normalizer.NormalizeTextAggressive(rawAlias).NormalizedText
				if alias == "" {
					continue
				}
				aliasSet[alias] = struct{}{}
			}
		}
		for alias := range aliasSet {
			aliasToKW[alias] = appendUniq(aliasToKW[alias], w)
		}
	}
	aliases := make([]string, 0, len(aliasToKW))
	for a := range aliasToKW {
		aliases = append(aliases, a)
	}
	sort.Strings(aliases)
	ac := NewAC()
	ac.Build(aliases)
	return &AliasIndex{Automaton: ac, AliasToKW: aliasToKW}
}

func (m *PinyinMatcher) loadAliasFile(path string) {
	if strings.TrimSpace(path) == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		word := strings.TrimSpace(parts[0])
		if word == "" {
			continue
		}
		right := strings.TrimSpace(parts[1])
		if right == "" {
			continue
		}
		for _, a := range strings.Split(right, "|") {
			a = strings.TrimSpace(a)
			if a == "" {
				continue
			}
			m.aliases[word] = appendUniq(m.aliases[word], a)
		}
	}
}

func defaultPinyinAliases() map[string][]string {
	// Built-in special aliases (short forms / common abbreviations).
	return map[string][]string{
		"傻逼":  {"shabi", "sha bi", "sha_bi", "sha-bi", "sb"},
		"傻比":  {"shabi", "sha bi", "sha_bi", "sha-bi", "sb"},
		"笨蛋":  {"bendan", "ben dan"},
		"操你":  {"caoni", "cao ni"},
		"操你妈": {"caonima", "cao ni ma", "cnm"},
		"草你妈": {"caonima", "cao ni ma", "cnm"},
		"他妈的": {"tamade", "ta ma de", "tmd"},
		"妈的":  {"made", "ma de", "md"},
		"鸡巴":  {"jiba", "ji ba"},
		"鸡吧":  {"jiba", "ji ba"},
		"屌":   {"diao"},
		"滚":   {"gun"},
	}
}

func (m *PinyinMatcher) autoPinyinAliases(word string) []string {
	if strings.TrimSpace(word) == "" {
		return nil
	}
	hasHan := false
	for _, r := range word {
		if unicode.Is(unicode.Han, r) {
			hasHan = true
			break
		}
	}
	if !hasHan {
		return nil
	}

	lazy := pinyin.LazyPinyin(word, m.args)
	if len(lazy) == 0 {
		return nil
	}

	parts := make([]string, 0, len(lazy))
	for _, s := range lazy {
		s = strings.TrimSpace(strings.ToLower(s))
		if s == "" {
			continue
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return nil
	}

	set := map[string]struct{}{}
	joined := strings.Join(parts, "")
	if joined != "" {
		set[joined] = struct{}{}
	}
	if len(parts) > 1 {
		set[strings.Join(parts, " ")] = struct{}{}
	}
	if m.enableInitials && len(parts) > 1 {
		var b strings.Builder
		for _, p := range parts {
			r := []rune(p)
			if len(r) == 0 {
				continue
			}
			b.WriteRune(r[0])
		}
		if b.Len() >= 2 {
			set[b.String()] = struct{}{}
		}
	}

	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
