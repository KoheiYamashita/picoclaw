package io.clawdroid.feature.chat

import io.clawdroid.core.domain.model.TtsConfig
import io.clawdroid.core.domain.model.TtsEngineInfo
import io.clawdroid.core.domain.model.TtsVoiceInfo
import io.clawdroid.core.domain.repository.TtsCatalogRepository
import io.clawdroid.core.domain.repository.TtsSettingsRepository
import io.clawdroid.feature.chat.voice.TextToSpeechWrapper
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.flowOf
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class SettingsViewModelTest {

    private val testDispatcher = UnconfinedTestDispatcher()

    private val ttsConfigFlow = MutableStateFlow(TtsConfig())
    private val enginesFlow = MutableStateFlow<List<TtsEngineInfo>>(emptyList())
    private val voicesFlow = MutableStateFlow<List<TtsVoiceInfo>>(emptyList())

    private val ttsSettingsRepository = mockk<TtsSettingsRepository>(relaxed = true) {
        every { ttsConfig } returns ttsConfigFlow
    }
    private val ttsCatalogRepository = mockk<TtsCatalogRepository> {
        every { availableEngines } returns enginesFlow
        every { availableVoices } returns voicesFlow
    }
    private val ttsWrapper = mockk<TextToSpeechWrapper>(relaxed = true)

    private lateinit var viewModel: SettingsViewModel

    @BeforeEach
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        viewModel = SettingsViewModel(ttsSettingsRepository, ttsCatalogRepository, ttsWrapper)
    }

    @AfterEach
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `init collects ttsConfig`() {
        val config = TtsConfig(enginePackageName = "com.google.tts", speechRate = 1.5f)
        ttsConfigFlow.value = config

        assertEquals(config, viewModel.uiState.value.ttsConfig)
    }

    @Test
    fun `init collects available engines`() {
        val engines = listOf(TtsEngineInfo("com.google.tts", "Google TTS"))
        enginesFlow.value = engines

        assertEquals(engines, viewModel.uiState.value.availableEngines)
    }

    @Test
    fun `init collects available voices`() {
        val voices = listOf(TtsVoiceInfo("en-us-1", "English US", "en-US"))
        voicesFlow.value = voices

        assertEquals(voices, viewModel.uiState.value.availableVoices)
    }

    @Test
    fun `onEngineSelected delegates to repository`() = runTest {
        viewModel.onEngineSelected("com.google.tts")

        coVerify { ttsSettingsRepository.updateEngine("com.google.tts") }
    }

    @Test
    fun `onVoiceSelected delegates to repository`() = runTest {
        viewModel.onVoiceSelected("en-us-x-sfg")

        coVerify { ttsSettingsRepository.updateVoiceName("en-us-x-sfg") }
    }

    @Test
    fun `onSpeechRateChanged delegates to repository`() = runTest {
        viewModel.onSpeechRateChanged(1.5f)

        coVerify { ttsSettingsRepository.updateSpeechRate(1.5f) }
    }

    @Test
    fun `onPitchChanged delegates to repository`() = runTest {
        viewModel.onPitchChanged(0.8f)

        coVerify { ttsSettingsRepository.updatePitch(0.8f) }
    }

    @Test
    fun `onTestSpeak calls ttsWrapper and manages isTesting`() = runTest {
        coEvery { ttsWrapper.speak(any()) } returns true

        viewModel.onTestSpeak()

        coVerify { ttsWrapper.speak("これはテスト音声です。This is a test.") }
        assertFalse(viewModel.uiState.value.isTesting)
    }
}
