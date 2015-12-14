// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translit

import "unicode"

/*
	@todo implmement https://github.com/magento/magento2/blob/master/lib%2Finternal%2FMagento%2FFramework%2FFilter%2FTranslit.php#L345
	needs maybe a re-architecture
*/

// Runes replaces characters in a string using a conversion table,
// e.g. ™ => tm; © => c; @ => at; € => euro
// Magento compatible
func Runes(str []rune) []rune {
	// @todo thread safe
	i := 0
	for i < len(str) {
		if to, ok := charMap[str[i]]; ok {
			if len(to) < 1 {
				i++
				continue
			}
			if len(to) == 1 {
				str[i] = to[0]
				i++
				continue
			}
			str = append(str[:i], append(to, str[i+1:]...)...)
			i = 0 // reset i an re-run the str because length changed
		} else {
			i++
		}
	}
	return str
}

// URL same as Runes() but removes after conversion any non 0-9A-Za-z
// characters and replaces them with a dash. Magento compatible.
// This function is responsible for the slug generation for product/category/cms URLs.
func URL(str []rune) []rune {
	str = Runes(str)
	for i, r := range str {
		str[i] = unicode.ToLower(r)
		switch {
		case '0' <= r && r <= '9':
			continue
		case 'a' <= r && r <= 'z':
			continue
		case 'A' <= r && r <= 'Z':
			continue
		case r == '-':
			continue
		default:
			str[i] = '-'
		}
	}
	i := 0
	// remove multiple dashes
	for i < len(str) {
		j := i + 1
		if str[i] == '-' && j < len(str) && str[j] == '-' { // look ahead
			str = append(str[:i], str[j:]...)
			i = 0 // reset i due to change in length of str
			continue
		}
		i++
	}
	// trim
	if len(str) > 0 {
		if str[0] == '-' {
			str = str[1:]
		}
		if str[len(str)-1] == '-' {
			str = str[:len(str)-1]
		}
	}
	return str
}

// URLUnicode takes care of unicode characters. @todo
func URLUnicode(str []rune) []rune {
	// @see https://github.com/raintreeinc/knowledgebase/blob/master/kb/slug.go
	// and var translitConvertTableExtended below
	return nil
}

