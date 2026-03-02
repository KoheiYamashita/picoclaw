package io.clawdroid.core.data.repository

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
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
        val LISTEN_BEEP = booleanPreferencesKey("listen_beep_enabled")
    }

    override val sttConfig: Flow<SttConfig> = context.sttDataStore.data.map { prefs ->
        SttConfig(
            listenBeepEnabled = prefs[Keys.LISTEN_BEEP] ?: true
        )
    }

    override suspend fun updateListenBeepEnabled(enabled: Boolean) {
        context.sttDataStore.edit { it[Keys.LISTEN_BEEP] = enabled }
    }
}
