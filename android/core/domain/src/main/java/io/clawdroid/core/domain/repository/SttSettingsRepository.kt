package io.clawdroid.core.domain.repository

import io.clawdroid.core.domain.model.SttConfig
import kotlinx.coroutines.flow.Flow

interface SttSettingsRepository {
    val sttConfig: Flow<SttConfig>
    suspend fun updateListenBeepUri(uri: String)
}
