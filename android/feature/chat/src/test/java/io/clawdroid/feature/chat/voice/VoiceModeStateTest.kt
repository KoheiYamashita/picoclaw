package io.clawdroid.feature.chat.voice

import io.clawdroid.core.domain.model.VoicePhase
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Test

class VoiceModeStateTest {

    @Test
    fun `default isActive is false`() {
        val state = VoiceModeState()

        assertFalse(state.isActive)
    }

    @Test
    fun `default phase is IDLE`() {
        val state = VoiceModeState()

        assertEquals(VoicePhase.IDLE, state.phase)
    }

    @Test
    fun `default recognizedText is empty`() {
        val state = VoiceModeState()

        assertEquals("", state.recognizedText)
    }

    @Test
    fun `default responseText is empty`() {
        val state = VoiceModeState()

        assertEquals("", state.responseText)
    }

    @Test
    fun `default statusText is null`() {
        val state = VoiceModeState()

        assertNull(state.statusText)
    }

    @Test
    fun `default errorMessage is null`() {
        val state = VoiceModeState()

        assertNull(state.errorMessage)
    }

    @Test
    fun `default amplitudeNormalized is 0`() {
        val state = VoiceModeState()

        assertEquals(0f, state.amplitudeNormalized)
    }

    @Test
    fun `default isCameraActive is false`() {
        val state = VoiceModeState()

        assertFalse(state.isCameraActive)
    }

    @Test
    fun `default isScreenCaptureActive is false`() {
        val state = VoiceModeState()

        assertFalse(state.isScreenCaptureActive)
    }

    @Test
    fun `default chatHistory is empty`() {
        val state = VoiceModeState()

        assertTrue(state.chatHistory.isEmpty())
    }
}
