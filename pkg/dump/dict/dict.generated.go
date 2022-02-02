package dict

type FixDict interface {
	Version() string
	TagName(tag int) (string, bool)
	ValueName(tag int, value string) (string, bool)
}

type emptyDict struct {
	fixver string
}

func (d *emptyDict) Version() string {
	return d.fixver
}

func (d *emptyDict) TagName(tag int) (string, bool) {
	return "", false
}

func (d *emptyDict) ValueName(tag int, value string) (string, bool) {
	return "", false
}

var dicts map[string]FixDict

func init() {
	dicts = make(map[string]FixDict)
}

func Get(fixver string) FixDict {
	dict, ok := dicts[fixver]
	if ok {
		return dict
	}

	switch fixver {
	case "FIX.4.0":
		dict = &FIX40{}
		dicts[fixver] = dict
		return dict
	case "FIX.4.1":
		dict = &FIX41{}
		dicts[fixver] = dict
		return dict
	case "FIX.4.2":
		dict = &FIX42{}
		dicts[fixver] = dict
		return dict
	case "FIX.4.3":
		dict = &FIX43{}
		dicts[fixver] = dict
		return dict
	case "FIX.4.4":
		dict = &FIX44{}
		dicts[fixver] = dict
		return dict
	case "FIX.5.0":
		dict = &FIX50{}
		dicts[fixver] = dict
		return dict
	case "FIX.5.0SP1":
		dict = &FIX50SP1{}
		dicts[fixver] = dict
		return dict
	case "FIX.5.0SP2+CME":
		dict = &FIX50SP2CME{}
		dicts[fixver] = dict
		return dict
	case "FIX.5.0SP2":
		dict = &FIX50SP2{}
		dicts[fixver] = dict
		return dict
	case "FIX.5.0SP2.EP115":
		dict = &FIX50SP2EP115{}
		dicts[fixver] = dict
		return dict
	case "FIXT.1.1":
		dict = &FIXT11{}
		dicts[fixver] = dict
		return dict
	}

	dict = &emptyDict{fixver}
	dicts[fixver] = dict
	return dict
}
