package io.clawdroid.core.data.remote.dto

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class WsIncomingTest {

    private val json = Json { ignoreUnknownKeys = true }

    @Test
    fun `serialize minimal WsIncoming`() {
        val dto = WsIncoming(content = "Hello")

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<WsIncoming>(jsonStr)

        assertEquals("Hello", decoded.content)
        assertNull(decoded.senderId)
        assertNull(decoded.images)
        assertNull(decoded.inputMode)
        assertNull(decoded.type)
        assertNull(decoded.requestId)
    }

    @Test
    fun `serialize full WsIncoming`() {
        val dto = WsIncoming(
            content = "Hi",
            senderId = "user1",
            images = listOf("base64img"),
            inputMode = "voice",
            type = "tool_response",
            requestId = "req-123",
        )

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<WsIncoming>(jsonStr)

        assertEquals("Hi", decoded.content)
        assertEquals("user1", decoded.senderId)
        assertEquals(listOf("base64img"), decoded.images)
        assertEquals("voice", decoded.inputMode)
        assertEquals("tool_response", decoded.type)
        assertEquals("req-123", decoded.requestId)
    }

    @Test
    fun `deserialize with SerialName mapping`() {
        val jsonStr = """{"content":"test","sender_id":"s1","input_mode":"text","request_id":"r1"}"""

        val dto = json.decodeFromString<WsIncoming>(jsonStr)

        assertEquals("test", dto.content)
        assertEquals("s1", dto.senderId)
        assertEquals("text", dto.inputMode)
        assertEquals("r1", dto.requestId)
    }

    @Test
    fun `deserialize ignores unknown keys`() {
        val jsonStr = """{"content":"test","unknown_field":"value"}"""

        val dto = json.decodeFromString<WsIncoming>(jsonStr)

        assertEquals("test", dto.content)
    }
}
