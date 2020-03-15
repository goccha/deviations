package deviations

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type TestResult struct {
	Number int
	Name   string
	Point  int
}

var yamada = []TestResult{
	{
		Name:  "山田太郎",
		Point: 100,
	},
	{
		Name:  "山田次郎",
		Point: 50,
	},
	{
		Name:  "山田三郎",
		Point: 80,
	},
	{
		Name:  "山田四郎",
		Point: 30,
	},
	{
		Name:  "山田五郎",
		Point: 0,
	},
}
var sato = []TestResult{
	{
		Name:  "佐藤太郎",
		Point: 80,
	},
	{
		Name:  "佐藤次郎",
		Point: 50,
	},
	{
		Name:  "佐藤三郎",
		Point: 80,
	},
	{
		Name:  "佐藤四郎",
		Point: 30,
	},
	{
		Name:  "佐藤五郎",
		Point: 0,
	},
}
var suzuki = []TestResult{
	{
		Name:  "鈴木太郎",
		Point: 60,
	},
	{
		Name:  "鈴木次郎",
		Point: 90,
	},
	{
		Name:  "鈴木三郎",
		Point: 20,
	},
	{
		Name:  "鈴木四郎",
		Point: 20,
	},
	{
		Name:  "鈴木五郎",
		Point: 0,
	},
}

func preset() (*Calc, *Calc, *Calc) {
	y := New()
	for _, v := range yamada {
		y.AddInt(v.Point, v)
	}
	sa := New()
	for _, v := range sato {
		sa.AddInt(v.Point, v)
	}
	su := New()
	for _, v := range suzuki {
		su.AddInt(v.Point, v)
	}
	return y, sa, su
}

func TestCalc_Sort(t *testing.T) {
	y, sa, su := preset()
	ysasu := y.Union(sa).Union(su)
	ysasu.Sort(false).ForEach(func(elm Element) bool {
		fmt.Printf("%s\n", elm.String())
		return true
	}, false)
	fmt.Printf("== desc ==\n")
	ysasu.Sort(true).ForEach(func(elm Element) bool {
		fmt.Printf("%s\n", elm.String())
		return true
	}, false)
}

func TestCalc_ForEach(t *testing.T) {
	y, sa, su := preset()
	ysasu := y.Union(sa).Union(su)
	ysasu.Sort(false).ForEach(func(elm Element) bool {
		fmt.Printf("%s\n", elm.String())
		return true
	}, false)
	fmt.Printf("== desc ==\n")
	ysasu.ForEach(func(elm Element) bool {
		fmt.Printf("%s\n", elm.String())
		return true
	}, true)
}

func TestCalc_Union(t *testing.T) {
	y, sa, su := preset()
	avg := y.Avg()
	fmt.Printf("yamada.AVG=%f\n", avg)
	avg = sa.Avg()
	fmt.Printf("sato.AVG=%f\n", avg)
	avg = su.Avg()
	fmt.Printf("suzuki.AVG=%f\n", avg)
	ysa := y.Union(sa)
	avg = ysa.Avg()
	fmt.Printf("yamada&suzuki.AVG=%f\n", avg)
	ysasu := ysa.Union(su)
	avg = ysasu.Avg()
	fmt.Printf("yamada&suzuki&suzuki.AVG=%f\n", avg)
}

func TestCalc_Intersection(t *testing.T) {
	y, sa, su := preset()
	avg := y.Avg()
	fmt.Printf("yamada.AVG=%f\n", avg)
	avg = sa.Avg()
	fmt.Printf("sato.AVG=%f\n", avg)
	avg = su.Avg()
	fmt.Printf("suzuki.AVG=%f\n", avg)
	ysa := y.Intersection(sa)
	avg = ysa.Avg()
	fmt.Printf("yamada&suzuki.AVG=%f\n", avg)
	ysasu := ysa.Intersection(su)
	avg = ysasu.Avg()
	fmt.Printf("yamada&suzuki&suzuki.AVG=%f\n", avg)
}

func TestCalc_Difference(t *testing.T) {
	y, sa, su := preset()
	ys, ya := y.Difference(sa)
	if ys.Len() != 1 {
		t.Errorf("expected=1, actual=%d", ys.Len())
	}
	if ya.Len() != 0 {
		t.Errorf("expected=0, actual=%d", ya.Len())
	}
	sas, sau := sa.Difference(su)
	if sas.Len() != 4 {
		t.Errorf("expected=4, actual=%d", sas.Len())
	}
	if sau.Len() != 4 {
		t.Errorf("expected=4, actual=%d", sau.Len())
	}
}

