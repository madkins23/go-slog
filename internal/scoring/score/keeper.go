package score

import (
	_ "embed"
	"fmt"
	"html/template"
	"sort"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
)

const specialChar = '+'

type KeeperTag string

type Keeper struct {
	tag      KeeperTag
	x, y     Axis
	filter   Filter
	handlers []data.HandlerTag
	tests    []data.TestTag
	doc      template.HTML
	KeeperOptions
}

type KeeperOptions struct {
	ChartCaption, Title template.HTML
}

const (
	noChartCaption = "No chart caption yet!!!"
	noKeeperTitle  = "No keeper title yet!!!"
)

func NewKeeper(tag KeeperTag, x, y Axis, doc template.HTML, options *KeeperOptions, filter Filter) *Keeper {
	k := &Keeper{
		tag:    tag,
		x:      x,
		y:      y,
		filter: filter,
		doc:    doc,
	}
	if options != nil {
		k.KeeperOptions = *options
	}
	if k.KeeperOptions.ChartCaption == "" {
		k.KeeperOptions.ChartCaption = noChartCaption
	}
	if k.KeeperOptions.Title == "" {
		k.KeeperOptions.Title = noKeeperTitle
	}
	return k
}

func (k *Keeper) Setup(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := k.x.Setup(bench, warns); err != nil {
		return fmt.Errorf("initialize x: %w", err)
	}
	if err := k.y.Setup(bench, warns); err != nil {
		return fmt.Errorf("initialize y: %w", err)
	}
	k.handlers = make([]data.HandlerTag, 0)
	k.tests = make([]data.TestTag, 0)
	for _, hdlr := range bench.HandlerTags() {
		if k.filter == nil || k.filter.Keep(warns.HandlerName(hdlr)) {
			k.handlers = append(k.handlers, hdlr)
		}
	}
	for _, test := range bench.TestTags() {
		if bench.HasTest(test) {
			if k.filter == nil || k.filter.Keep(warns.TestName(test)) {
				k.tests = append(k.tests, test)
			}
		}
	}
	return nil
}

func (k *Keeper) Name() string {
	return string(k.tag)
}

func (k *Keeper) Tag() KeeperTag {
	if k == nil {
		return ""
	}
	return k.tag
}

func (k *Keeper) HandlerTags() []data.HandlerTag {
	return k.handlers
}

func (k *Keeper) TestTags() []data.TestTag {
	return k.tests
}

func (k *Keeper) ChartCaption() template.HTML {
	return k.KeeperOptions.ChartCaption
}

func (k *Keeper) ChartTitle() template.HTML {
	return k.KeeperOptions.Title
}

// Summary returns documentation related to the current scorekeeper object.
func (k *Keeper) Summary() template.HTML {
	return k.doc
}

var (
	//go:embed doc/overview.md
	overviewMD   string
	overviewHTML template.HTML
)

// Overview returns documentation applicable to all scorekeepers.
func (k *Keeper) Overview() template.HTML {
	if overviewHTML == "" {
		overviewHTML = markdown.TemplateHTML(overviewMD, false)
	}
	return overviewHTML
}

func (k *Keeper) Axes() map[string]Axis {
	return map[string]Axis{"X": k.x, "Y": k.y}
}

func (k *Keeper) X() Axis {
	return k.x
}

func (k *Keeper) Y() Axis {
	return k.y
}

// -----------------------------------------------------------------------------

var (
	keepers    map[KeeperTag]*Keeper
	keeperTags []KeeperTag
)

func AddKeeper(keeper *Keeper) error {
	if keepers == nil {
		keepers = make(map[KeeperTag]*Keeper)
	}
	if _, found := keepers[keeper.tag]; found {
		return fmt.Errorf("duplicate keeper '%s'", keeper.tag)
	} else {
		keepers[keeper.tag] = keeper
		keeperTags = nil
	}
	return nil
}

func GetKeeper(tag KeeperTag) *Keeper {
	return keepers[tag]
}

func Keepers() []KeeperTag {
	if keeperTags == nil {
		keeperTags = make([]KeeperTag, 0, len(keepers))
		for tag := range keepers {
			keeperTags = append(keeperTags, tag)
		}
		sort.Slice(keeperTags, func(i, j int) bool {
			if keeperTags[i][0] == specialChar {
				if keeperTags[j][0] != specialChar {
					return false
				}
			} else if keeperTags[j][0] == specialChar {
				return true
			}
			return keeperTags[i] < keeperTags[j]
		})
	}
	return keeperTags
}
