package warn

import (
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

type HandlerData struct {
	*common.HandlerData
	byLevel map[warning.Level]*score.Average
	counts  map[warning.Level]uint64
}

func NewHandlerData() *HandlerData {
	hd := &HandlerData{
		HandlerData: common.NewHandlerData(),
		byLevel:     make(map[warning.Level]*score.Average),
		counts:      make(map[warning.Level]uint64),
	}
	return hd
}

func (hd *HandlerData) CountFor(level warning.Level) uint64 {
	return hd.counts[level]
}

func (hd *HandlerData) SetCountFor(level warning.Level, count uint64) {
	hd.counts[level] = count
}

func (hd *HandlerData) ByLevel(level warning.Level) *score.Average {
	if hd.byLevel[level] == nil {
		hd.byLevel[level] = &score.Average{}
	}
	return hd.byLevel[level]
}

/*
   Detail Table
     per hdlr
       sub score
         average of scores from all tests of hdlr
         show after test columns before 'regular' scores
         may be just another score type
         - $axis.Score(score.SubScore)
         - scores[subscore]
     per hdlr/test
       score
         converted score from counts for all levels for hdlr/test
         - $axis.ByTest(test)
         - byTest[test].score
   Levels Table
     per hdlr
       sub score
         average of scores from all levels of hdlr
         show after test columns before 'regular' scores
         may be just another score type
         - $axis.Score(score.SubScore)
         - scores[subscore]
     per hdlr/level
       score
         converted score from counts for specified level for all hdlr tests
         show after test columns before 'regular' scores
         may be just another score type
         - $warn.ScoreForLevel(hdlr, level)
   Specific Level Table
     per hdlr
       sub score
         average of scores from all tests of handler for specified level
         show after test columns before 'regular' scores
         may be just another score type
         - $axis.Score(score.SubScore)
         - scores[subscore]
     per hdlr/test
       raw count for specified hdlr/test
       - $axis.ByTest(test)[<level>]
       - byTest[test][<level>]
*/
