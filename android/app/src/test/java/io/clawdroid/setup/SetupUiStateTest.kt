package io.clawdroid.setup

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

class SetupUiStateTest {

    @Nested
    inner class GatewayPortError {

        @Test
        fun `empty port returns null`() {
            val state = SetupUiState(gatewayPort = "")

            assertNull(state.gatewayPortError)
        }

        @Test
        fun `valid port returns null`() {
            val state = SetupUiState(gatewayPort = "8080")

            assertNull(state.gatewayPortError)
        }

        @Test
        fun `non-numeric returns Invalid number`() {
            val state = SetupUiState(gatewayPort = "abc")

            assertEquals("Invalid number", state.gatewayPortError)
        }

        @Test
        fun `port 0 returns range error`() {
            val state = SetupUiState(gatewayPort = "0")

            assertEquals("1-65535", state.gatewayPortError)
        }

        @Test
        fun `port 65536 returns range error`() {
            val state = SetupUiState(gatewayPort = "65536")

            assertEquals("1-65535", state.gatewayPortError)
        }

        @Test
        fun `port 1 is valid`() {
            val state = SetupUiState(gatewayPort = "1")

            assertNull(state.gatewayPortError)
        }

        @Test
        fun `port 65535 is valid`() {
            val state = SetupUiState(gatewayPort = "65535")

            assertNull(state.gatewayPortError)
        }
    }

    @Nested
    inner class CanProceedStep1 {

        @Test
        fun `true when port and apiKey are valid`() {
            val state = SetupUiState(gatewayPort = "18790", gatewayApiKey = "key-123")

            assertTrue(state.canProceedStep1)
        }

        @Test
        fun `false when port is empty`() {
            val state = SetupUiState(gatewayPort = "", gatewayApiKey = "key-123")

            assertFalse(state.canProceedStep1)
        }

        @Test
        fun `false when apiKey is empty`() {
            val state = SetupUiState(gatewayPort = "18790", gatewayApiKey = "")

            assertFalse(state.canProceedStep1)
        }

        @Test
        fun `false when port has error`() {
            val state = SetupUiState(gatewayPort = "99999", gatewayApiKey = "key")

            assertFalse(state.canProceedStep1)
        }
    }

    @Test
    fun `default values`() {
        val state = SetupUiState()

        assertEquals(0, state.currentStep)
        assertFalse(state.loading)
        assertNull(state.error)
        assertEquals("18790", state.gatewayPort)
        assertEquals("", state.gatewayApiKey)
        assertFalse(state.step1Done)
        assertEquals("", state.llmModel)
        assertFalse(state.step2Skipped)
        assertEquals("", state.workspace)
        assertFalse(state.step3Skipped)
    }
}
