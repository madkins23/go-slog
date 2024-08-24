package warn

import (
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

type HandlerData struct {
	byTest   map[data.TestTag]uint
	subScore map[warning.Level]*score.Average

	//byTest        map[data.TestTag]*Average
	//scores        map[score.Type]score.Value
	//subScore      map[Weight]*Average
	//rollup        map[RollOver]*Average
}

/*
   Detail Table
     per hdlr
       sub score
         average of scores from all tests of hdlr
         show after test columns before 'regular' scores
         may be just another score type
     per hdlr/test
       score
         converted score from counts for all levels for hdlr/test
   Levels Table
     per hdlr
       sub score
         $axis.Score(score.SubScore)
         average of scores from all levels of hdlr
         show after test columns before 'regular' scores
         may be just another score type
     per hdlr/level
       score
         converted score from counts for specified level for all hdlr tests
         show after test columns before 'regular' scores
         may be just another score type
   Specific Level Table
     per hdlr
       sub score
         average of scores from all tests of handler for specified level
         show after test columns before 'regular' scores
         may be just another score type
     per hdlr/level
       raw count for specified level
         actual count for that level for hdlr/test
*/
