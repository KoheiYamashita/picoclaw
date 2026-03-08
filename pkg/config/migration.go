package config

// ConfigVersion is the current config schema version.
// Increment this when adding new fields that must appear in existing config files.
const ConfigVersion = 3

// migrateConfig runs version-based migrations on cfg.
// Returns true if migrations were applied and config should be re-saved.
func migrateConfig(cfg *Config) bool {
	if cfg.Version >= ConfigVersion {
		return false
	}

	migrations := []func(*Config){
		migrateV0ToV1,
		migrateV1ToV2,
		migrateV2ToV3,
	}

	for i := cfg.Version; i < ConfigVersion && i < len(migrations); i++ {
		migrations[i](cfg)
	}

	cfg.Version = ConfigVersion
	return true
}

func migrateV0ToV1(cfg *Config) {
	// queue_messages: Go zero value (false) is the correct default.
	// Version bump + re-save writes the new field into config.json.
}

func migrateV1ToV2(cfg *Config) {
	cfg.Agents.Defaults.ShowErrors = true
	cfg.Agents.Defaults.ShowWarnings = true
}

func migrateV2ToV3(cfg *Config) {
	cfg.Tools.Android.Categories = AndroidCategories{
		Alarm:         true,
		Calendar:      true,
		Contacts:      false,
		Communication: false,
		Media:         true,
		Navigation:    true,
		DeviceControl: true,
		Settings:      true,
		Web:           true,
		Clipboard:     true,
	}

	// Set all actions to enabled by default, then disable any that were
	// listed in the old DisabledActions string slice.
	actions := DefaultAndroidActions()

	disabled := make(map[string]bool, len(cfg.Tools.Android.DisabledActions))
	for _, d := range cfg.Tools.Android.DisabledActions {
		disabled[d] = true
	}

	if len(disabled) > 0 {
		disableAndroidActions(&actions, disabled)
	}

	cfg.Tools.Android.Actions = actions
	cfg.Tools.Android.DisabledActions = nil
}

// disableAndroidActions sets action fields to false for any action name present in disabled.
func disableAndroidActions(a *AndroidActions, disabled map[string]bool) {
	// Alarm
	if disabled["set_alarm"] {
		a.Alarm.SetAlarm = false
	}
	if disabled["set_timer"] {
		a.Alarm.SetTimer = false
	}
	if disabled["dismiss_alarm"] {
		a.Alarm.DismissAlarm = false
	}
	if disabled["show_alarms"] {
		a.Alarm.ShowAlarms = false
	}
	// Calendar
	if disabled["create_event"] {
		a.Calendar.CreateEvent = false
	}
	if disabled["query_events"] {
		a.Calendar.QueryEvents = false
	}
	if disabled["update_event"] {
		a.Calendar.UpdateEvent = false
	}
	if disabled["delete_event"] {
		a.Calendar.DeleteEvent = false
	}
	if disabled["list_calendars"] {
		a.Calendar.ListCalendars = false
	}
	if disabled["add_reminder"] {
		a.Calendar.AddReminder = false
	}
	// Contacts
	if disabled["search_contacts"] {
		a.Contacts.SearchContacts = false
	}
	if disabled["get_contact_detail"] {
		a.Contacts.GetContactDetail = false
	}
	if disabled["add_contact"] {
		a.Contacts.AddContact = false
	}
	// Communication
	if disabled["dial"] {
		a.Communication.Dial = false
	}
	if disabled["compose_sms"] {
		a.Communication.ComposeSMS = false
	}
	if disabled["compose_email"] {
		a.Communication.ComposeEmail = false
	}
	// Media
	if disabled["media_play_pause"] {
		a.Media.PlayPause = false
	}
	if disabled["media_next"] {
		a.Media.Next = false
	}
	if disabled["media_previous"] {
		a.Media.Previous = false
	}
	if disabled["play_music_search"] {
		a.Media.PlayMusicSearch = false
	}
	// Navigation
	if disabled["navigate"] {
		a.Navigation.Navigate = false
	}
	if disabled["search_nearby"] {
		a.Navigation.SearchNearby = false
	}
	if disabled["show_map"] {
		a.Navigation.ShowMap = false
	}
	if disabled["get_current_location"] {
		a.Navigation.GetCurrentLocation = false
	}
	// Device Control
	if disabled["flashlight"] {
		a.DeviceControl.Flashlight = false
	}
	if disabled["set_volume"] {
		a.DeviceControl.SetVolume = false
	}
	if disabled["set_ringer_mode"] {
		a.DeviceControl.SetRingerMode = false
	}
	if disabled["set_dnd"] {
		a.DeviceControl.SetDND = false
	}
	if disabled["set_brightness"] {
		a.DeviceControl.SetBrightness = false
	}
	// Settings
	if disabled["open_settings"] {
		a.Settings.OpenSettings = false
	}
	// Web
	if disabled["open_url"] {
		a.Web.OpenURL = false
	}
	if disabled["web_search"] {
		a.Web.WebSearch = false
	}
	// Clipboard
	if disabled["clipboard_copy"] {
		a.Clipboard.ClipboardCopy = false
	}
	if disabled["clipboard_read"] {
		a.Clipboard.ClipboardRead = false
	}
}
