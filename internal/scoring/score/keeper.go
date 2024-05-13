package score

import (
	"fmt"
	"html/template"
	"log/slog"
	"sort"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/test"
)

type KeeperTag string

type Keeper struct {
	tag  KeeperTag
	x, y Axis
}

func NewKeeper(tag KeeperTag, x, y Axis) *Keeper {
	return &Keeper{
		tag: tag,
		x:   x,
		y:   y,
	}
}

func (k *Keeper) Initialize(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := k.x.Initialize(bench, warns); err != nil {
		return fmt.Errorf("initialize x: %w", err)
	}
	if err := k.y.Initialize(bench, warns); err != nil {
		return fmt.Errorf("initialize y: %w", err)
	}
	// TODO: Do something else?
	return nil
}

func (k *Keeper) DocOverview() template.HTML {
	slog.Error("TBD", "method", "DocOverview")
	return "<strong>TBD</strong>"
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

// Initialize all score.Keeper objects with the specified benchmark and marking data.
//
// Do not call this from within an init() function,
// it is dependent on configuration done during init() functions
// defined in package internal/scoring/keeper.
func Initialize(bench *data.Benchmarks, warns *data.Warnings) error {
	for name, keeper := range keepers {
		if err := keeper.Initialize(bench, warns); err != nil {
			return fmt.Errorf("initialize '%s': %w", name, err)
		}
	}
	return nil
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