// charMap courtesy: https://github.com/magento/magento2/blob/master/lib%2Finternal%2FMagento%2FFramework%2FFilter%2FTranslit.php
var charMap = map[rune][]rune{
	'€': []rune{'e', 'u', 'r', 'o'},
	'&': []rune{'a', 'n', 'd'},
	'@': []rune{'a', 't'},
	'©': []rune{'c'},
	'®': []rune{'r'},
	'À': []rune{'a'},
	'Á': []rune{'a'},
	'Â': []rune{'a'},
	'Ä': []rune{'a'},
	'Å': []rune{'a'},
	'Æ': []rune{'a', 'e'},
	'Ç': []rune{'c'},
	'È': []rune{'e'},
	'É': []rune{'e'},
	'Ë': []rune{'e'},
	'Ì': []rune{'i'},
	'Í': []rune{'i'},
	'Î': []rune{'i'},
	'Ï': []rune{'i'},
	'Ò': []rune{'o'},
	'Ó': []rune{'o'},
	'Ô': []rune{'o'},
	'Õ': []rune{'o'},
	'Ö': []rune{'o'},
	'Ø': []rune{'o'},
	'Ù': []rune{'u'},
	'Ú': []rune{'u'},
	'Û': []rune{'u'},
	'Ü': []rune{'u'},
	'Ý': []rune{'y'},
	'ß': []rune{'s', 's'},
	'à': []rune{'a'},
	'á': []rune{'a'},
	'â': []rune{'a'},
	'ä': []rune{'a'},
	'å': []rune{'a'},
	'æ': []rune{'a', 'e'},
	'ç': []rune{'c'},
	'è': []rune{'e'},
	'é': []rune{'e'},
	'ê': []rune{'e'},
	'ë': []rune{'e'},
	'ì': []rune{'i'},
	'í': []rune{'i'},
	'î': []rune{'i'},
	'ï': []rune{'i'},
	'ò': []rune{'o'},
	'ó': []rune{'o'},
	'ô': []rune{'o'},
	'õ': []rune{'o'},
	'ö': []rune{'o'},
	'ø': []rune{'o'},
	'ù': []rune{'u'},
	'ú': []rune{'u'},
	'û': []rune{'u'},
	'ü': []rune{'u'},
	'ý': []rune{'y'},
	'þ': []rune{'p'},
	'ÿ': []rune{'y'},
	'Ā': []rune{'a'},
	'ā': []rune{'a'},
	'Ă': []rune{'a'},
	'ă': []rune{'a'},
	'Ą': []rune{'a'},
	'ą': []rune{'a'},
	'Ć': []rune{'c'},
	'ć': []rune{'c'},
	'Ĉ': []rune{'c'},
	'ĉ': []rune{'c'},
	'Ċ': []rune{'c'},
	'ċ': []rune{'c'},
	'Č': []rune{'c'},
	'č': []rune{'c'},
	'Ď': []rune{'d'},
	'ď': []rune{'d'},
	'Đ': []rune{'d'},
	'đ': []rune{'d'},
	'Ē': []rune{'e'},
	'ē': []rune{'e'},
	'Ĕ': []rune{'e'},
	'ĕ': []rune{'e'},
	'Ė': []rune{'e'},
	'ė': []rune{'e'},
	'Ę': []rune{'e'},
	'ę': []rune{'e'},
	'Ě': []rune{'e'},
	'ě': []rune{'e'},
	'Ĝ': []rune{'g'},
	'ĝ': []rune{'g'},
	'Ğ': []rune{'g'},
	'ğ': []rune{'g'},
	'Ġ': []rune{'g'},
	'ġ': []rune{'g'},
	'Ģ': []rune{'g'},
	'ģ': []rune{'g'},
	'Ĥ': []rune{'h'},
	'ĥ': []rune{'h'},
	'Ħ': []rune{'h'},
	'ħ': []rune{'h'},
	'Ĩ': []rune{'i'},
	'ĩ': []rune{'i'},
	'Ī': []rune{'i'},
	'ī': []rune{'i'},
	'Ĭ': []rune{'i'},
	'ĭ': []rune{'i'},
	'Į': []rune{'i'},
	'į': []rune{'i'},
	'İ': []rune{'i'},
	'ı': []rune{'i'},
	'Ĳ': []rune{'i', 'j'},
	'ĳ': []rune{'i', 'j'},
	'Ĵ': []rune{'j'},
	'ĵ': []rune{'j'},
	'Ķ': []rune{'k'},
	'ķ': []rune{'k'},
	'ĸ': []rune{'k'},
	'Ĺ': []rune{'l'},
	'ĺ': []rune{'l'},
	'Ļ': []rune{'l'},
	'ļ': []rune{'l'},
	'Ľ': []rune{'l'},
	'ľ': []rune{'l'},
	'Ŀ': []rune{'l'},
	'ŀ': []rune{'l'},
	'Ł': []rune{'l'},
	'ł': []rune{'l'},
	'Ń': []rune{'n'},
	'ń': []rune{'n'},
	'Ņ': []rune{'n'},
	'ņ': []rune{'n'},
	'Ň': []rune{'n'},
	'ň': []rune{'n'},
	'ŉ': []rune{'n'},
	'Ŋ': []rune{'n'},
	'ŋ': []rune{'n'},
	'Ō': []rune{'o'},
	'ō': []rune{'o'},
	'Ŏ': []rune{'o'},
	'ŏ': []rune{'o'},
	'Ő': []rune{'o'},
	'ő': []rune{'o'},
	'Œ': []rune{'o', 'e'},
	'œ': []rune{'o', 'e'},
	'Ŕ': []rune{'r'},
	'ŕ': []rune{'r'},
	'Ŗ': []rune{'r'},
	'ŗ': []rune{'r'},
	'Ř': []rune{'r'},
	'ř': []rune{'r'},
	'Ś': []rune{'s'},
	'ś': []rune{'s'},
	'Ŝ': []rune{'s'},
	'ŝ': []rune{'s'},
	'Ş': []rune{'s'},
	'ş': []rune{'s'},
	'Š': []rune{'s'},
	'š': []rune{'s'},
	'Ţ': []rune{'t'},
	'ţ': []rune{'t'},
	'Ť': []rune{'t'},
	'ť': []rune{'t'},
	'Ŧ': []rune{'t'},
	'ŧ': []rune{'t'},
	'Ũ': []rune{'u'},
	'ũ': []rune{'u'},
	'Ū': []rune{'u'},
	'ū': []rune{'u'},
	'Ŭ': []rune{'u'},
	'ŭ': []rune{'u'},
	'Ů': []rune{'u'},
	'ů': []rune{'u'},
	'Ű': []rune{'u'},
	'ű': []rune{'u'},
	'Ų': []rune{'u'},
	'ų': []rune{'u'},
	'Ŵ': []rune{'w'},
	'ŵ': []rune{'w'},
	'Ŷ': []rune{'y'},
	'ŷ': []rune{'y'},
	'Ÿ': []rune{'y'},
	'Ź': []rune{'z'},
	'ź': []rune{'z'},
	'Ż': []rune{'z'},
	'ż': []rune{'z'},
	'Ž': []rune{'z'},
	'ž': []rune{'z'},
	'ſ': []rune{'z'},
	'Ə': []rune{'e'},
	'ƒ': []rune{'f'},
	'Ơ': []rune{'o'},
	'ơ': []rune{'o'},
	'Ư': []rune{'u'},
	'ư': []rune{'u'},
	'Ǎ': []rune{'a'},
	'ǎ': []rune{'a'},
	'Ǐ': []rune{'i'},
	'ǐ': []rune{'i'},
	'Ǒ': []rune{'o'},
	'ǒ': []rune{'o'},
	'Ǔ': []rune{'u'},
	'ǔ': []rune{'u'},
	'Ǖ': []rune{'u'},
	'ǖ': []rune{'u'},
	'Ǘ': []rune{'u'},
	'ǘ': []rune{'u'},
	'Ǚ': []rune{'u'},
	'ǚ': []rune{'u'},
	'Ǜ': []rune{'u'},
	'ǜ': []rune{'u'},
	'Ǻ': []rune{'a'},
	'ǻ': []rune{'a'},
	'Ǽ': []rune{'a', 'e'},
	'ǽ': []rune{'a', 'e'},
	'Ǿ': []rune{'o'},
	'ǿ': []rune{'o'},
	'ə': []rune{'e'},
	'Ё': []rune{'j', 'o'},
	'Є': []rune{'e'},
	'І': []rune{'i'},
	'Ї': []rune{'i'},
	'А': []rune{'a'},
	'Б': []rune{'b'},
	'В': []rune{'v'},
	'Г': []rune{'g'},
	'Д': []rune{'d'},
	'Е': []rune{'e'},
	'Ж': []rune{'z', 'h'},
	'З': []rune{'z'},
	'И': []rune{'i'},
	'Й': []rune{'j'},
	'К': []rune{'k'},
	'Л': []rune{'l'},
	'М': []rune{'m'},
	'Н': []rune{'n'},
	'О': []rune{'o'},
	'П': []rune{'p'},
	'Р': []rune{'r'},
	'С': []rune{'s'},
	'Т': []rune{'t'},
	'У': []rune{'u'},
	'Ф': []rune{'f'},
	'Х': []rune{'h'},
	'Ц': []rune{'c'},
	'Ч': []rune{'c', 'h'},
	'Ш': []rune{'s', 'h'},
	'Щ': []rune{'s', 'c', 'h'},
	'Ъ': []rune{'-'},
	'Ы': []rune{'y'},
	'Ь': []rune{'-'},
	'Э': []rune{'j', 'e'},
	'Ю': []rune{'j', 'u'},
	'Я': []rune{'j', 'a'},
	'а': []rune{'a'},
	'б': []rune{'b'},
	'в': []rune{'v'},
	'г': []rune{'g'},
	'д': []rune{'d'},
	'е': []rune{'e'},
	'ж': []rune{'z', 'h'},
	'з': []rune{'z'},
	'и': []rune{'i'},
	'й': []rune{'j'},
	'к': []rune{'k'},
	'л': []rune{'l'},
	'м': []rune{'m'},
	'н': []rune{'n'},
	'о': []rune{'o'},
	'п': []rune{'p'},
	'р': []rune{'r'},
	'с': []rune{'s'},
	'т': []rune{'t'},
	'у': []rune{'u'},
	'ф': []rune{'f'},
	'х': []rune{'h'},
	'ц': []rune{'c'},
	'ч': []rune{'c', 'h'},
	'ш': []rune{'s', 'h'},
	'щ': []rune{'s', 'c', 'h'},
	'ъ': []rune{'-'},
	'ы': []rune{'y'},
	'ь': []rune{'-'},
	'э': []rune{'j', 'e'},
	'ю': []rune{'j', 'u'},
	'я': []rune{'j', 'a'},
	'ё': []rune{'j', 'o'},
	'є': []rune{'e'},
	'і': []rune{'i'},
	'ї': []rune{'i'},
	'Ґ': []rune{'g'},
	'ґ': []rune{'g'},
	'א': []rune{'a'},
	'ב': []rune{'b'},
	'ג': []rune{'g'},
	'ד': []rune{'d'},
	'ה': []rune{'h'},
	'ו': []rune{'v'},
	'ז': []rune{'z'},
	'ח': []rune{'h'},
	'ט': []rune{'t'},
	'י': []rune{'i'},
	'ך': []rune{'k'},
	'כ': []rune{'k'},
	'ל': []rune{'l'},
	'ם': []rune{'m'},
	'מ': []rune{'m'},
	'ן': []rune{'n'},
	'נ': []rune{'n'},
	'ס': []rune{'s'},
	'ע': []rune{'e'},
	'ף': []rune{'p'},
	'פ': []rune{'p'},
	'ץ': []rune{'C'},
	'צ': []rune{'c'},
	'ק': []rune{'q'},
	'ר': []rune{'r'},
	'ש': []rune{'w'},
	'ת': []rune{'t'},
	'™': []rune{'t', 'm'},
}

