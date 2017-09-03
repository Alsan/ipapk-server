package middleware

import (
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/phinexdaz/ipapk"
	"github.com/phinexdaz/ipapk-server/models"
	"github.com/phinexdaz/ipapk-server/serializers"
	"github.com/phinexdaz/ipapk-server/utils"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	ext := models.BundleFileExtension(filepath.Ext(file.Filename))
	if !ext.IsValid() {
		return
	}

	uuid := utils.NewUUID()
	filename := filepath.Join("data", "app", uuid+string(ext.PlatformType().Extention()))

	if err := c.SaveUploadedFile(file, filename); err != nil {
		return
	}

	app, err := ipapk.NewAppParser(filename)
	if err != nil {
		return
	}

	icon := uuid + ".png"
	if err := utils.SaveIcon(app.Icon, filepath.Join("data", "icon", icon)); err != nil {
		return
	}

	bundle := new(models.Bundle)
	bundle.UUID = uuid
	bundle.PlatformType = ext.PlatformType()
	bundle.Name = app.Name
	bundle.BundleId = app.BundleId
	bundle.Version = app.Version
	bundle.Build = app.Build
	bundle.Size = app.Size

	if err := models.AddBundle(bundle); err != nil {
		return
	}

	url := location.Get(c)

	c.JSON(http.StatusOK, &serializers.BundleJSON{
		UUID:       uuid,
		Name:       bundle.Name,
		Platform:   bundle.PlatformType.String(),
		BundleId:   bundle.BundleId,
		Version:    bundle.Version,
		Build:      bundle.Build,
		InstallUrl: bundle.GetInstallUrl(url.String()),
		QRCodeUrl:  url.String() + "/qrcode/" + uuid,
		IconUrl:    url.String() + "/icon/" + icon,
		Downloads:  bundle.Downloads,
	})
}

func QRCode(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	url := location.Get(c)

	data := fmt.Sprintf("%v/bundles/%v?_t=%v", url.String(), bundle.UUID, time.Now().Unix())
	img, err := qrcode.Encode(data, qrcode.Medium, -5)
	if err != nil {
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}

func GetBundle(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	url := location.Get(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"bundle":     bundle,
		"installUrl": bundle.GetInstallUrl(url.String()),
		"qrCodeUrl":  url.String() + "/qrcode/" + bundle.UUID,
		"iconUrl":    url.String() + "/icon/" + bundle.UUID + ".png",
	})
}

func GetVersions(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	versions, err := bundle.GetVersions()
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "version.html", gin.H{
		"versions": versions,
		"uuid":     bundle.UUID,
	})
}

func GetBuilds(c *gin.Context) {
	uuid := c.Param("uuid")
	version := c.Param("ver")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	builds, err := bundle.GetBuilds(version)
	if err != nil {
		return
	}

	url := location.Get(c)

	var bundles []serializers.BundleJSON
	for _, v := range builds {
		bundles = append(bundles, serializers.BundleJSON{
			UUID:       v.UUID,
			Name:       v.Name,
			Platform:   v.PlatformType.String(),
			BundleId:   v.BundleId,
			Version:    v.Version,
			Build:      v.Build,
			InstallUrl: v.GetInstallUrl(url.String()),
			QRCodeUrl:  url.String() + "/qrcode/" + uuid,
			IconUrl:    url.String() + "/icon/" + uuid + ".png",
			Downloads:  v.Downloads,
		})
	}

	c.HTML(http.StatusOK, "build.html", gin.H{
		"builds": bundles,
	})
}

func Plist(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	if bundle.PlatformType != models.BundlePlatformTypeIOS {
		return
	}

	url := location.Get(c)

	ipaUrl := url.String() + "/ipa/" + bundle.UUID

	data, err := models.NewPlist(bundle.Name, bundle.Version, bundle.BundleId, ipaUrl).Marshall()
	if err != nil {
		return
	}

	c.Data(http.StatusOK, "application/x-plist", data)
}

func DownloadIPA(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	if bundle.PlatformType != models.BundlePlatformTypeIOS {
		return
	}

	filename := bundle.UUID + string(bundle.PlatformType.Extention())
	file, err := ioutil.ReadFile(filepath.Join("data", "app", filename))
	if err != nil {
		return
	}

	bundle.UpdateDownload()

	c.Header("Content-Disposition", "attachment;filename="+filename)
	c.Header("Content-Length", strconv.Itoa(int(bundle.Size)))
	c.Data(http.StatusOK, "application/octet-stream", file)
}

func DownloadAPK(c *gin.Context) {
	uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUUID(uuid)
	if err != nil {
		return
	}

	if bundle.PlatformType != models.BundlePlatformTypeAndroid {
		return
	}

	filename := bundle.UUID + string(bundle.PlatformType.Extention())
	file, err := ioutil.ReadFile(filepath.Join("data", "app", filename))
	if err != nil {
		return
	}

	bundle.UpdateDownload()

	c.Header("Content-Disposition", "attachment;filename="+filename)
	c.Header("Content-Length", strconv.Itoa(int(bundle.Size)))
	c.Data(http.StatusOK, "application/vnd.android.package-archive", file)
}
