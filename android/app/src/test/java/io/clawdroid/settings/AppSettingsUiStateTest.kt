package io.clawdroid.settings

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

class AppSettingsUiStateTest {

    @Nested
    inner class HttpPortError {

        @Test
        fun `empty port returns null`() {
            val state = AppSettingsUiState(httpPort = "")

            assertNull(state.httpPortError)
        }

        @Test
        fun `valid port returns null`() {
            val state = AppSettingsUiState(httpPort = "8080")

            assertNull(state.httpPortError)
        }

        @Test
        fun `non-numeric returns Invalid number`() {
            val state = AppSettingsUiState(httpPort = "abc")

            assertEquals("Invalid number", state.httpPortError)
        }

        @Test
        fun `port 0 returns range error`() {
            val state = AppSettingsUiState(httpPort = "0")

            assertEquals("1-65535", state.httpPortError)
        }

        @Test
        fun `port 65536 returns range error`() {
            val state = AppSettingsUiState(httpPort = "65536")

            assertEquals("1-65535", state.httpPortError)
        }

        @Test
        fun `port 1 is valid`() {
            val state = AppSettingsUiState(httpPort = "1")

            assertNull(state.httpPortError)
        }

        @Test
        fun `port 65535 is valid`() {
            val state = AppSettingsUiState(httpPort = "65535")

            assertNull(state.httpPortError)
        }
    }

    @Nested
    inner class HasErrors {

        @Test
        fun `false when port is valid`() {
            val state = AppSettingsUiState(httpPort = "8080")

            assertFalse(state.hasErrors)
        }

        @Test
        fun `true when port has error`() {
            val state = AppSettingsUiState(httpPort = "abc")

            assertTrue(state.hasErrors)
        }

        @Test
        fun `false when port is empty`() {
            val state = AppSettingsUiState(httpPort = "")

            assertFalse(state.hasErrors)
        }
    }

    @Test
    fun `default values`() {
        val state = AppSettingsUiState()

        assertEquals("", state.apiKey)
        assertEquals("18790", state.httpPort)
        assertFalse(state.saving)
        assertNull(state.error)
    }
}
