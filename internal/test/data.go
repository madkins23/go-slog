package test

import (
	"log/slog"
	"math"
	"net"
	"time"

	"github.com/madkins23/go-slog/infra"
)

const (
	Message = "This is a message. No, really!"
)

var (
	anything = []any{"alpha", "omega"}
	duration = time.Minute + 3*time.Second + 14*time.Millisecond
	ip       = net.IPv4(123, 231, 213, 23)
	ipNet    = &net.IPNet{IP: ip, Mask: []byte{0xFF, 0xFF, 0xFF, 0x80}}
	Level    = slog.LevelInfo
	mac, _   = net.ParseMAC("00:00:5e:00:53:01") // Hopefully no errors here.
	Now      = time.Now()
)

var (
	Attributes = []slog.Attr{
		slog.Time("when", Now),
		slog.Duration("howLong", duration),
		slog.String("Goober", "Snoofus"),
		slog.Bool("boolean", true),
		slog.Float64("pi", math.Pi),
		slog.Int("skidoo", 23),
		slog.Int64("minus", -64),
		slog.Uint64("unsigned", 79),
		slog.Any("any", anything),
		slog.Any("ip", ip),
		slog.Any("ipNet", ipNet),
		slog.Any("macAddr", mac),
		slog.Group("group",
			slog.String("name", "Beatles"),
			infra.EmptyAttr(),
			slog.Float64("pi", math.Pi),
			infra.EmptyAttr(),
			slog.Group("subGroup",
				infra.EmptyAttr(),
				slog.String("name", "Rolling Stones"),
				infra.EmptyAttr()))}
	AttributeMap = map[string]any{
		"howLong":  float64(duration),
		"when":     Now.Format(time.RFC3339Nano),
		"Goober":   "Snoofus",
		"boolean":  true,
		"pi":       math.Pi,
		"skidoo":   float64(23),
		"minus":    float64(-64),
		"unsigned": float64(79),
		"any":      anything,
		"ip":       ip.String(),
		"ipNet":    ipNet.String(),
		"macAddr":  mac.String(),
		"group": map[string]any{
			"name": "Beatles",
			"pi":   math.Pi,
			"subGroup": map[string]any{
				"name": "Rolling Stones",
			},
		},
	}
)

var EscapeCases = map[string]string{
	"3": `3`,
	"The quick brown fox jumped over the lazy dog.":    `The quick brown fox jumped over the lazy dog.`,
	"Control characters:  \b, \f, \n, \r, \t":          `Control characters:  \b, \f, \n, \r, \t`,
	"Quote and slashes: \", \\, /":                     `Quote and slashes: \", \\, \/`,
	"UTF8 Characters: ϢӦֆĒ͖̈́Ͳ  ĦĪǂǼɆψϠѬӜԪ":             `UTF8 Characters: ϢӦֆĒ͖̈́Ͳ  ĦĪǂǼɆψϠѬӜԪ`,
	"Unicode Characters: 😀🐦‍🔥⛓️‍💥🍋‍🟩  ظۇ  ❂✈☯  亳亴亵亶亷亸": `Unicode Characters: 😀🐦‍🔥⛓️‍💥🍋‍🟩  ظۇ  ❂✈☯  亳亴亵亶亷亸`,
}
