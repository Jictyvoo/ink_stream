package deviceprof

import (
	"strings"

	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

type (
	Resolution    = inktypes.ImageDimensions
	DeviceProfile struct {
		Name       string
		Resolution Resolution
		Palette    PaletteType
		Scale      float64
	}
)

var defaultProfiles = map[DeviceType]DeviceProfile{
	DeviceKindle1:  {"Kindle 1", Resolution{Width: 600, Height: 670}, Palette4, 1.8},
	DeviceKindle11: {"Kindle 11/12", Resolution{Width: 1072, Height: 1448}, Palette16, 1.8},
	DeviceKindle2:  {"Kindle 2", Resolution{Width: 600, Height: 670}, Palette15, 1.8},
	DeviceKindleKeyboardTouch: {
		"Kindle Keyboard/Touch", Resolution{Width: 600, Height: 800}, Palette16, 1.8,
	},
	DeviceKindle:      {"Kindle", Resolution{Width: 600, Height: 800}, Palette16, 1.8},
	DeviceKindleDXDXG: {"Kindle DX/DXG", Resolution{Width: 824, Height: 1000}, Palette16, 1.8},
	DeviceKindlePaperwhite1_2: {
		"Kindle Paperwhite 1/2", Resolution{Width: 758, Height: 1024}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite3_4_Voyage_Oasis: {
		"Kindle Paperwhite 3/4/Voyage/Oasis", Resolution{Width: 1072, Height: 1448}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite5_SignatureEdition: {
		"Kindle Paperwhite 5/Signature Edition", Resolution{Width: 1236, Height: 1648}, Palette16, 1.8,
	},
	DeviceKindlePaperwhite6: {
		"Kindle Paperwhite 6", Resolution{Width: 1264, Height: 1680}, Palette16, 1.8,
	},
	DeviceKindleColorsoft: {
		"Kindle Colorsoft", Resolution{Width: 1264, Height: 1680}, Palette16, 1.8,
	},
	DeviceKindleOasis_2_3: {
		"Kindle Oasis 2/3", Resolution{Width: 1264, Height: 1680}, Palette16, 1.8,
	},
	DeviceKindleScribe: {
		"Kindle Scribe", Resolution{Width: 1860, Height: 2480}, Palette16, 1.8,
	},
	DeviceKoboMini_Touch: {
		"Kobo Mini/Touch", Resolution{Width: 600, Height: 800}, Palette16, 1.8,
	},
	DeviceKoboGlo:    {"Kobo Glo", Resolution{Width: 768, Height: 1024}, Palette16, 1.8},
	DeviceKoboGloHD:  {"Kobo Glo HD", Resolution{Width: 1072, Height: 1448}, Palette16, 1.8},
	DeviceKoboAura:   {"Kobo Aura", Resolution{Width: 758, Height: 1024}, Palette16, 1.8},
	DeviceKoboAuraHD: {"Kobo Aura HD", Resolution{Width: 1080, Height: 1440}, Palette16, 1.8},
	DeviceKoboAuraH2O: {
		"Kobo Aura H2O", Resolution{Width: 1080, Height: 1430}, Palette16, 1.8,
	},
	DeviceKoboAuraONE: {
		"Kobo Aura ONE", Resolution{Width: 1404, Height: 1872}, Palette16, 1.8,
	},
	DeviceKoboNia: {"Kobo Nia", Resolution{Width: 758, Height: 1024}, Palette16, 1.8},
	DeviceKoboClaraHD_KoboClara2E: {
		"Kobo Clara HD/Kobo Clara 2E", Resolution{Width: 1072, Height: 1448}, Palette16, 1.8,
	},
	DeviceKoboClaraColour: {
		"Kobo Clara Colour", Resolution{Width: 1072, Height: 1448}, Palette16, 1.8,
	},
	DeviceKoboLibraH2O_KoboLibra2: {
		"Kobo Libra H2O/Kobo Libra 2", Resolution{Width: 1264, Height: 1680}, Palette16, 1.8,
	},
	DeviceKoboLibraColour: {
		"Kobo Libra Colour", Resolution{Width: 1264, Height: 1680}, Palette16, 1.8,
	},
	DeviceKoboForma:  {"Kobo Forma", Resolution{Width: 1440, Height: 1920}, Palette16, 1.8},
	DeviceKoboSage:   {"Kobo Sage", Resolution{Width: 1440, Height: 1920}, Palette16, 1.8},
	DeviceKoboElipsa: {"Kobo Elipsa", Resolution{Width: 1404, Height: 1872}, Palette16, 1.8},
	DeviceOther:      {"Other", Resolution{}, Palette16, 1.8},
}

func Profile(name DeviceType) (DeviceProfile, bool) {
	prof, found := defaultProfiles[name]
	if !found {
		for key, dProf := range defaultProfiles {
			if strings.EqualFold(string(key), string(name)) {
				return dProf, true
			}
		}
	}
	return prof, found
}
