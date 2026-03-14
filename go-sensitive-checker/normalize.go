package main

import (
	"strings"
	"unicode"
)

// NormalizeResult stores normalized text and mapping back to original rune index.
type NormalizeResult struct {
	OriginalText   string   `json:"original_text"`
	NormalizedText string   `json:"normalized_text"`
	Steps          []string `json:"steps"`
	IndexMap       []int    `json:"-"`
}

type Normalizer struct {
	separatorRunes map[rune]struct{}
	replaceMap     map[rune]rune
}

type phraseMapping struct {
	from []rune
	to   []rune
}

type mappingProfile struct {
	enabled bool
	charMap map[rune]rune
	phrases []phraseMapping
}

func NewNormalizer() *Normalizer {
	seps := map[rune]struct{}{}
	for _, r := range []rune("-_*|/\\.,，。!?！？:：;；'\"`~·()（）[]{}<>《》+=%^&#@$") {
		seps[r] = struct{}{}
	}
	return &Normalizer{
		separatorRunes: seps,
		replaceMap: map[rune]rune{
			'@': 'a',
			'4': 'a',
			'0': 'o',
			'1': 'i',
			'!': 'i',
			'3': 'e',
			'5': 's',
			'$': 's',
			'7': 't',
			'+': 't',
			'8': 'b',
		},
	}
}

// NormalizeText runs baseline normalization for detection.
func (n *Normalizer) NormalizeText(text string) NormalizeResult {
	return n.normalize(text, false, DefaultDetectOptions())
}

// NormalizeTextAggressive runs stronger normalization used by fuzzy/pinyin pipelines.
func (n *Normalizer) NormalizeTextAggressive(text string) NormalizeResult {
	return n.normalize(text, true, DefaultDetectOptions())
}

func (n *Normalizer) NormalizeTextWithOptions(text string, options DetectOptions) NormalizeResult {
	return n.normalize(text, false, options)
}

func (n *Normalizer) NormalizeTextAggressiveWithOptions(text string, options DetectOptions) NormalizeResult {
	return n.normalize(text, true, options)
}

func (n *Normalizer) normalize(text string, aggressive bool, options DetectOptions) NormalizeResult {
	baseRunes := []rune(text)
	normRunes := make([]rune, 0, len(baseRunes))
	normIdx := make([]int, 0, len(baseRunes))
	for i, r := range baseRunes {
		normRunes = append(normRunes, unicode.ToLower(fullToHalf(r)))
		normIdx = append(normIdx, i)
	}

	profile := n.buildMappingProfile(options)
	if profile.enabled && len(profile.phrases) > 0 {
		normRunes, normIdx = applyPhraseMappings(normRunes, normIdx, profile.phrases)
	}

	out := make([]rune, 0, len(normRunes))
	idxMap := make([]int, 0, len(normRunes))

	for i, r := range normRunes {
		originIdx := normIdx[i]
		if profile.enabled {
			if mapped, ok := profile.charMap[r]; ok {
				r = mapped
			}
		}
		if unicode.IsSpace(r) {
			continue
		}
		if _, ok := n.separatorRunes[r]; ok {
			continue
		}
		if aggressive {
			// Drop non letter/number in aggressive mode to improve obfuscation tolerance.
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
				continue
			}
		}
		out = append(out, r)
		idxMap = append(idxMap, originIdx)
	}

	steps := []string{"trim_space", "remove_common_separators", "lowercase", "fullwidth_to_halfwidth"}
	if profile.enabled {
		steps = append(steps, "term_mapping")
	}
	if aggressive {
		steps = append(steps, "aggressive_noise_removal", "replace_confusable_chars")
	}

	return NormalizeResult{
		OriginalText:   text,
		NormalizedText: strings.TrimSpace(string(out)),
		Steps:          steps,
		IndexMap:       idxMap,
	}
}

func (n *Normalizer) buildMappingProfile(options DetectOptions) mappingProfile {
	profile := mappingProfile{
		enabled: options.EnableTermMapping,
		charMap: map[rune]rune{},
		phrases: nil,
	}
	if !profile.enabled {
		return profile
	}

	if options.MappingMode != MappingModeOverride {
		for k, v := range n.replaceMap {
			profile.charMap[k] = v
		}
	}

	phrases := make([]phraseMapping, 0, len(options.CustomMappings))
	for _, mapping := range options.CustomMappings {
		from := normalizeMappingRunes(mapping.From)
		to := normalizeMappingRunes(mapping.To)
		if len(from) == 0 || len(to) == 0 {
			continue
		}
		if len(from) == 1 && len(to) == 1 {
			profile.charMap[from[0]] = to[0]
			continue
		}
		phrases = append(phrases, phraseMapping{from: from, to: to})
	}

	if len(phrases) > 1 {
		// Longest match first to avoid short phrase抢占.
		for i := 0; i < len(phrases)-1; i++ {
			for j := i + 1; j < len(phrases); j++ {
				if len(phrases[j].from) > len(phrases[i].from) {
					phrases[i], phrases[j] = phrases[j], phrases[i]
				}
			}
		}
	}
	profile.phrases = phrases
	return profile
}

func normalizeMappingRunes(s string) []rune {
	rs := []rune(strings.TrimSpace(s))
	if len(rs) == 0 {
		return nil
	}
	out := make([]rune, 0, len(rs))
	for _, r := range rs {
		out = append(out, unicode.ToLower(fullToHalf(r)))
	}
	return out
}

func applyPhraseMappings(runes []rune, idxMap []int, mappings []phraseMapping) ([]rune, []int) {
	if len(runes) == 0 || len(mappings) == 0 {
		return runes, idxMap
	}

	out := make([]rune, 0, len(runes))
	outIdx := make([]int, 0, len(idxMap))

	for i := 0; i < len(runes); {
		applied := false
		for _, mp := range mappings {
			if len(mp.from) == 0 || i+len(mp.from) > len(runes) {
				continue
			}
			if !hasRunePrefix(runes[i:i+len(mp.from)], mp.from) {
				continue
			}

			lastSourceOffset := len(mp.from) - 1
			for outPos, rr := range mp.to {
				out = append(out, rr)
				sourceOffset := outPos
				if sourceOffset > lastSourceOffset {
					sourceOffset = lastSourceOffset
				}
				outIdx = append(outIdx, idxMap[i+sourceOffset])
			}
			i += len(mp.from)
			applied = true
			break
		}
		if applied {
			continue
		}
		out = append(out, runes[i])
		outIdx = append(outIdx, idxMap[i])
		i++
	}

	return out, outIdx
}

func hasRunePrefix(src, prefix []rune) bool {
	if len(src) != len(prefix) {
		return false
	}
	for i := range prefix {
		if src[i] != prefix[i] {
			return false
		}
	}
	return true
}

func fullToHalf(r rune) rune {
	// Full-width space.
	if r == 12288 {
		return 32
	}
	// Full-width ASCII variants.
	if r >= 65281 && r <= 65374 {
		return r - 65248
	}
	return r
}

func mapToOriginalIndex(idxMap []int, idx int) int {
	if idx < 0 || idx >= len(idxMap) {
		return idx
	}
	return idxMap[idx]
}
