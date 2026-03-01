package io.clawdroid.feature.chat

import io.clawdroid.core.domain.model.ConnectionState
import io.clawdroid.feature.chat.voice.VoiceModeState
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Test

class ChatUiStateTest {

    @Test
    fun `default state has empty messages`() {
        val state = ChatUiState()

        assertTrue(state.messages.isEmpty())
    }

    @Test
    fun `default state has DISCONNECTED connection`() {
        val state = ChatUiState()

        assertEquals(ConnectionState.DISCONNECTED, state.connectionState)
    }

    @Test
    fun `default state has empty inputText`() {
        val state = ChatUiState()

        assertEquals("", state.inputText)
    }

    @Test
    fun `default state has empty pendingImages`() {
        val state = ChatUiState()

        assertTrue(state.pendingImages.isEmpty())
    }

    @Test
    fun `default state has no error`() {
        val state = ChatUiState()

        assertNull(state.error)
    }

    @Test
    fun `default state has no statusLabel`() {
        val state = ChatUiState()

        assertNull(state.statusLabel)
    }

    @Test
    fun `default state has default voiceModeState`() {
        val state = ChatUiState()

        assertEquals(VoiceModeState(), state.voiceModeState)
    }

    @Test
    fun `default isLoadingMore is false`() {
        val state = ChatUiState()

        assertEquals(false, state.isLoadingMore)
    }

    @Test
    fun `default canLoadMore is true`() {
        val state = ChatUiState()

        assertEquals(true, state.canLoadMore)
    }
}
