package test

var EscapeCases = map[string]string{
	"3": `3`,
	"The quick brown fox jumped over the lazy dog.":    `The quick brown fox jumped over the lazy dog.`,
	"Control characters:  \b, \f, \n, \r, \t":          `Control characters:  \b, \f, \n, \r, \t`,
	"Quote and slashes: \", \\, /":                     `Quote and slashes: \", \\, \/`,
	"UTF8 Characters: ϢӦֆĒ͖̈́Ͳ  ĦĪǂǼɆψϠѬӜԪ":             `UTF8 Characters: ϢӦֆĒ͖̈́Ͳ  ĦĪǂǼɆψϠѬӜԪ`,
	"Unicode Characters: 😀🐦‍🔥⛓️‍💥🍋‍🟩  ظۇ  ❂✈☯  亳亴亵亶亷亸": `Unicode Characters: 😀🐦‍🔥⛓️‍💥🍋‍🟩  ظۇ  ❂✈☯  亳亴亵亶亷亸`,
}
