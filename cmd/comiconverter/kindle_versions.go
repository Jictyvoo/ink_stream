package main

type PaletteType int

const (
	Palette4 PaletteType = iota
	Palette15
	Palette16
)

type DeviceProfile struct {
	Name       string
	Resolution [2]int
	Palette    PaletteType
	Scale      float64
}

type DeviceType string

//goland:noinspection GoSnakeCaseUsage
const (
	KDeviceKindle1                            DeviceType = "K1"
	KDeviceKindle11                           DeviceType = "K11"
	KDeviceKindle2                            DeviceType = "K2"
	KDeviceKindleKeyboardTouch                DeviceType = "K34"
	KDeviceKindle                             DeviceType = "K578"
	KDeviceKindleDXDXG                        DeviceType = "KDX"
	KDeviceKindlePaperwhite1_2                DeviceType = "KPW"
	KDeviceKindlePaperwhite3_4_Voyage_Oasis   DeviceType = "KV"
	KDeviceKindlePaperwhite5_SignatureEdition DeviceType = "KPW5"
	KDeviceKindleOasis_2_3                    DeviceType = "KO"
	KDeviceKindleScribe                       DeviceType = "KS"
	KDeviceKoboMini_Touch                     DeviceType = "KoMT"
	KDeviceKoboGlo                            DeviceType = "KoG"
	KDeviceKoboGloHD                          DeviceType = "KoGHD"
	KDeviceKoboAura                           DeviceType = "KoA"
	KDeviceKoboAuraHD                         DeviceType = "KoAHD"
	KDeviceKoboAuraH2O                        DeviceType = "KoAH2O"
	KDeviceKoboAuraONE                        DeviceType = "KoAO"
	KDeviceKoboNia                            DeviceType = "KoN"
	KDeviceKoboClaraHD_KoboClara2E            DeviceType = "KoC"
	KDeviceKoboClaraColour                    DeviceType = "KoCC"
	KDeviceKoboLibraH2O_KoboLibra2            DeviceType = "KoL"
	KDeviceKoboLibraColour                    DeviceType = "KoLC"
	KDeviceKoboForma                          DeviceType = "KoF"
	KDeviceKoboSage                           DeviceType = "KoS"
	KDeviceKoboElipsa                         DeviceType = "KoE"
	KDeviceOther                              DeviceType = "OTHER"
)

var DeviceProfiles = map[DeviceType]DeviceProfile{
	KDeviceKindle1:  {"Kindle 1", [2]int{600, 670}, Palette4, 1.8},
	KDeviceKindle11: {"Kindle 11", [2]int{1072, 1448}, Palette16, 1.8},
	KDeviceKindle2:  {"Kindle 2", [2]int{600, 670}, Palette15, 1.8},
	KDeviceKindleKeyboardTouch: {
		"Kindle Keyboard/Touch",
		[2]int{600, 800},
		Palette16,
		1.8,
	},
	KDeviceKindle:      {"Kindle", [2]int{600, 800}, Palette16, 1.8},
	KDeviceKindleDXDXG: {"Kindle DX/DXG", [2]int{824, 1000}, Palette16, 1.8},
	KDeviceKindlePaperwhite1_2: {
		"Kindle Paperwhite 1/2",
		[2]int{758, 1024},
		Palette16,
		1.8,
	},
	KDeviceKindlePaperwhite3_4_Voyage_Oasis: {
		"Kindle Paperwhite 3/4/Voyage/Oasis",
		[2]int{1072, 1448},
		Palette16,
		1.8,
	},
	KDeviceKindlePaperwhite5_SignatureEdition: {
		"Kindle Paperwhite 5/Signature Edition",
		[2]int{1236, 1648},
		Palette16,
		1.8,
	},
	KDeviceKindleOasis_2_3: {
		"Kindle Oasis 2/3",
		[2]int{1264, 1680},
		Palette16,
		1.8,
	},
	KDeviceKindleScribe: {
		"Kindle Scribe",
		[2]int{1860, 2480},
		Palette16,
		1.8,
	},
	KDeviceKoboMini_Touch: {
		"Kobo Mini/Touch",
		[2]int{600, 800},
		Palette16,
		1.8,
	},
	KDeviceKoboGlo:    {"Kobo Glo", [2]int{768, 1024}, Palette16, 1.8},
	KDeviceKoboGloHD:  {"Kobo Glo HD", [2]int{1072, 1448}, Palette16, 1.8},
	KDeviceKoboAura:   {"Kobo Aura", [2]int{758, 1024}, Palette16, 1.8},
	KDeviceKoboAuraHD: {"Kobo Aura HD", [2]int{1080, 1440}, Palette16, 1.8},
	KDeviceKoboAuraH2O: {
		"Kobo Aura H2O",
		[2]int{1080, 1430},
		Palette16,
		1.8,
	},
	KDeviceKoboAuraONE: {
		"Kobo Aura ONE",
		[2]int{1404, 1872},
		Palette16,
		1.8,
	},
	KDeviceKoboNia: {"Kobo Nia", [2]int{758, 1024}, Palette16, 1.8},
	KDeviceKoboClaraHD_KoboClara2E: {
		"Kobo Clara HD/Kobo Clara 2E",
		[2]int{1072, 1448},
		Palette16,
		1.8,
	},
	KDeviceKoboClaraColour: {
		"Kobo Clara Colour",
		[2]int{1072, 1448},
		Palette16,
		1.8,
	},
	KDeviceKoboLibraH2O_KoboLibra2: {
		"Kobo Libra H2O/Kobo Libra 2",
		[2]int{1264, 1680},
		Palette16,
		1.8,
	},
	KDeviceKoboLibraColour: {
		"Kobo Libra Colour",
		[2]int{1264, 1680},
		Palette16,
		1.8,
	},
	KDeviceKoboForma:  {"Kobo Forma", [2]int{1440, 1920}, Palette16, 1.8},
	KDeviceKoboSage:   {"Kobo Sage", [2]int{1440, 1920}, Palette16, 1.8},
	KDeviceKoboElipsa: {"Kobo Elipsa", [2]int{1404, 1872}, Palette16, 1.8},
	KDeviceOther:      {"Other", [2]int{0, 0}, Palette16, 1.8},
}
