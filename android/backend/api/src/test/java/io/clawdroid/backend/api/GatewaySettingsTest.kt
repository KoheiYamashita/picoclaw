package io.clawdroid.backend.api

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class GatewaySettingsTest {

    @Test
    fun `default httpPort is 18790`() {
        val settings = GatewaySettings()

        assertEquals(18790, settings.httpPort)
    }

    @Test
    fun `default apiKey is empty`() {
        val settings = GatewaySettings()

        assertEquals("", settings.apiKey)
    }

    @Test
    fun `httpBaseUrl uses httpPort`() {
        val settings = GatewaySettings(httpPort = 8080)

        assertEquals("http://127.0.0.1:8080", settings.httpBaseUrl)
    }

    @Test
    fun `httpBaseUrl with default port`() {
        val settings = GatewaySettings()

        assertEquals("http://127.0.0.1:18790", settings.httpBaseUrl)
    }

    @Test
    fun `custom values are preserved`() {
        val settings = GatewaySettings(httpPort = 9999, apiKey = "secret-key")

        assertEquals(9999, settings.httpPort)
        assertEquals("secret-key", settings.apiKey)
    }
}
