package score

import (
	_ "embed"
	"fmt"
	"html/template"
	"sort"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/test"
)

type KeeperTag string

type Keeper struct {
	tag  KeeperTag
	x, y Axis
	doc  template.HTML
}

func NewKeeper(tag KeeperTag, x, y Axis, doc template.HTML) *Keeper {
	return &Keeper{
		tag: tag,
		x:   x,
		y:   y,
		doc: doc,
	}
}

func (k *Keeper) Setup(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := k.x.Setup(bench, warns); err != nil {
		return fmt.Errorf("initialize x: %w", err)
	}
	if err := k.y.Setup(bench, warns); err != nil {
		return fmt.Errorf("initialize y: %w", err)
	}
	// TODO: Do something else?
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

// Documentation returns documentation related to the current scorekeeper object.
func (k *Keeper) Documentation() template.HTML {
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
	test.Debugf(1, ">>> AddKeeper(%s)", keeper.Tag())
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
			return keeperTags[i] < keeperTags[j]
		})
	}
	return keeperTags
}
