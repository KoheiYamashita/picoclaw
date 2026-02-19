package io.picoclaw.android.core.domain.repository

import io.picoclaw.android.core.domain.model.TtsEngineInfo
import io.picoclaw.android.core.domain.model.TtsVoiceInfo
import kotlinx.coroutines.flow.StateFlow

interface TtsCatalogRepository {
    val availableEngines: StateFlow<List<TtsEngineInfo>>
    val availableVoices: StateFlow<List<TtsVoiceInfo>>
}
