package io.clawdroid.core.domain.model

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class TtsConfigTest {

    @Test
    fun `default enginePackageName is null`() {
        val config = TtsConfig()

        assertNull(config.enginePackageName)
    }

    @Test
    fun `default voiceName is null`() {
        val config = TtsConfig()

        assertNull(config.voiceName)
    }

    @Test
    fun `default speechRate is 1_0`() {
        val config = TtsConfig()

        assertEquals(1.0f, config.speechRate)
    }

    @Test
    fun `default pitch is 1_0`() {
        val config = TtsConfig()

        assertEquals(1.0f, config.pitch)
    }

    @Test
    fun `custom values are preserved`() {
        val config = TtsConfig(
            enginePackageName = "com.google.tts",
            voiceName = "en-us-x-sfg",
            speechRate = 1.5f,
            pitch = 0.8f,
        )

        assertEquals("com.google.tts", config.enginePackageName)
        assertEquals("en-us-x-sfg", config.voiceName)
        assertEquals(1.5f, config.speechRate)
        assertEquals(0.8f, config.pitch)
    }
}
