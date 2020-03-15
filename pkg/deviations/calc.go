package deviations

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Round 四捨五入
func Round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

// New 計算機を生成する
func New() *Calc {
	c := &Calc{
		elements: make([]Element, 0, 10),
		min:      math.NaN(),
	}
	c.ranking = newRanking(c)
	return c
}

// Calc　集合計算機
type Calc struct {
	elements              []Element
	total                 float64  //　合計
	avg                   float64  // 平均
	dispersion            float64  // 分散
	totalSquaredDeviation float64  // 偏差の二乗の合計
	standardDeviation     float64  // 標準偏差
	ranking               *Ranking // ランキング
	max                   float64
	min                   float64
	summed                bool
}

// Len 集合に含まれるElementの数を返す
func (c *Calc) Len() int {
	return len(c.elements)
}

// Sort ソート
func (c *Calc) Sort(desc bool) *Calc {
	if desc {
		sort.Slice(c.elements, func(i, j int) bool {
			return c.elements[i].value < c.elements[j].value
		})
	} else {
		sort.Slice(c.elements, func(i, j int) bool {
			return c.elements[i].value > c.elements[j].value
		})
	}
	return c
}

// CustomSort 任意ソート
func (c *Calc) CustomSort(f func(i, j int) bool) *Calc {
	sort.Slice(c.elements, f)
	return c
}

// Search リストから指定の点数のElementリストを検索
func (c *Calc) Search(value float64) []Element {
	c.Sort(false)
	index := sort.Search(len(c.elements), func(i int) bool {
		return c.elements[i].value <= value
	})
	if index < len(c.elements) {
		key := floatingKey{key: value}
		if cnt, ok := c.ranking.values[key]; ok {
			last := index + cnt
			return c.elements[index:last]
		}
	}
	return []Element{}
}

// AddInt 整数を追加する
func (c *Calc) AddInt(v int, a ...interface{}) *Calc {
	return c.Add(float64(v), a...)
}

// add Elementを追加する
func (c *Calc) add(elm Element) *Calc {
	v := elm.value
	c.elements = append(c.elements, elm)
	c.total += v
	c.ranking.Add(v)
	if v > c.max {
		c.max = v
	}
	if c.min == math.NaN() || v < c.min {
		c.min = v
	}
	c.summed = false
	return c
}

// Add 数値を追加する
func (c *Calc) Add(v float64, a ...interface{}) *Calc {
	var elm Element
	if a != nil {
		if len(a) == 1 {
			elm = Element{value: v, Attached: a[0]}
		} else {
			elm = Element{value: v, Attached: a}
		}
	} else {
		elm = Element{value: v}
	}
	return c.add(elm)
}

// Clone 集合のクローンを生成する
func (c *Calc) Clone() *Calc {
	elms := make([]Element, len(c.elements), len(c.elements))
	copy(elms, c.elements)
	return &Calc{
		elements: elms,
		total:    c.total,
		ranking:  c.ranking.Clone(),
		max:      c.max,
		min:      c.min,
	}
}

// Union 和集合
func (c *Calc) Union(s *Calc) *Calc {
	calc := New()
	for _, elm := range c.elements {
		calc.add(elm)
	}
	for _, elm := range s.elements {
		calc.add(elm)
	}
	return calc
}

// Intersection　積集合
func (c *Calc) Intersection(s *Calc) *Calc {
	calc := New()
	for _, elm := range c.elements {
		key := floatingKey{key: elm.value}
		if _, ok := s.ranking.values[key]; ok {
			calc.add(elm)
		}
	}
	for _, elm := range s.elements {
		key := floatingKey{key: elm.value}
		if _, ok := c.ranking.values[key]; ok {
			calc.add(elm)
		}
	}
	return calc
}

