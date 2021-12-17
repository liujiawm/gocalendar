package gocalendar

// (*SolarTermItem) clone
func (sti *SolarTermItem) clone() *SolarTermItem {
	if sti == nil {
		return nil
	}

	t := sti.Time.AddDate(0, 0, 0)
	return &SolarTermItem{
		Index: sti.Index,
		Name:  sti.Name,
		Time:  &t,
	}
}

// (*StarSignItem) clone
func (ssi *StarSignItem) clone() *StarSignItem {
	if ssi == nil {
		return nil
	}

	return &StarSignItem{
		Index: ssi.Index,
		Name:  ssi.Name,
	}
}

// (*FestivalItem) clone
func (fi *FestivalItem) clone() *FestivalItem {
	if fi == nil {
		return nil
	}

	var shows []string
	var secondary []string

	shows = append(shows, fi.Show...)
	secondary = append(secondary, fi.Secondary...)

	return &FestivalItem{
		Show:      shows,
		Secondary: secondary,
	}
}

// (*GZItem) clone
func (gzi *GZItem) clone() *GZItem {
	if gzi == nil {
		return nil
	}

	return &GZItem{
		HSI: gzi.HSI,
		HSN: gzi.HSN,
		EBI: gzi.EBI,
		EBN: gzi.EBN,
	}
}

// (*GZ)clone
func (gz *GZ) clone() *GZ {
	if gz == nil {
		return nil
	}

	return &GZ{
		Year:  gz.Year.clone(),
		Month: gz.Month.clone(),
		Day:   gz.Day.clone(),
		Hour:  gz.Hour.clone(),
	}
}

// (*LunarDate) clone
func (ld *LunarDate) clone() *LunarDate {
	if ld == nil {
		return nil
	}

	return &LunarDate{
		Year:          ld.Year,
		Month:         ld.Month,
		Day:           ld.Day,
		MonthName:     ld.MonthName,
		DayName:       ld.DayName,
		LeapStr:       ld.LeapStr,
		YearLeapMonth: ld.YearLeapMonth,
		AnimalIndex:   ld.AnimalIndex,
		AnimalName:    ld.AnimalName,
		YearGZ:        ld.YearGZ.clone(),
		Festival:      ld.Festival.clone(),
	}
}

// (*CalendarItem) clone
func (ci *CalendarItem) clone() *CalendarItem {
	if ci == nil {
		return nil
	}

	t := ci.Time.AddDate(0, 0, 0)
	return &CalendarItem{
		Time:         &t,
		IsAccidental: ci.IsAccidental,
		IsToday:      ci.IsToday,
		Festival:     ci.Festival.clone(),
		SolarTerm:    ci.SolarTerm.clone(),
		GZ:           ci.GZ.clone(),
		LunarDate:    ci.LunarDate.clone(),
		StarSign:     ci.StarSign.clone(),
	}
}

// (*Calendar) clone 克隆一个Calendar
// 该克隆不对临时数据克隆,而是清空临时数据
func (c *Calendar) Clone() *Calendar {
	if c == nil {
		return nil
	}

	rawT := c.rawTime.AddDate(0, 0, 0)

	var items []*CalendarItem
	if c.Items != nil {
		for _, itemv := range c.Items {
			items = append(items, itemv.clone())
		}
	}

	return &Calendar{
		Items:    items,
		config:   c.config.clone(),
		loc:      c.loc,
		rawTime:  &rawT,
		tempData: newCalendarTempData(),
	}
}
