package config

const (
	TEMP_DIR               = "temp"
	BUILD_DIR              = "build"
	REDIRECTS_FILE         = "_redirects"
	TEMPLATE_HTML          = "templates/index.html"
	GH_REPO_OWNER          = "sethcottle"
	GH_REPO_NAME           = "littlelink"
	DOWNLOAD_TAG_DEF_VER   = "v2.3.4"
	DOWNLOAD_ZIP_NAME      = "littlelink-%s.zip"
	BRANDS_CSS_FILE        = "%s/css/brands.css"
	IMAGES_ICONS           = "%s/images/icons/%s.svg"
	BUTTON_CLASS_NAME           = ".button.button-%s"
	BUTTON_DETAILS_DEF_BRAND    = "web"
	BUTTON_DETAILS_DEF_NAME     = "Google"
	BUTTON_DETAILS_DEF_ICON     = "generic-website"
	BUTTON_DETAILS_DEF_LINK     = "https://google.com"
)

var FILES_TO_DELETE = []string{
	"privacy.html",
	"images/littlelink.png",
	"images/littlelink.svg",
	"images/littlelink@2x.png",
	"LICENSE.md",
	"README.md",
	".gitignore",
}
