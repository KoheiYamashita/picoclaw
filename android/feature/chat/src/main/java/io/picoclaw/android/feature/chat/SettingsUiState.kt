package io.picoclaw.android.feature.chat

import io.picoclaw.android.core.domain.model.TtsConfig
import io.picoclaw.android.core.domain.model.TtsEngineInfo
import io.picoclaw.android.core.domain.model.TtsVoiceInfo

data class SettingsUiState(
    val ttsConfig: TtsConfig = TtsConfig(),
    val availableEngines: List<TtsEngineInfo> = emptyList(),
    val availableVoices: List<TtsVoiceInfo> = emptyList(),
    val isTesting: Boolean = false
)