// Difference 差集合
func (c *Calc) Difference(s *Calc) (*Calc, *Calc) {
	s1 := New()
	for _, elm := range c.elements {
		key := floatingKey{key: elm.value}
		if _, ok := s.ranking.values[key]; !ok {
			s1.add(elm)
		}
	}
	s2 := New()
	for _, elm := range s.elements {
		key := floatingKey{key: elm.value}
		if _, ok := c.ranking.values[key]; !ok {
			s2.add(elm)
		}
	}
	return s1, s2
}

// SymmetricDifference 対称差集合
func (c *Calc) SymmetricDifference(s *Calc) *Calc {
	s1 := New()
	for _, elm := range c.elements {
		key := floatingKey{key: elm.value}
		if _, ok := s.ranking.values[key]; !ok {
			s1.add(elm)
		}
	}
	for _, elm := range s.elements {
		key := floatingKey{key: elm.value}
		if _, ok := c.ranking.values[key]; !ok {
			s1.add(elm)
		}
	}
	return s1
}

// Min 最小値を返す
func (c *Calc) Min() float64 {
	return c.min
}

// Max 最大値を返す
func (c *Calc) Max() float64 {
	return c.max
}

// Total　合計値を返す
func (c *Calc) Total() float64 {
	return c.total
}

// Avg　平均値を返す
func (c *Calc) Avg() float64 {
	c.Sum()
	return c.avg
}

// Dispersion 分散を返す
func (c *Calc) Dispersion() float64 {
	c.Sum()
	return c.dispersion
}

// TotalSquaredDeviation 二乗偏差の合計値を返す
func (c *Calc) TotalSquaredDeviation() float64 {
	c.Sum()
	return c.totalSquaredDeviation
}

// Deviation 指定した値の偏差を返す
func (c *Calc) Deviation(v float64) float64 {
	c.Sum()
	elm := Element{
		value: v,
	}
	elm.SquaredDeviation(c.Avg())
	return elm.deviation
}

// StandardDeviation 標準偏差を返す
func (c *Calc) StandardDeviation() float64 {
	c.Sum()
	return c.standardDeviation
}

// Sum 集計処理
func (c *Calc) Sum() *Calc {
	if !c.summed {
		c.avg = c.total / float64(len(c.elements))
		for _, elm := range c.elements {
			c.totalSquaredDeviation += elm.SquaredDeviation(c.avg)
		}
		c.dispersion = c.totalSquaredDeviation / float64(len(c.elements))
		c.standardDeviation = math.Sqrt(c.dispersion)
		c.summed = true
	}
	return c
}

// Extract 指定したフィルターで選択したElementを保持するCalcを返す
func (c *Calc) Extract(f func(elm Element) bool) *Calc {
	cal := New()
	for _, e := range c.elements {
		if f(e) {
			cal.add(e)
		}
	}
	return cal
}

// ForEach 集合に含まれるエレメントループ
func (c *Calc) ForEach(f func(elm Element) bool, reverse bool) {
	last := len(c.elements) - 1
	for i, elm := range c.elements {
		if reverse {
			elm = c.elements[last-i]
		}
		f(elm)
	}
}

// Contains 集合の中に指定した値が含まれるか判定する
func (c *Calc) Contains(v float64) bool {
	key := floatingKey{key: v}
	if _, ok := c.ranking.values[key]; ok {
		return true
	}
	return false
}

// Ranking ランキングを返す
func (c *Calc) Ranking() *Ranking {
	return c.ranking
}

// DeviationValue 指定した値の偏差値を返す
func (c *Calc) DeviationValue(v float64) float64 {
	c.Sum()
	elm := Element{
		value: v,
	}
	elm.SquaredDeviation(c.Avg())
	return elm.deviationValue(c.StandardDeviation())
}

// Element 集計要素
type Element struct {
	value     float64 // 値
	deviation float64 // 偏差
	Attached  interface{}
}

