package official

import (
	"encoding/json"
	"fmt"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webcore/registry/upload/checker"
	"github.com/coscms/webfront/library/settings"
	"github.com/webx-top/echo"
)

func onRequestGet(ctx echo.Context) error {
	baseCfg := config.FromFile().Settings().Base
	formModel := &settings.SettingsMultilingual{
		SiteName:            baseCfg.String(`siteName`),
		SiteSlogan:          baseCfg.String(`siteSlogan`),
		SiteMetaKeywords:    baseCfg.String(`siteMetaKeywords`),
		SiteMetaDescription: baseCfg.String(`siteMetaDescription`),
		SiteAnnouncement:    baseCfg.String(`siteAnnouncement`),
	}
	form := formbuilder.New(ctx,
		formModel,
		formbuilder.ConfigFile(`official/settings/base`),
		formbuilder.MultilingualFields(`SiteName`, `SiteSlogan`, `SiteMetaKeywords`, `SiteMetaDescription`, `SiteAnnouncement`),
	)
	baseMultilinguals := settings.GetBaseMultilinguals()
	setInput := func(langCode string, item *settings.SettingsMultilingual) {
		form.SetLangInput(langCode, `siteName`, item.SiteName)
		form.SetLangInput(langCode, `siteSlogan`, item.SiteSlogan)
		form.SetLangInput(langCode, `siteMetaKeywords`, item.SiteMetaKeywords)
		form.SetLangInput(langCode, `siteMetaDescription`, item.SiteMetaDescription)
		form.SetLangInput(langCode, `siteAnnouncement`, item.SiteAnnouncement)
	}
	setInput(form.Languages().Default, formModel)
	if baseMultilinguals != nil {
		for langCode, item := range *baseMultilinguals {
			setInput(langCode, &item)
		}
	}
	form.Generate()
	form.Config().SetDefaultValue(func(fieldName string) string {
		return ctx.Form(fieldName)
	})

	uploadURL := checker.BackendUploadURL(``, `group`, `base`, `key`, `siteAnnouncement`)
	finderURL := `!` + backend.URLFor(`/finder`) + `?multiple=1`
	siteAnnouncementField := form.MultilingualField(config.FromFile().Language.Default, `siteAnnouncement`, `siteAnnouncement`)
	siteAnnouncementField.SetParam(`action`, finderURL)
	siteAnnouncementField.SetParam(`data-upload-url`, uploadURL)
	return nil
}

func onRequestPost(ctx echo.Context) error {
	formModel := &settings.SettingsMultilingual{}
	form := formbuilder.New(ctx,
		formModel,
		formbuilder.ConfigFile(`official/settings/base`),
		formbuilder.AllowedNames(`siteName`, `siteSlogan`, `siteMetaKeywords`, `siteMetaDescription`, `siteAnnouncement`),
		formbuilder.MultilingualFields(`SiteName`, `SiteSlogan`, `SiteMetaKeywords`, `SiteMetaDescription`, `SiteAnnouncement`),
	)
	form.OnPost(func() error {
		ctx.Request().Form().Set(`base[siteName][value]`, formModel.SiteName)
		ctx.Request().Form().Set(`base[siteSlogan][value]`, formModel.SiteSlogan)
		ctx.Request().Form().Set(`base[siteMetaKeywords][value]`, formModel.SiteMetaKeywords)
		ctx.Request().Form().Set(`base[siteMetaDescription][value]`, formModel.SiteMetaDescription)
		ctx.Request().Form().Set(`base[siteAnnouncement][value]`, formModel.SiteAnnouncement)
		mls := settings.SettingsMultilinguals{}
		ldf := form.Languages().Default
		for _, langCode := range form.Languages().AllList {
			if langCode == ldf {
				continue
			}
			mls[langCode] = settings.SettingsMultilingual{
				SiteName:            form.GetLangInput(langCode, `siteName`),
				SiteSlogan:          form.GetLangInput(langCode, `siteSlogan`),
				SiteMetaKeywords:    form.GetLangInput(langCode, `siteMetaKeywords`),
				SiteMetaDescription: form.GetLangInput(langCode, `siteMetaDescription`),
				SiteAnnouncement:    form.GetLangInput(langCode, `siteAnnouncement`),
			}
		}
		b, err := json.Marshal(mls)
		if err != nil {
			return fmt.Errorf(`failed to json marshal(%#v): %w`, mls, err)
		}
		ctx.Request().Form().Set(`base[multilingual][value]`, string(b))
		return nil
	})
	form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}

	return nil
}