//var translitConvertTableExtended = map[rune]string{
//	'\U00000021': "excl",
//	'\U00000022': "quot",
//	'\U00000023': "num",
//	'\U00000024': "dollar",
//	'\U00000025': "percnt",
//	'\U00000026': "amp",
//	'\U00000027': "apos",
//	'\U00000028': "lpar",
//	'\U00000029': "rpar",
//	'\U0000002A': "ast",
//	'\U0000002B': "plus",
//	'\U0000002C': "comma",
//	'\U0000002E': "period",
//	'\U0000002F': "sol",
//	'\U0000003A': "colon",
//	'\U0000003B': "semi",
//	'\U0000003C': "lt",
//	'\U0000003D': "equals",
//	'\U0000003E': "gt",
//	'\U0000003F': "quest",
//	'\U00000040': "commat",
//	'\U0000005B': "lsqb",
//	'\U0000005C': "bsol",
//	'\U0000005D': "rsqb",
//	'\U0000005E': "hat",
//	'\U0000005F': "lowbar",
//	'\U00000060': "grave",
//	'\U0000007B': "lcub",
//	'\U0000007C': "vert",
//	'\U0000007D': "rcub",
//	'\U000000A1': "iexcl",
//	'\U000000A2': "cent",
//	'\U000000A3': "pound",
//	'\U000000A4': "curren",
//	'\U000000A5': "yen",
//	'\U000000A6': "brvbar",
//	'\U000000A7': "sect",
//	'\U000000A8': "uml",
//	'\U000000A9': "copy",
//	'\U000000AB': "laquo",
//	'\U000000AC': "not",
//	'\U000000AE': "reg",
//	'\U000000AF': "macr",
//	'\U000000B0': "deg",
//	'\U000000B1': "pm",
//	'\U000000B4': "acute",
//	'\U000000B6': "para",
//	'\U000000B7': "middot",
//	'\U000000B8': "cedil",
//	'\U000000BB': "raquo",
//	'\U000000BF': "iquest",
//	'\U000000D7': "times",
//	'\U000000F7': "div",
//	'\U000002D8': "breve",
//	'\U000002D9': "dot",
//	'\U000002DA': "ring",
//	'\U000002DB': "ogon",
//	'\U000002DC': "tilde",
//	'\U000002DD': "dblac",
//	'\U000003F6': "bepsi",
//	'\U00002010': "dash",
//	'\U00002013': "ndash",
//	'\U00002014': "mdash",
//	'\U00002015': "horbar",
//	'\U00002016': "vert",
//	'\U00002018': "lsquo",
//	'\U00002019': "rsquo",
//	'\U0000201A': "sbquo",
//	'\U0000201C': "ldquo",
//	'\U0000201D': "rdquo",
//	'\U0000201E': "bdquo",
//	'\U00002020': "dagger",
//	'\U00002021': "dagger",
//	'\U00002022': "bull",
//	'\U00002025': "nldr",
//	'\U00002026': "mldr",
//	'\U00002030': "permil",
//	'\U00002031': "pertenk",
//	'\U00002032': "prime",
//	'\U00002033': "prime",
//	'\U00002034': "tprime",
//	'\U00002035': "bprime",
//	'\U00002039': "lsaquo",
//	'\U0000203A': "rsaquo",
//	'\U0000203E': "oline",
//	'\U00002041': "caret",
//	'\U00002043': "hybull",
//	'\U00002044': "frasl",
//	'\U0000204F': "bsemi",
//	'\U00002057': "qprime",
//	'\U000020AC': "euro",
//	'\U00002105': "incare",
//	'\U00002116': "numero",
//	'\U00002117': "copysr",
//	'\U00002118': "wp",
//	'\U0000211E': "rx",
//	'\U00002122': "trade",
//	'\U00002127': "mho",
//	'\U00002129': "iiota",
//	'\U00002190': "larr",
//	'\U00002191': "uarr",
//	'\U00002192': "rarr",
//	'\U00002193': "darr",
//	'\U00002194': "harr",
//	'\U00002195': "varr",
//	'\U00002196': "nwarr",
//	'\U00002197': "nearr",
//	'\U00002198': "searr",
//	'\U00002199': "swarr",
//	'\U0000219A': "nlarr",
//	'\U0000219B': "nrarr",
//	'\U0000219D': "rarrw",
//	'\U0000219E': "larr",
//	'\U0000219F': "uarr",
//	'\U000021A0': "rarr",
//	'\U000021A1': "darr",
//	'\U000021A2': "larrtl",
//	'\U000021A3': "rarrtl",
//	'\U000021A4': "mapstoleft",
//	'\U000021A5': "mapstoup",
//	'\U000021A6': "map",
//	'\U000021A7': "mapstodown",
//	'\U000021A9': "larrhk",
//	'\U000021AA': "rarrhk",
//	'\U000021AB': "larrlp",
//	'\U000021AC': "rarrlp",
//	'\U000021AD': "harrw",
//	'\U000021AE': "nharr",
//	'\U000021B0': "lsh",
//	'\U000021B1': "rsh",
//	'\U000021B2': "ldsh",
//	'\U000021B3': "rdsh",
//	'\U000021B5': "crarr",
//	'\U000021B6': "cularr",
//	'\U000021B7': "curarr",
//	'\U000021BA': "olarr",
//	'\U000021BB': "orarr",
//	'\U000021BC': "lharu",
//	'\U000021BD': "lhard",
//	'\U000021BE': "uharr",
//	'\U000021BF': "uharl",
//	'\U000021C0': "rharu",
//	'\U000021C1': "rhard",
//	'\U000021C2': "dharr",
//	'\U000021C3': "dharl",
//	'\U000021C4': "rlarr",
//	'\U000021C5': "udarr",
//	'\U000021C6': "lrarr",
//	'\U000021C7': "llarr",
//	'\U000021C8': "uuarr",
//	'\U000021C9': "rrarr",
//	'\U000021CA': "ddarr",
//	'\U000021CB': "lrhar",
//	'\U000021CC': "rlhar",
//	'\U000021CD': "nlarr",
//	'\U000021CE': "nharr",
//	'\U000021CF': "nrarr",
//	'\U000021D0': "larr",
//	'\U000021D1': "uarr",
//	'\U000021D2': "rarr",
//	'\U000021D3': "darr",
//	'\U000021D4': "iff",
//	'\U000021D5': "varr",
//	'\U000021D6': "nwarr",
//	'\U000021D7': "nearr",
//	'\U000021D8': "searr",
//	'\U000021D9': "swarr",
//	'\U000021DA': "laarr",
//	'\U000021DB': "raarr",
//	'\U000021DD': "zigrarr",
//	'\U000021E4': "larrb",
//	'\U000021E5': "rarrb",
//	'\U000021F5': "duarr",
//	'\U000021FD': "loarr",
//	'\U000021FE': "roarr",
//	'\U000021FF': "hoarr",
//	'\U00002200': "forall",
//	'\U00002201': "comp",
//	'\U00002202': "part",
//	'\U00002203': "exist",
//	'\U00002204': "nexist",
//	'\U00002205': "empty",
//	'\U00002207': "del",
//	'\U00002208': "in",
//	'\U00002209': "notin",
//	'\U0000220B': "ni",
//	'\U0000220C': "notni",
//	'\U0000220F': "prod",
//	'\U00002210': "coprod",
//	'\U00002211': "sum",
//	'\U00002212': "minus",
//	'\U00002213': "mp",
//	'\U00002214': "plusdo",
//	'\U00002216': "setmn",
//	'\U00002217': "lowast",
//	'\U00002218': "compfn",
//	'\U0000221A': "sqrt",
//	'\U0000221D': "prop",
//	'\U0000221E': "infin",
//	'\U0000221F': "angrt",
//	'\U00002220': "ang",
//	'\U00002221': "angmsd",
//	'\U00002222': "angsph",
//	'\U00002223': "mid",
//	'\U00002224': "nmid",
//	'\U00002225': "par",
//	'\U00002226': "npar",
//	'\U00002227': "and",
//	'\U00002228': "or",
//	'\U00002229': "cap",
//	'\U0000222A': "cup",
//	'\U0000222B': "int",
//	'\U0000222C': "int",
//	'\U0000222D': "tint",
//	'\U0000222E': "oint",
//	'\U0000222F': "conint",
//	'\U00002230': "cconint",
//	'\U00002231': "cwint",
//	'\U00002232': "cwconint",
//	'\U00002233': "awconint",
//	'\U00002234': "there4",
//	'\U00002235': "becaus",
//	'\U00002236': "ratio",
//	'\U00002237': "colon",
//	'\U00002238': "minusd",
//	'\U0000223A': "mddot",
//	'\U0000223B': "homtht",
//	'\U0000223C': "sim",
//	'\U0000223D': "bsim",
//	'\U0000223E': "ac",
//	'\U0000223F': "acd",
//	'\U00002240': "wr",
//	'\U00002241': "nsim",
//	'\U00002242': "esim",
//	'\U00002243': "sime",
//	'\U00002244': "nsime",
//	'\U00002245': "cong",
//	'\U00002246': "simne",
//	'\U00002247': "ncong",
//	'\U00002248': "ap",
//	'\U00002249': "nap",
//	'\U0000224A': "ape",
//	'\U0000224B': "apid",
//	'\U0000224C': "bcong",
//	'\U0000224D': "cupcap",
//	'\U0000224E': "bump",
//	'\U0000224F': "bumpe",
//	'\U00002250': "doteq",
//	'\U00002251': "edot",
//	'\U00002252': "efdot",
//	'\U00002253': "erdot",
//	'\U00002254': "assign",
//	'\U00002255': "ecolon",
//	'\U00002256': "ecir",
//	'\U00002257': "cire",
//	'\U00002259': "wedgeq",
//	'\U0000225A': "veeeq",
//	'\U0000225C': "trie",
//	'\U0000225F': "equest",
//	'\U00002260': "ne",
//	'\U00002261': "equiv",
//	'\U00002262': "nequiv",
//	'\U00002264': "le",
//	'\U00002265': "ge",
//	'\U00002266': "le",
//	'\U00002267': "ge",
//	'\U00002268': "lne",
//	'\U00002269': "gne",
//	'\U0000226A': "lt",
//	'\U0000226B': "gt",
//	'\U0000226C': "twixt",
//	'\U0000226D': "notcupcap",
//	'\U0000226E': "nlt",
//	'\U0000226F': "ngt",
//	'\U00002270': "nle",
//	'\U00002271': "nge",
//	'\U00002272': "lsim",
//	'\U00002273': "gsim",
//	'\U00002274': "nlsim",
//	'\U00002275': "ngsim",
//	'\U00002276': "lg",
//	'\U00002277': "gl",
//	'\U00002278': "ntlg",
//	'\U00002279': "ntgl",
//	'\U0000227A': "pr",
//	'\U0000227B': "sc",
//	'\U0000227C': "prcue",
//	'\U0000227D': "sccue",
//	'\U0000227E': "prsim",
//	'\U0000227F': "scsim",
//	'\U00002280': "npr",
//	'\U00002281': "nsc",
//	'\U00002282': "sub",
//	'\U00002283': "sup",
//	'\U00002284': "nsub",
//	'\U00002285': "nsup",
//	'\U00002286': "sube",
//	'\U00002287': "supe",
//	'\U00002288': "nsube",
//	'\U00002289': "nsupe",
//	'\U0000228A': "subne",
//	'\U0000228B': "supne",
//	'\U0000228D': "cupdot",
//	'\U0000228E': "uplus",
//	'\U0000228F': "sqsub",
//	'\U00002290': "sqsup",
//	'\U00002291': "sqsube",
//	'\U00002292': "sqsupe",
//	'\U00002293': "sqcap",
//	'\U00002294': "sqcup",
//	'\U00002295': "oplus",
//	'\U00002296': "ominus",
//	'\U00002297': "otimes",
//	'\U00002298': "osol",
//	'\U00002299': "odot",
//	'\U0000229A': "ocir",
//	'\U0000229B': "oast",
//	'\U0000229D': "odash",
//	'\U0000229E': "plusb",
//	'\U0000229F': "minusb",
//	'\U000022A0': "timesb",
//	'\U000022A1': "sdotb",
//	'\U000022A2': "vdash",
//	'\U000022A3': "dashv",
//	'\U000022A4': "top",
//	'\U000022A5': "bot",
//	'\U000022A7': "models",
//	'\U000022A8': "vdash",
//	'\U000022A9': "vdash",
//	'\U000022AA': "vvdash",
//	'\U000022AB': "vdash",
//	'\U000022AC': "nvdash",
//	'\U000022AD': "nvdash",
//	'\U000022AE': "nvdash",
//	'\U000022AF': "nvdash",
//	'\U000022B0': "prurel",
//	'\U000022B2': "vltri",
//	'\U000022B3': "vrtri",
//	'\U000022B4': "ltrie",
//	'\U000022B5': "rtrie",
//	'\U000022B6': "origof",
//	'\U000022B7': "imof",
//	'\U000022B8': "mumap",
//	'\U000022B9': "hercon",
//	'\U000022BA': "intcal",
//	'\U000022BB': "veebar",
//	'\U000022BD': "barvee",
//	'\U000022BE': "angrtvb",
//	'\U000022BF': "lrtri",
//	'\U000022C0': "wedge",
//	'\U000022C1': "vee",
//	'\U000022C2': "xcap",
//	'\U000022C3': "xcup",
//	'\U000022C4': "diam",
//	'\U000022C5': "sdot",
//	'\U000022C6': "star",
//	'\U000022C7': "divonx",
//	'\U000022C8': "bowtie",
//	'\U000022C9': "ltimes",
//	'\U000022CA': "rtimes",
//	'\U000022CB': "lthree",
//	'\U000022CC': "rthree",
//	'\U000022CD': "bsime",
//	'\U000022CE': "cuvee",
//	'\U000022CF': "cuwed",
//	'\U000022D0': "sub",
//	'\U000022D1': "sup",
//	'\U000022D2': "cap",
//	'\U000022D3': "cup",
//	'\U000022D4': "fork",
//	'\U000022D5': "epar",
//	'\U000022D6': "ltdot",
//	'\U000022D7': "gtdot",
//	'\U000022D8': "ll",
//	'\U000022D9': "gg",
//	'\U000022DA': "leg",
//	'\U000022DB': "gel",
//	'\U000022DE': "cuepr",
//	'\U000022DF': "cuesc",
//	'\U000022E0': "nprcue",
//	'\U000022E1': "nsccue",
//	'\U000022E2': "nsqsube",
//	'\U000022E3': "nsqsupe",
//	'\U000022E6': "lnsim",
//	'\U000022E7': "gnsim",
//	'\U000022E8': "prnsim",
//	'\U000022E9': "scnsim",
//	'\U000022EA': "nltri",
//	'\U000022EB': "nrtri",
//	'\U000022EC': "nltrie",
//	'\U000022ED': "nrtrie",
//	'\U000022EE': "vellip",
//	'\U000022EF': "ctdot",
//	'\U000022F0': "utdot",
//	'\U000022F1': "dtdot",
//	'\U000022F2': "disin",
//	'\U000022F3': "isinsv",
//	'\U000022F4': "isins",
//	'\U000022F5': "isindot",
//	'\U000022F6': "notinvc",
//	'\U000022F7': "notinvb",
//	'\U000022F9': "isine",
//	'\U000022FA': "nisd",
//	'\U000022FB': "xnis",
//	'\U000022FC': "nis",
//	'\U000022FD': "notnivc",
//	'\U000022FE': "notnivb",
//	'\U00002305': "barwed",
//	'\U00002306': "barwed",
//	'\U00002308': "lceil",
//	'\U00002309': "rceil",
//	'\U0000230A': "lfloor",
//	'\U0000230B': "rfloor",
//	'\U0000230C': "drcrop",
//	'\U0000230D': "dlcrop",
//	'\U0000230E': "urcrop",
//	'\U0000230F': "ulcrop",
//	'\U00002310': "bnot",
//	'\U00002312': "profline",
//	'\U00002313': "profsurf",
//	'\U00002315': "telrec",
//	'\U00002316': "target",
//	'\U0000231C': "ulcorn",
//	'\U0000231D': "urcorn",
//	'\U0000231E': "dlcorn",
//	'\U0000231F': "drcorn",
//	'\U00002322': "frown",
//	'\U00002323': "smile",
//	'\U0000232D': "cylcty",
//	'\U0000232E': "profalar",
//	'\U00002336': "topbot",
//	'\U0000233D': "ovbar",
//	'\U0000233F': "solbar",
//	'\U0000237C': "angzarr",
//	'\U000023B0': "lmoust",
//	'\U000023B1': "rmoust",
//	'\U000023B4': "tbrk",
//	'\U000023B5': "bbrk",
//	'\U000023B6': "bbrktbrk",
//	'\U000023DC': "overparenthesis",
//	'\U000023DD': "underparenthesis",
//	'\U000023DE': "overbrace",
//	'\U000023DF': "underbrace",
//	'\U000023E2': "trpezium",
//	'\U000023E7': "elinters",
//	'\U00002423': "blank",
//	'\U000024C8': "os",
//	'\U00002500': "boxh",
//	'\U00002502': "boxv",
//	'\U0000250C': "boxdr",
//	'\U00002510': "boxdl",
//	'\U00002514': "boxur",
//	'\U00002518': "boxul",
//	'\U0000251C': "boxvr",
//	'\U00002524': "boxvl",
//	'\U0000252C': "boxhd",
//	'\U00002534': "boxhu",
//	'\U0000253C': "boxvh",
//	'\U00002550': "boxh",
//	'\U00002551': "boxv",
//	'\U00002552': "boxdr",
//	'\U00002553': "boxdr",
//	'\U00002554': "boxdr",
//	'\U00002555': "boxdl",
//	'\U00002556': "boxdl",
//	'\U00002557': "boxdl",
//	'\U00002558': "boxur",
//	'\U00002559': "boxur",
//	'\U0000255A': "boxur",
//	'\U0000255B': "boxul",
//	'\U0000255C': "boxul",
//	'\U0000255D': "boxul",
//	'\U0000255E': "boxvr",
//	'\U0000255F': "boxvr",
//	'\U00002560': "boxvr",
//	'\U00002561': "boxvl",
//	'\U00002562': "boxvl",
//	'\U00002563': "boxvl",
//	'\U00002564': "boxhd",
//	'\U00002565': "boxhd",
//	'\U00002566': "boxhd",
//	'\U00002567': "boxhu",
//	'\U00002568': "boxhu",
//	'\U00002569': "boxhu",
//	'\U0000256A': "boxvh",
//	'\U0000256B': "boxvh",
//	'\U0000256C': "boxvh",
//	'\U00002580': "uhblk",
//	'\U00002584': "lhblk",
//	'\U00002588': "block",
//	'\U00002591': "blk14",
//	'\U00002592': "blk12",
//	'\U00002593': "blk34",
//	'\U000025A1': "squ",
//	'\U000025AA': "squf",
//	'\U000025AB': "emptyverysmallsquare",
//	'\U000025AD': "rect",
//	'\U000025AE': "marker",
//	'\U000025B1': "fltns",
//	'\U000025B3': "xutri",
//	'\U000025B4': "utrif",
//	'\U000025B5': "utri",
//	'\U000025B8': "rtrif",
//	'\U000025B9': "rtri",
//	'\U000025BD': "xdtri",
//	'\U000025BE': "dtrif",
//	'\U000025BF': "dtri",
//	'\U000025C2': "ltrif",
//	'\U000025C3': "ltri",
//	'\U000025CA': "loz",
//	'\U000025CB': "cir",
//	'\U000025EC': "tridot",
//	'\U000025EF': "xcirc",
//	'\U000025F8': "ultri",
//	'\U000025F9': "urtri",
//	'\U000025FA': "lltri",
//	'\U000025FB': "emptysmallsquare",
//	'\U000025FC': "filledsmallsquare",
//	'\U00002605': "starf",
//	'\U00002606': "star",
//	'\U0000260E': "phone",
//	'\U00002640': "female",
//	'\U00002642': "male",
//	'\U00002660': "spades",
//	'\U00002663': "clubs",
//	'\U00002665': "hearts",
//	'\U00002666': "diams",
//	'\U0000266A': "sung",
//	'\U0000266D': "flat",
//	'\U0000266E': "natur",
//	'\U0000266F': "sharp",
//	'\U00002713': "check",
//	'\U00002717': "cross",
//	'\U00002720': "malt",
//	'\U00002736': "sext",
//	'\U00002758': "verticalseparator",
//	'\U00002772': "lbbrk",
//	'\U00002773': "rbbrk",
//	'\U000027C8': "bsolhsub",
//	'\U000027C9': "suphsol",
//	'\U000027E6': "lobrk",
//	'\U000027E7': "robrk",
//	'\U000027E8': "lang",
//	'\U000027E9': "rang",
//	'\U000027EA': "lang",
//	'\U000027EB': "rang",
//	'\U000027EC': "loang",
//	'\U000027ED': "roang",
//	'\U000027F5': "xlarr",
//	'\U000027F6': "xrarr",
//	'\U000027F7': "xharr",
//	'\U000027F8': "xlarr",
//	'\U000027F9': "xrarr",
//	'\U000027FA': "xharr",
//	'\U000027FC': "xmap",
//	'\U000027FF': "dzigrarr",
//	'\U00002902': "nvlarr",
//	'\U00002903': "nvrarr",
//	'\U00002904': "nvharr",
//	'\U00002905': "map",
//	'\U0000290C': "lbarr",
//	'\U0000290D': "rbarr",
//	'\U0000290E': "lbarr",
//	'\U0000290F': "rbarr",
//	'\U00002910': "rbarr",
//	'\U00002911': "ddotrahd",
//	'\U00002912': "uparrowbar",
//	'\U00002913': "downarrowbar",
//	'\U00002916': "rarrtl",
//	'\U00002919': "latail",
//	'\U0000291A': "ratail",
//	'\U0000291B': "latail",
//	'\U0000291C': "ratail",
//	'\U0000291D': "larrfs",
//	'\U0000291E': "rarrfs",
//	'\U0000291F': "larrbfs",
//	'\U00002920': "rarrbfs",
//	'\U00002923': "nwarhk",
//	'\U00002924': "nearhk",
//	'\U00002925': "searhk",
//	'\U00002926': "swarhk",
//	'\U00002927': "nwnear",
//	'\U00002928': "toea",
//	'\U00002929': "tosa",
//	'\U0000292A': "swnwar",
//	'\U00002933': "rarrc",
//	'\U00002935': "cudarrr",
//	'\U00002936': "ldca",
//	'\U00002937': "rdca",
//	'\U00002938': "cudarrl",
//	'\U00002939': "larrpl",
//	'\U0000293C': "curarrm",
//	'\U0000293D': "cularrp",
//	'\U00002945': "rarrpl",
//	'\U00002948': "harrcir",
//	'\U00002949': "uarrocir",
//	'\U0000294A': "lurdshar",
//	'\U0000294B': "ldrushar",
//	'\U0000294E': "leftrightvector",
//	'\U0000294F': "rightupdownvector",
//	'\U00002950': "downleftrightvector",
//	'\U00002951': "leftupdownvector",
//	'\U00002952': "leftvectorbar",
//	'\U00002953': "rightvectorbar",
//	'\U00002954': "rightupvectorbar",
//	'\U00002955': "rightdownvectorbar",
//	'\U00002956': "downleftvectorbar",
//	'\U00002957': "downrightvectorbar",
//	'\U00002958': "leftupvectorbar",
//	'\U00002959': "leftdownvectorbar",
//	'\U0000295A': "leftteevector",
//	'\U0000295B': "rightteevector",
//	'\U0000295C': "rightupteevector",
//	'\U0000295D': "rightdownteevector",
//	'\U0000295E': "downleftteevector",
//	'\U0000295F': "downrightteevector",
//	'\U00002960': "leftupteevector",
//	'\U00002961': "leftdownteevector",
//	'\U00002962': "lhar",
//	'\U00002963': "uhar",
//	'\U00002964': "rhar",
//	'\U00002965': "dhar",
//	'\U00002966': "luruhar",
//	'\U00002967': "ldrdhar",
//	'\U00002968': "ruluhar",
//	'\U00002969': "rdldhar",
//	'\U0000296A': "lharul",
//	'\U0000296B': "llhard",
//	'\U0000296C': "rharul",
//	'\U0000296D': "lrhard",
//	'\U0000296E': "udhar",
//	'\U0000296F': "duhar",
//	'\U00002970': "roundimplies",
//	'\U00002971': "erarr",
//	'\U00002972': "simrarr",
//	'\U00002973': "larrsim",
//	'\U00002974': "rarrsim",
//	'\U00002975': "rarrap",
//	'\U00002976': "ltlarr",
//	'\U00002978': "gtrarr",
//	'\U00002979': "subrarr",
//	'\U0000297B': "suplarr",
//	'\U0000297C': "lfisht",
//	'\U0000297D': "rfisht",
//	'\U0000297E': "ufisht",
//	'\U0000297F': "dfisht",
//	'\U00002985': "lopar",
//	'\U00002986': "ropar",
//	'\U0000298B': "lbrke",
//	'\U0000298C': "rbrke",
//	'\U0000298D': "lbrkslu",
//	'\U0000298E': "rbrksld",
//	'\U0000298F': "lbrksld",
//	'\U00002990': "rbrkslu",
//	'\U00002991': "langd",
//	'\U00002992': "rangd",
//	'\U00002993': "lparlt",
//	'\U00002994': "rpargt",
//	'\U00002995': "gtlpar",
//	'\U00002996': "ltrpar",
//	'\U0000299A': "vzigzag",
//	'\U0000299C': "vangrt",
//	'\U0000299D': "angrtvbd",
//	'\U000029A4': "ange",
//	'\U000029A5': "range",
//	'\U000029A6': "dwangle",
//	'\U000029A7': "uwangle",
//	'\U000029A8': "angmsdaa",
//	'\U000029A9': "angmsdab",
//	'\U000029AA': "angmsdac",
//	'\U000029AB': "angmsdad",
//	'\U000029AC': "angmsdae",
//	'\U000029AD': "angmsdaf",
//	'\U000029AE': "angmsdag",
//	'\U000029AF': "angmsdah",
//	'\U000029B0': "bemptyv",
//	'\U000029B1': "demptyv",
//	'\U000029B2': "cemptyv",
//	'\U000029B3': "raemptyv",
//	'\U000029B4': "laemptyv",
//	'\U000029B5': "ohbar",
//	'\U000029B6': "omid",
//	'\U000029B7': "opar",
//	'\U000029B9': "operp",
//	'\U000029BB': "olcross",
//	'\U000029BC': "odsold",
//	'\U000029BE': "olcir",
//	'\U000029BF': "ofcir",
//	'\U000029C0': "olt",
//	'\U000029C1': "ogt",
//	'\U000029C2': "cirscir",
//	'\U000029C3': "cire",
//	'\U000029C4': "solb",
//	'\U000029C5': "bsolb",
//	'\U000029C9': "boxbox",
//	'\U000029CD': "trisb",
//	'\U000029CE': "rtriltri",
//	'\U000029CF': "lefttrianglebar",
//	'\U000029D0': "righttrianglebar",
//	'\U000029DC': "iinfin",
//	'\U000029DD': "infintie",
//	'\U000029DE': "nvinfin",
//	'\U000029E3': "eparsl",
//	'\U000029E4': "smeparsl",
//	'\U000029E5': "eqvparsl",
//	'\U000029EB': "lozf",
//	'\U000029F4': "ruledelayed",
//	'\U000029F6': "dsol",
//	'\U00002A00': "xodot",
//	'\U00002A01': "xoplus",
//	'\U00002A02': "xotime",
//	'\U00002A04': "xuplus",
//	'\U00002A06': "xsqcup",
//	'\U00002A0C': "qint",
//	'\U00002A0D': "fpartint",
//	'\U00002A10': "cirfnint",
//	'\U00002A11': "awint",
//	'\U00002A12': "rppolint",
//	'\U00002A13': "scpolint",
//	'\U00002A14': "npolint",
//	'\U00002A15': "pointint",
//	'\U00002A16': "quatint",
//	'\U00002A17': "intlarhk",
//	'\U00002A22': "pluscir",
//	'\U00002A23': "plusacir",
//	'\U00002A24': "simplus",
//	'\U00002A25': "plusdu",
//	'\U00002A26': "plussim",
//	'\U00002A27': "plustwo",
//	'\U00002A29': "mcomma",
//	'\U00002A2A': "minusdu",
//	'\U00002A2D': "loplus",
//	'\U00002A2E': "roplus",
//	'\U00002A2F': "cross",
//	'\U00002A30': "timesd",
//	'\U00002A31': "timesbar",
//	'\U00002A33': "smashp",
//	'\U00002A34': "lotimes",
//	'\U00002A35': "rotimes",
//	'\U00002A36': "otimesas",
//	'\U00002A37': "otimes",
//	'\U00002A38': "odiv",
//	'\U00002A39': "triplus",
//	'\U00002A3A': "triminus",
//	'\U00002A3B': "tritime",
//	'\U00002A3C': "iprod",
//	'\U00002A3F': "amalg",
//	'\U00002A40': "capdot",
//	'\U00002A42': "ncup",
//	'\U00002A43': "ncap",
//	'\U00002A44': "capand",
//	'\U00002A45': "cupor",
//	'\U00002A46': "cupcap",
//	'\U00002A47': "capcup",
//	'\U00002A48': "cupbrcap",
//	'\U00002A49': "capbrcup",
//	'\U00002A4A': "cupcup",
//	'\U00002A4B': "capcap",
//	'\U00002A4C': "ccups",
//	'\U00002A4D': "ccaps",
//	'\U00002A50': "ccupssm",
//	'\U00002A53': "and",
//	'\U00002A54': "or",
//	'\U00002A55': "andand",
//	'\U00002A56': "oror",
//	'\U00002A57': "orslope",
//	'\U00002A58': "andslope",
//	'\U00002A5A': "andv",
//	'\U00002A5B': "orv",
//	'\U00002A5C': "andd",
//	'\U00002A5D': "ord",
//	'\U00002A5F': "wedbar",
//	'\U00002A66': "sdote",
//	'\U00002A6A': "simdot",
//	'\U00002A6D': "congdot",
//	'\U00002A6E': "easter",
//	'\U00002A6F': "apacir",
//	'\U00002A70': "ape",
//	'\U00002A71': "eplus",
//	'\U00002A72': "pluse",
//	'\U00002A73': "esim",
//	'\U00002A74': "colone",
//	'\U00002A75': "equal",
//	'\U00002A77': "eddot",
//	'\U00002A78': "equivdd",
//	'\U00002A79': "ltcir",
//	'\U00002A7A': "gtcir",
//	'\U00002A7B': "ltquest",
//	'\U00002A7C': "gtquest",
//	'\U00002A7D': "les",
//	'\U00002A7E': "ges",
//	'\U00002A7F': "lesdot",
//	'\U00002A80': "gesdot",
//	'\U00002A81': "lesdoto",
//	'\U00002A82': "gesdoto",
//	'\U00002A83': "lesdotor",
//	'\U00002A84': "gesdotol",
//	'\U00002A85': "lap",
//	'\U00002A86': "gap",
//	'\U00002A87': "lne",
//	'\U00002A88': "gne",
//	'\U00002A89': "lnap",
//	'\U00002A8A': "gnap",
//	'\U00002A8B': "leg",
//	'\U00002A8C': "gel",
//	'\U00002A8D': "lsime",
//	'\U00002A8E': "gsime",
//	'\U00002A8F': "lsimg",
//	'\U00002A90': "gsiml",
//	'\U00002A91': "lge",
//	'\U00002A92': "gle",
//	'\U00002A93': "lesges",
//	'\U00002A94': "gesles",
//	'\U00002A95': "els",
//	'\U00002A96': "egs",
//	'\U00002A97': "elsdot",
//	'\U00002A98': "egsdot",
//	'\U00002A99': "el",
//	'\U00002A9A': "eg",
//	'\U00002A9D': "siml",
//	'\U00002A9E': "simg",
//	'\U00002A9F': "simle",
//	'\U00002AA0': "simge",
//	'\U00002AA1': "lessless",
//	'\U00002AA2': "greatergreater",
//	'\U00002AA4': "glj",
//	'\U00002AA5': "gla",
//	'\U00002AA6': "ltcc",
//	'\U00002AA7': "gtcc",
//	'\U00002AA8': "lescc",
//	'\U00002AA9': "gescc",
//	'\U00002AAA': "smt",
//	'\U00002AAB': "lat",
//	'\U00002AAC': "smte",
//	'\U00002AAD': "late",
//	'\U00002AAE': "bumpe",
//	'\U00002AAF': "pre",
//	'\U00002AB0': "sce",
//	'\U00002AB3': "pre",
//	'\U00002AB4': "sce",
//	'\U00002AB5': "prne",
//	'\U00002AB6': "scne",
//	'\U00002AB7': "prap",
//	'\U00002AB8': "scap",
//	'\U00002AB9': "prnap",
//	'\U00002ABA': "scnap",
//	'\U00002ABB': "pr",
//	'\U00002ABC': "sc",
//	'\U00002ABD': "subdot",
//	'\U00002ABE': "supdot",
//	'\U00002ABF': "subplus",
//	'\U00002AC0': "supplus",
//	'\U00002AC1': "submult",
//	'\U00002AC2': "supmult",
//	'\U00002AC3': "subedot",
//	'\U00002AC4': "supedot",
//	'\U00002AC5': "sube",
//	'\U00002AC6': "supe",
//	'\U00002AC7': "subsim",
//	'\U00002AC8': "supsim",
//	'\U00002ACB': "subne",
//	'\U00002ACC': "supne",
//	'\U00002ACF': "csub",
//	'\U00002AD0': "csup",
//	'\U00002AD1': "csube",
//	'\U00002AD2': "csupe",
//	'\U00002AD3': "subsup",
//	'\U00002AD4': "supsub",
//	'\U00002AD5': "subsub",
//	'\U00002AD6': "supsup",
//	'\U00002AD7': "suphsub",
//	'\U00002AD8': "supdsub",
//	'\U00002AD9': "forkv",
//	'\U00002ADA': "topfork",
//	'\U00002ADB': "mlcp",
//	'\U00002AE4': "dashv",
//	'\U00002AE6': "vdashl",
//	'\U00002AE7': "barv",
//	'\U00002AE8': "vbar",
//	'\U00002AE9': "vbarv",
//	'\U00002AEB': "vbar",
//	'\U00002AEC': "not",
//	'\U00002AED': "bnot",
//	'\U00002AEE': "rnmid",
//	'\U00002AEF': "cirmid",
//	'\U00002AF0': "midcir",
//	'\U00002AF1': "topcir",
//	'\U00002AF2': "nhpar",
//	'\U00002AF3': "parsim",
//	'\U00002AFD': "parsl",
//}