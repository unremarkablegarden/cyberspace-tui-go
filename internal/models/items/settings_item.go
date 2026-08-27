package items

type SettingsSection uint8

const (
	SettingsSectionTheme SettingsSection = iota
)

type SettingsItem struct {
	Name  string
	Field SettingsSection
}
