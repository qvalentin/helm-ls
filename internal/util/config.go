package util

type HelmlsConfiguration struct {
	YamllsConfiguration YamllsConfiguration `json:"yamlls,omitempty"`
}

type YamllsConfiguration struct {
	Enabled        bool           `json:"enabled,omitempty"`
	Path           string         `json:"path,omitempty"`
	YamllsSettings YamllsSettings `json:"config,omitempty"`
}

var DefaultConfig = HelmlsConfiguration{
	YamllsConfiguration: YamllsConfiguration{
		Enabled:        true,
		Path:           "yaml-language-server",
		YamllsSettings: DefaultYamllsSettings,
	},
}

type YamllsSettings struct {
	Schemas    map[string]string `json:"schemas"`
	Completion bool              `json:"completion"`
	Hover      bool              `json:"hover"`
}

var DefaultYamllsSettings = YamllsSettings{
	Schemas:    map[string]string{"kubernetes": "**"},
	Completion: true,
	Hover:      true,
}
