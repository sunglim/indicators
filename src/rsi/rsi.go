package rsi

type record struct {
	rawData  float32
	upward   float32
	downward float32
	avgGain  float32
	avgLoss  float32
	rs       float32
	rsi      float32
}

type sheet struct {
	records []record
}

type RSICalculator struct {
	data  []float32
	Sheet sheet
}

func NewRSICalculator(data []float32) *RSICalculator {
	calculator := &RSICalculator{
		data:  data,
		Sheet: sheet{},
	}
	return calculator
}

func upward(before, after float32) float32 {
	if after > before {
		return (after - before) / before
	}
	return 0
}

func downward(before, after float32) float32 {
	if after < before {
		return (before - after) / before
	}
	return 0
}

func (r *RSICalculator) calculateAVGGainAndLoss() {
	// calculate Last 13 items's average
	for i := 0; i < len(r.data); i++ {
		// 0 index is not calculated.
		// 1 - 13
		if i < 14 {
			continue
		} else if i == 14 {
			totalUpward := float32(0)
			totalDownward := float32(0)
			for k := 1; k <= 14; k++ {
				totalUpward += float32(r.Sheet.records[k].upward)
				totalDownward += float32(r.Sheet.records[k].downward)
			}
			r.Sheet.records[i].avgGain = float32(totalUpward / 14)
			r.Sheet.records[i].avgLoss = float32(totalDownward / 14)

			r.Sheet.records[i].rs = r.Sheet.records[i].avgGain / r.Sheet.records[i].avgLoss
			r.Sheet.records[i].rsi = 100 - (100 / (r.Sheet.records[i].rs + 1))
		} else {
			previousAvgGain := r.Sheet.records[i-1].avgGain
			previousAvgGain = previousAvgGain*13 + r.Sheet.records[i].upward

			r.Sheet.records[i].avgGain = previousAvgGain / 14

			previousAvgLoss := r.Sheet.records[i-1].avgLoss
			previousAvgLoss = previousAvgLoss*13 + r.Sheet.records[i].downward

			r.Sheet.records[i].avgLoss = previousAvgLoss / 14

			r.Sheet.records[i].rs = r.Sheet.records[i].avgGain / r.Sheet.records[i].avgLoss
			r.Sheet.records[i].rsi = 100 - (100 / (r.Sheet.records[i].rs + 1))
		}
	}
}

// Fills out upward, downward.
// RawData | Upward | Downward
// ----------------------------
// 245.706665 | 0 | 0
// 244.919998 | 0 | 245.706665 - 244.919998
func (r *RSICalculator) calculateUpwardDownward() {
	for i := 0; i < len(r.data); i++ {
		if i == 0 {
			r.Sheet.records = append(r.Sheet.records, record{rawData: r.data[i]})
			continue
		}

		r.Sheet.records = append(r.Sheet.records, record{
			rawData:  r.data[i],
			upward:   upward(r.data[i-1], r.data[i]),
			downward: downward(r.data[i-1], r.data[i]),
		})
	}
}

func (r *RSICalculator) Load() {
	// We can optimize the calculation.
	// O(N)
	r.calculateUpwardDownward()
	// O(N)
	r.calculateAVGGainAndLoss()
}

func (r *RSICalculator) Calculate() {
}
