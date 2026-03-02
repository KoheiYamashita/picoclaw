package io.clawdroid.core.data.repository

import android.content.Context
import android.media.RingtoneManager
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import io.clawdroid.core.domain.model.SttConfig
import io.clawdroid.core.domain.repository.SttSettingsRepository
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private val Context.sttDataStore: DataStore<Preferences> by preferencesDataStore(name = "stt_settings")

class SttSettingsRepositoryImpl(
    private val context: Context
) : SttSettingsRepository {

    private object Keys {
        val LISTEN_BEEP_URI = stringPreferencesKey("listen_beep_uri")
        val LISTEN_BEEP_LEGACY = booleanPreferencesKey("listen_beep_enabled")
    }

    private val defaultUri: String
        get() = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION).toString()

    override val sttConfig: Flow<SttConfig> = context.sttDataStore.data.map { prefs ->
        val uri = if (Keys.LISTEN_BEEP_URI in prefs) {
            prefs[Keys.LISTEN_BEEP_URI] ?: defaultUri
        } else {
            // Migration: read legacy boolean key
            val legacyEnabled = prefs[Keys.LISTEN_BEEP_LEGACY] ?: true
            if (legacyEnabled) defaultUri else ""
        }
        SttConfig(listenBeepUri = uri)
    }

    override suspend fun updateListenBeepUri(uri: String) {
        context.sttDataStore.edit {
            it[Keys.LISTEN_BEEP_URI] = uri
            it.remove(Keys.LISTEN_BEEP_LEGACY)
        }
    }
}
