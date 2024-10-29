package deviceprof

type (
	Resolution    [2]uint
	DeviceProfile struct {
		Name       string
		Resolution Resolution
		Palette    PaletteType
		Scale      float64
	}
)

var defaultProfiles = map[DeviceType]DeviceProfile{
	DeviceKindle1:  {"Kindle 1", Resolution{600, 670}, Palette4, 1.8},
	DeviceKindle11: {"Kindle 11/12", Resolution{1072, 1448}, Palette16, 1.8},
	DeviceKindle2:  {"Kindle 2", Resolution{600, 670}, Palette15, 1.8},
	DeviceKindleKeyboardTouch: {
		"Kindle Keyboard/Touch", Resolution{600, 800}, Palette16, 1.8,
	},
	DeviceKindle:      {"Kindle", Resolution{600, 800}, Palette16, 1.8},
	DeviceKindleDXDXG: {"Kindle DX/DXG", Resolution{824, 1000}, Palette16, 1.8},
	DeviceKindlePaperwhite1_2: {
		"Kindle Paperwhite 1/2", Resolution{758, 1024}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite3_4_Voyage_Oasis: {
		"Kindle Paperwhite 3/4/Voyage/Oasis", Resolution{1072, 1448}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite5_SignatureEdition: {
		"Kindle Paperwhite 5/Signature Edition", Resolution{1236, 1648}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite6: {
		"Kindle Paperwhite 6", Resolution{1264, 1680}, Palette16, 1.8,
	},
	DeviceKindleColorsoft: {
		"Kindle Colorsoft", Resolution{1264, 1680}, Palette16, 1.8,
	},
	DeviceKindleOasis_2_3: {
		"Kindle Oasis 2/3", Resolution{1264, 1680}, Palette16, 1.8,
	},
	DeviceKindleScribe: {
		"Kindle Scribe", Resolution{1860, 2480}, Palette16, 1.8,
	},
	DeviceKoboMini_Touch: {
		"Kobo Mini/Touch", Resolution{600, 800}, Palette16, 1.8,
	},
	DeviceKoboGlo:    {"Kobo Glo", Resolution{768, 1024}, Palette16, 1.8},
	DeviceKoboGloHD:  {"Kobo Glo HD", Resolution{1072, 1448}, Palette16, 1.8},
	DeviceKoboAura:   {"Kobo Aura", Resolution{758, 1024}, Palette16, 1.8},
	DeviceKoboAuraHD: {"Kobo Aura HD", Resolution{1080, 1440}, Palette16, 1.8},
	DeviceKoboAuraH2O: {
		"Kobo Aura H2O", Resolution{1080, 1430}, Palette16, 1.8,
	},
	DeviceKoboAuraONE: {
		"Kobo Aura ONE", Resolution{1404, 1872}, Palette16, 1.8,
	},
	DeviceKoboNia: {"Kobo Nia", Resolution{758, 1024}, Palette16, 1.8},
	DeviceKoboClaraHD_KoboClara2E: {
		"Kobo Clara HD/Kobo Clara 2E", Resolution{1072, 1448}, Palette16, 1.8,
	},
	DeviceKoboClaraColour: {
		"Kobo Clara Colour", Resolution{1072, 1448}, Palette16, 1.8,
	},
	DeviceKoboLibraH2O_KoboLibra2: {
		"Kobo Libra H2O/Kobo Libra 2", Resolution{1264, 1680}, Palette16, 1.8,
	},
	DeviceKoboLibraColour: {
		"Kobo Libra Colour", Resolution{1264, 1680}, Palette16, 1.8,
	},
	DeviceKoboForma:  {"Kobo Forma", Resolution{1440, 1920}, Palette16, 1.8},
	DeviceKoboSage:   {"Kobo Sage", Resolution{1440, 1920}, Palette16, 1.8},
	DeviceKoboElipsa: {"Kobo Elipsa", Resolution{1404, 1872}, Palette16, 1.8},
	DeviceOther:      {"Other", Resolution{0, 0}, Palette16, 1.8},
}

func Profile(name DeviceType) (DeviceProfile, bool) {
	prof, found := defaultProfiles[name]
	return prof, found
}
