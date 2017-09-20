package middleware

import (
	"bytes"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/phinexdaz/ipapk"
	"github.com/phinexdaz/ipapk-server/conf"
	"github.com/phinexdaz/ipapk-server/models"
	"github.com/phinexdaz/ipapk-server/serializers"
	"github.com/satori/go.uuid"
	"image/png"
	"net/http"
	"path/filepath"
	"time"
)

func Upload(c *gin.Context) {
	changelog := c.PostForm("changelog")
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	ext := models.BundleFileExtension(filepath.Ext(file.Filename))
	if !ext.IsValid() {
		return
	}

	_uuid := uuid.NewV4().String()
	filename := filepath.Join(".data", _uuid+string(ext.PlatformType().Extention()))

	if err := c.SaveUploadedFile(file, filename); err != nil {
		return
	}

	app, err := ipapk.NewAppParser(filename)
	if err != nil {
		return
	}

	iconBuffer := new(bytes.Buffer)
	if err := png.Encode(iconBuffer, app.Icon); err != nil {
		return
	}

	bundle := new(models.Bundle)
	bundle.UUID = _uuid
	bundle.PlatformType = ext.PlatformType()
	bundle.Name = app.Name
	bundle.BundleId = app.BundleId
	bundle.Version = app.Version
	bundle.Build = app.Build
	bundle.Size = app.Size
	bundle.Icon = iconBuffer.Bytes()
	bundle.ChangeLog = changelog

	if err := models.AddBundle(bundle); err != nil {
		return
	}

	c.JSON(http.StatusOK, &serializers.BundleJSON{
		UUID:       _uuid,
		Name:       bundle.Name,
		Platform:   bundle.PlatformType.String(),
		BundleId:   bundle.BundleId,
		Version:    bundle.Version,
		Build:      bundle.Build,
		InstallUrl: bundle.GetInstallUrl(conf.AppConfig.ProxyURL()),
		QRCodeUrl:  conf.AppConfig.ProxyURL() + "/qrcode/" + _uuid,
		IconUrl:    conf.AppConfig.ProxyURL() + "/icon/" + _uuid,
		Changelog:  bundle.ChangeLog,
		Downloads:  bundle.Downloads,
	})
}

func GetQRCode(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	data := fmt.Sprintf("%v/bundle/%v?_t=%v", conf.AppConfig.ProxyURL(), bundle.UUID, time.Now().Unix())
	code, err := qr.Encode(data, qr.L, qr.Unicode)
	if err != nil {
		return
	}
	code, err = barcode.Scale(code, 160, 160)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, code); err != nil {
		return
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

func GetIcon(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	c.Data(http.StatusOK, "image/png", bundle.Icon)
}

func GetChangelog(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "change.html", gin.H{
		"changelog": bundle.ChangeLog,
	})
}

func GetBundle(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"bundle":     bundle,
		"installUrl": bundle.GetInstallUrl(conf.AppConfig.ProxyURL()),
		"qrCodeUrl":  conf.AppConfig.ProxyURL() + "/qrcode/" + bundle.UUID,
		"iconUrl":    conf.AppConfig.ProxyURL() + "/icon/" + bundle.UUID,
	})
}

func GetVersions(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
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
	_uuid := c.Param("uuid")
	version := c.Param("ver")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	builds, err := bundle.GetBuilds(version)
	if err != nil {
		return
	}

	var bundles []serializers.BundleJSON
	for _, v := range builds {
		bundles = append(bundles, serializers.BundleJSON{
			UUID:       v.UUID,
			Name:       v.Name,
			Platform:   v.PlatformType.String(),
			BundleId:   v.BundleId,
			Version:    v.Version,
			Build:      v.Build,
			InstallUrl: v.GetInstallUrl(conf.AppConfig.ProxyURL()),
			QRCodeUrl:  conf.AppConfig.ProxyURL() + "/qrcode/" + v.UUID,
			IconUrl:    conf.AppConfig.ProxyURL() + "/icon/" + v.UUID,
			Changelog:  bundle.ChangeLog,
			Downloads:  v.Downloads,
		})
	}

	c.HTML(http.StatusOK, "build.html", gin.H{
		"builds": bundles,
	})
}

func GetPlist(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	if bundle.PlatformType != models.BundlePlatformTypeIOS {
		return
	}

	ipaUrl := conf.AppConfig.ProxyURL() + "/ipa/" + bundle.UUID

	data, err := models.NewPlist(bundle.Name, bundle.Version, bundle.BundleId, ipaUrl).Marshall()
	if err != nil {
		return
	}

	c.Data(http.StatusOK, "application/x-plist", data)
}

func DownloadAPP(c *gin.Context) {
	_uuid := c.Param("uuid")

	bundle, err := models.GetBundleByUID(_uuid)
	if err != nil {
		return
	}

	bundle.UpdateDownload()

	downloadUrl := conf.AppConfig.ProxyURL() + "/app/" + bundle.UUID + string(bundle.PlatformType.Extention())
	c.Redirect(http.StatusMovedPermanently, downloadUrl)
}
