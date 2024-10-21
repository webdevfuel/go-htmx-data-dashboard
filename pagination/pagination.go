package pagination

type Page struct {
	Value    int64
	Disabled bool
}

func NextPage(page, totalPages int64) *Page {
	if page >= totalPages {
		return &Page{
			Value:    0,
			Disabled: true,
		}
	}
	return &Page{
		Value:    page + 1,
		Disabled: false,
	}
}

func PrevPage(page int64) *Page {
	if page == 1 {
		return &Page{
			Value:    0,
			Disabled: true,
		}
	}
	return &Page{
		Value:    page - 1,
		Disabled: false,
	}
}