// String 文字列表現を返す
func (elm *Element) String() string {
	buf := strings.Builder{}
	buf.WriteString("value=")
	buf.WriteString(strconv.FormatFloat(elm.value, 'f', 0, 64))
	buf.WriteString(",deviation=")
	buf.WriteString(strconv.FormatFloat(elm.deviation, 'f', 0, 64))
	buf.WriteString(",Attached=")
	buf.WriteString(fmt.Sprintf("%v", elm.Attached))
	return buf.String()
}

// Value 値を返す
func (elm *Element) Value() float64 {
	return elm.value
}

// Deviation 偏差を返す
func (elm *Element) Deviation() float64 {
	return elm.deviation
}

// SquaredDeviation　二乗偏差を返す
func (elm *Element) SquaredDeviation(avg float64) float64 {
	elm.deviation = elm.value - avg
	return elm.deviation * elm.deviation
}

// DeviationValue 偏差値を返す
func (elm *Element) deviationValue(std float64) float64 {
	return elm.deviation/std*10 + 50 // （得点－平均点）÷標準偏差×10＋50
}

// floatingKey ランキングキー
type floatingKey struct {
	key float64
}

// Ranking ランキング
type Ranking struct {
	c      *Calc
	keys   []floatingKey
	values map[floatingKey]int
	total  int
	sorted bool
}

// newRanking ランキングオブジェクトを生成する
func newRanking(c *Calc) *Ranking {
	return &Ranking{
		c:      c,
		keys:   make([]floatingKey, 0, 10),
		values: make(map[floatingKey]int),
	}
}

// Clone Rankingのクローンを生成する
func (r *Ranking) Clone() *Ranking {
	dst := make([]floatingKey, len(r.keys), len(r.keys))
	copy(dst, r.keys)
	m := make(map[floatingKey]int)
	for k, v := range r.values {
		m[k] = v
	}
	return &Ranking{
		keys:   dst,
		values: m,
		sorted: false,
	}
}

// Add ランキングに登録する
func (r *Ranking) Add(val float64) {
	key := floatingKey{key: val}
	if v, ok := r.values[key]; ok {
		r.values[key] = v + 1
	} else {
		r.values[key] = 1
		r.keys = append(r.keys, floatingKey{key: val})
	}
	r.total++
	r.sorted = false
}

// Len ランキングの数を返す
func (r *Ranking) Len() int {
	return len(r.values)
}

// sort キーをソートする
func (r *Ranking) sort() {
	if !r.sorted {
		sort.Slice(r.keys, func(i, j int) bool {
			return r.keys[i].key > r.keys[j].key
		})
		r.sorted = true
	}
}

// Rank 指定の値の現在ランクを返す
// ランキングに含まれない値の場合、0を返す
func (r *Ranking) Rank(val float64) int {
	r.sort()
	key := floatingKey{key: val}
	if _, ok := r.values[key]; ok {
		rank := 1
		for _, k := range r.keys {
			if k == key {
				break
			}
			rank += r.values[k]
		}
		return rank
	}
	return 0
}

// Value 指定した順位の値を返す
// 存在しない順位の場合、NaNを返す
func (r *Ranking) Value(rank int) float64 {
	if rank > 0 {
		r.sort()
		cnt := 1
		var prev *floatingKey
		if r.total >= rank {
			for _, k := range r.keys {
				cnt += r.values[k]
				if cnt == rank {
					return k.key
				} else if cnt > rank {
					if prev != nil {
						return prev.key
					}
					return k.key
				}
				prev = &k
			}
		}
	}
	return math.NaN()
}

// Elements 指定した順位の値と等しいElementのリストを返す
func (r *Ranking) Elements(rank int) []Element {
	val := r.Value(rank)
	if val != math.NaN() {
		return r.c.Search(val)
	}
	return []Element{}
}

// ForEach ランキングリストの順次処理
func (r *Ranking) ForEach(f func(key float64, cnt, rank int) bool) {
	r.sort()
	rank := 1
	for _, k := range r.keys {
		if !f(k.key, r.values[k], rank) {
			break
		}
		rank += r.values[k]
	}
}
