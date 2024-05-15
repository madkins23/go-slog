package score

import (
	"fmt"
	"html/template"
	"sort"

	"github.com/madkins23/go-slog/internal/data"
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

func (k *Keeper) X() Axis {
	return k.x
}

func (k *Keeper) Y() Axis {
	return k.y
}

func (k *Keeper) Tag() KeeperTag {
	if k == nil {
		return ""
	}
	return k.tag
}

func (k *Keeper) Documentation() template.HTML {
	return k.doc
}

func (k *Keeper) Exhibits() []Exhibit {
	var exhibits []Exhibit
	xExhibits := k.X().Exhibits()
	yExhibits := k.Y().Exhibits()
	size := len(xExhibits) + len(yExhibits)
	if size > 0 {
		exhibits = make([]Exhibit, 0, size)
		for _, exhibit := range xExhibits {
			exhibits = append(exhibits, exhibit)
		}
		for _, exhibit := range yExhibits {
			exhibits = append(exhibits, exhibit)
		}
	}
	return exhibits
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
