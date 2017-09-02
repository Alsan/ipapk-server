package serializers

type BundleJSON struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Platform   string `json:"platform"`
	BundleId   string `json:"bundleId"`
	Version    string `json:"version"`
	Build      string `json:"build"`
	InstallUrl string `json:"install_url"`
	QRCodeUrl  string `json:"qrcode_url"`
	IconUrl    string `json:"icon_url"`
	Downloads  uint64 `json:"downloads"`
}
