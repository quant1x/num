package stat

import num "gitee.com/quant1x/go-num/x64"

func (r RollingAndExpandingMixin) Count() (s Series) {
	if r.Series.Type() != SERIES_TYPE_BOOL {
		panic("不支持非bool序列")
	}
	values := make([]DType, r.Series.Len())
	for i, block := range r.GetBlocks() {
		if block.Len() == 0 {
			values[i] = 0
			continue
		}
		bs := block.Values().([]bool)
		values[i] = DType(num.Count(bs))
	}
	s = r.Series.Empty(SERIES_TYPE_DTYPE)
	s.Rename(r.Series.Name())
	s = s.Append(values)
	return
}
