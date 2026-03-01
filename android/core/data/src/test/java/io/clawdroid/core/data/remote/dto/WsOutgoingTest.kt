package io.clawdroid.core.data.remote.dto

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class WsOutgoingTest {

    private val json = Json { ignoreUnknownKeys = true }

    @Test
    fun `round-trip with minimal fields`() {
        val dto = WsOutgoing(content = "Hello world")

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<WsOutgoing>(jsonStr)

        assertEquals("Hello world", decoded.content)
        assertNull(decoded.type)
    }

    @Test
    fun `round-trip with type`() {
        val dto = WsOutgoing(content = "status text", type = "status")

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<WsOutgoing>(jsonStr)

        assertEquals("status text", decoded.content)
        assertEquals("status", decoded.type)
    }

    @Test
    fun `deserialize from server-like JSON`() {
        val jsonStr = """{"content":"Thinking...","type":"status"}"""

        val dto = json.decodeFromString<WsOutgoing>(jsonStr)

        assertEquals("Thinking...", dto.content)
        assertEquals("status", dto.type)
    }

    @Test
    fun `deserialize ignores unknown keys`() {
        val jsonStr = """{"content":"msg","type":null,"extra":"ignored"}"""

        val dto = json.decodeFromString<WsOutgoing>(jsonStr)

        assertEquals("msg", dto.content)
    }
}
