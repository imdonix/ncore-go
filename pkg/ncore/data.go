package ncore

type ParamSort string

const (
	SortName           ParamSort = "name"
	SortUpload         ParamSort = "fid"
	SortSize           ParamSort = "size"
	SortTimesCompleted ParamSort = "times_completed"
	SortSeeders        ParamSort = "seeders"
	SortLeechers       ParamSort = "leechers"
)

type SearchParamType string

const (
	TypeSDHun       SearchParamType = "xvid_hun"
	TypeSD          SearchParamType = "xvid"
	TypeDVDHun      SearchParamType = "dvd_hun"
	TypeDVD         SearchParamType = "dvd"
	TypeDVD9Hun     SearchParamType = "dvd9_hun"
	TypeDVD9        SearchParamType = "dvd9"
	TypeHDHun       SearchParamType = "hd_hun"
	TypeHD          SearchParamType = "hd"
	TypeSDSerHun    SearchParamType = "xvidser_hun"
	TypeSDSer       SearchParamType = "xvidser"
	TypeDVDSerHun   SearchParamType = "dvdser_hun"
	TypeDVDSer      SearchParamType = "dvdser"
	TypeHDSerHun    SearchParamType = "hdser_hun"
	TypeHDSer       SearchParamType = "hdser"
	TypeMP3Hun      SearchParamType = "mp3_hun"
	TypeMP3         SearchParamType = "mp3"
	TypeLosslessHun SearchParamType = "lossless_hun"
	TypeLossless    SearchParamType = "lossless"
	TypeClip        SearchParamType = "clip"
	TypeGameIso     SearchParamType = "game_iso"
	TypeGameRip     SearchParamType = "game_rip"
	TypeConsole     SearchParamType = "console"
	TypeEbookHun    SearchParamType = "ebook_hun"
	TypeEbook       SearchParamType = "ebook"
	TypeIso         SearchParamType = "iso"
	TypeMisc        SearchParamType = "misc"
	TypeMobil       SearchParamType = "mobil"
	TypeXXXImg      SearchParamType = "xxx_imageset"
	TypeXXXSD       SearchParamType = "xxx_xvid"
	TypeXXXDVD      SearchParamType = "xxx_dvd"
	TypeXXXHD       SearchParamType = "xxx_hd"
	TypeAllOwn      SearchParamType = "all_own"
)

var detailedParamMap = map[string]SearchParamType{
	"osszes_film_xvid_hun":       TypeSDHun,
	"osszes_film_xvid":           TypeSD,
	"osszes_film_dvd_hun":        TypeDVDHun,
	"osszes_film_dvd":            TypeDVD,
	"osszes_film_dvd9_hun":       TypeDVD9Hun,
	"osszes_film_dvd9":           TypeDVD9,
	"osszes_film_hd_hun":         TypeHDHun,
	"osszes_film_hd":             TypeHD,
	"osszes_sorozat_xvidser_hun": TypeSDSerHun,
	"osszes_sorozat_xvidser":     TypeSDSer,
	"osszes_sorozat_dvdser_hun":  TypeDVDSerHun,
	"osszes_sorozat_dvdser":      TypeDVDSer,
	"osszes_sorozat_hdser_hun":   TypeHDSerHun,
	"osszes_sorozat_hdser":       TypeHDSer,
	"osszes_zene_mp3_hun":        TypeMP3Hun,
	"osszes_zene_mp3":            TypeMP3,
	"osszes_zene_lossless_hun":   TypeLosslessHun,
	"osszes_zene_lossless":       TypeLossless,
	"osszes_zene_clip":           TypeClip,
	"osszes_jatek_game_iso":      TypeGameIso,
	"osszes_jatek_game_rip":      TypeGameRip,
	"osszes_jatek_console":       TypeConsole,
	"osszes_konyv_ebook_hun":     TypeEbookHun,
	"osszes_konyv_ebook":         TypeEbook,
	"osszes_program_iso":         TypeIso,
	"osszes_program_misc":        TypeMisc,
	"osszes_program_mobil":       TypeMobil,
	"osszes_xxx_xxx_imageset":    TypeXXXImg,
	"osszes_xxx_xxx_xvid":        TypeXXXSD,
	"osszes_xxx_xxx_dvd":         TypeXXXDVD,
	"osszes_xxx_xxx_hd":          TypeXXXHD,
}

type SearchParamWhere string

const (
	WhereName        SearchParamWhere = "name"
	WhereDescription SearchParamWhere = "leiras"
	WhereIMDB        SearchParamWhere = "imdb"
	WhereLabel       SearchParamWhere = "cimke"
)

type ParamSeq string

const (
	SeqAsc  ParamSeq = "ASC"
	SeqDesc ParamSeq = "DESC"
)

const (
	URLIndex           = "https://ncore.pro/index.php"
	URLLogin           = "https://ncore.pro/login.php"
	URLActivity        = "https://ncore.pro/hitnrun.php"
	URLRecommended     = "https://ncore.pro/recommended.php"
	URLDownloadPattern = "https://ncore.pro/torrents.php?oldal=%d&tipus=%s&miszerint=%s&hogyan=%s&mire=%s&miben=%s"
	URLDetailPattern   = "https://ncore.pro/torrents.php?action=details&id=%s"
	URLDownloadLink    = "https://ncore.pro/torrents.php?action=download&id=%s&key=%s"
	URLCookieDomain    = "ncore.pro"
)

var AllowedCookies = []string{"nick", "pass", "stilus", "nyelv", "PHPSESSID"}
