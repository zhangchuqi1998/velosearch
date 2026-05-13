package collection

import "github.com/zhangchuqi1998/velosearch/internal/distance"

type Metric int

const (
	MetricL2 Metric = iota
	MetricCosine
)

type Config struct {
	Name           string
	Dim            int
	Metric         Metric
	M              int
	EfConstruction int
}

func (m Metric) DistanceFunc() distance.DistanceFunc {
	switch m {
	case MetricL2:
		return distance.L2Squared
	case MetricCosine:
		return distance.Cosine
	}
	return distance.L2Squared
}