func TestCalc_SymmetricDifference(t *testing.T) {
	y, sa, su := preset()
	ysa := y.SymmetricDifference(sa)
	if ysa.Len() != 1 {
		t.Errorf("expected=1, actual=%d", ysa.Len())
	}
	sau := sa.SymmetricDifference(su)
	if sau.Len() != 8 {
		t.Errorf("expected=8, actual=%d", sau.Len())
	}
}

func TestCalc_Search(t *testing.T) {
	y, sa, su := preset()
	ysasu := y.Union(sa).Union(su)
	elements := ysasu.Search(30)
	for _, elm := range elements {
		fmt.Printf("%v\n", elm)
	}
}

func TestCalc_Extract(t *testing.T) {
	y, sa, su := preset()
	ysasu := y.Union(sa).Union(su)
	cal := ysasu.Extract(func(elm Element) bool {
		if elm.Attached != nil {
			result := elm.Attached.(TestResult)
			if strings.HasSuffix(result.Name, "三郎") {
				return true
			}
		}
		return false
	})
	if cal.Len() != 3 {
		t.Errorf("expected=3,actual=%d", cal.Len())
	}
}

func TestRanking_ForEach(t *testing.T) {
	start := time.Now()
	c := New()
	max := 1000000
	for i := 0; i < max; i++ {
		c.AddInt(rand.Intn(16001))
	}
	fmt.Printf("エレメント数=%d\n", c.Len())
	fmt.Printf("最大値=%f\n", c.Max())
	fmt.Printf("最小値=%f\n", c.Min())
	fmt.Printf("標準偏差=%f\n", c.StandardDeviation())
	fmt.Printf("平均=%f\n", c.Avg())
	fmt.Printf("偏差値（8000）=%f\n", c.DeviationValue(8000))
	c.Ranking().ForEach(func(key float64, cnt, rank int) bool {
		//fmt.Printf("%f(%d) = %d (%f)\n", key, cnt, rank, c.DeviationValue(key))
		return true
	})
	fmt.Printf("ランキング数=%d\n", c.Ranking().Len())
	fmt.Printf("elapsed=%vs\n", time.Now().Sub(start).Seconds())
}

func TestRanking_Rank(t *testing.T) {
	points := []float64{
		0, 0, 0, 0, 50, 50, 100, 100, 100, 100,
	}
	c := New()
	for _, p := range points {
		c.Add(p)
	}
	c = c.Sum()
	expected := fmt.Sprintf("%f", 44.721360)
	actual := fmt.Sprintf("%f", c.StandardDeviation())
	if expected != actual {
		t.Errorf("expected=%s, actual=%s", expected, actual)
	}
	for _, p := range points {
		v := fmt.Sprintf("%f", c.DeviationValue(p))
		switch p {
		case 0:
			expected = fmt.Sprintf("%f", 38.81966)
		case 50:
			expected = fmt.Sprintf("%f", 50.0)
		case 100:
			expected = fmt.Sprintf("%f", 61.18034)
		}
		if v != expected {
			t.Errorf("expected=%s, actual=%s", expected, v)
		}
	}
	rank := c.Ranking().Rank(50)
	if rank != 5 {
		t.Errorf("expected=%d, actual=%d", 5, rank)
	}
}

func TestRanking_Value(t *testing.T) {
	start := time.Now()
	c := New()
	max := 1000000
	for i := 0; i < max; i++ {
		c.AddInt(rand.Intn(16001))
	}
	rank := c.Len() / 2
	val := c.Ranking().Value(rank)
	fmt.Printf("Rank=%d, 値=%f\n", rank, val)
	fmt.Printf("elapsed=%vs\n", time.Now().Sub(start).Seconds())
	fmt.Printf("finish.\n")
}

func TestRanking_Elements(t *testing.T) {
	points := []float64{
		0, 0, 0, 0, 50, 50, 100, 100, 100, 100,
	}
	c := New()
	for _, p := range points {
		c.Add(p)
	}
	elements := c.Ranking().Elements(6)
	if len(elements) != 2 {
		t.Errorf("expected=2,actual=%d", len(elements))
	}
}
